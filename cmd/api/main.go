package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"warehouse/internal/config"
	"warehouse/internal/db"
	"warehouse/internal/handlers"
	"warehouse/internal/middleware"

	"github.com/wb-go/wbf/ginext"
	"github.com/wb-go/wbf/zlog"
)

func main() {
	zlog.Init()

	cfg, err := config.Load("./config/config.yml")
	if err != nil {
		zlog.Logger.Fatal().Err(err).Msg("Failed to load config")
	}

	database, err := db.NewDB(cfg.DSNString())

	if err != nil {
		zlog.Logger.Fatal().Err(err).Msg("Failed to connect database")
	}

	defer database.DB.Master.Close()

	handler := handlers.NewHandler(database, cfg.JWTSecret)

	r := ginext.New()

	r.Static("/static", "./static")
	r.GET("/", func(c *ginext.Context) {
		c.File("./static/index.html")
	})

	r.POST("/login", handler.Login)

	auth := r.Group("")
	auth.Use(middleware.AuthRequired(cfg.JWTSecret))

	auth.GET("/items", middleware.RequireRole("admin", "manager", "viewer"), handler.GetItems)
	auth.POST("/items", middleware.RequireRole("admin", "manager"), handler.CreateItem)
	auth.PUT("/items/:id", middleware.RequireRole("admin", "manager"), handler.UpdateItem)
	auth.DELETE("/items/:id", middleware.RequireRole("admin"), handler.DeleteItem)

	auth.GET("/items/:id/history", middleware.RequireRole("admin", "manager", "viewer"), handler.ItemHistory)
	auth.GET("/history", middleware.RequireRole("admin", "manager", "viewer"), handler.GetHistory)

	srv := &http.Server{
		Addr:    fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port),
		Handler: r,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			zlog.Logger.Fatal().Err(err).Msg("Failed to start server")
		}
	}()

	zlog.Logger.Info().Str("addr", srv.Addr).Msg("Server started")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	zlog.Logger.Info().Msg("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		zlog.Logger.Fatal().Err(err).Msg("Server forced to shutdown")
	}

	zlog.Logger.Info().Msg("Server exited")
}
