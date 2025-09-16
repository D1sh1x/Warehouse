package handlers

import (
	"database/sql"
	"net/http"
	"strconv"
	"time"
	"warehouse/internal/db"
	"warehouse/internal/middleware"
	"warehouse/internal/models"

	"github.com/golang-jwt/jwt/v5"
	"github.com/wb-go/wbf/ginext"
	"github.com/wb-go/wbf/zlog"
)

type Handler struct {
	storage   *db.DB
	jwtSecret []byte
}

func NewHandler(db *db.DB, jwtSecret []byte) *Handler {
	return &Handler{
		storage:   db,
		jwtSecret: jwtSecret,
	}
}

func (h *Handler) CreateItem(c *ginext.Context) {
	var item models.Item
	if err := c.ShouldBindJSON(&item); err != nil {
		zlog.Logger.Error().Err(err).Msg("Failed to bind JSON")
		c.JSON(http.StatusBadRequest, ginext.H{"error": err.Error()})
		return
	}

	err := h.storage.CreateItem(c.Request.Context(), &item)
	if err != nil {
		zlog.Logger.Error().Err(err).Msg("Failed to create item")
		c.JSON(http.StatusInternalServerError, ginext.H{"error": err.Error()})
		return
	}

	zlog.Logger.Info().Msg("Item created successfully")
	c.JSON(http.StatusOK, item)
}

func (h *Handler) GetItems(c *ginext.Context) {

	items, err := h.storage.GetItems(c.Request.Context())

	if err != nil {
		zlog.Logger.Error().Err(err).Msg("Failed to get items")
		c.JSON(http.StatusInternalServerError, ginext.H{"error": err.Error()})
		return
	}

	zlog.Logger.Info().Int("count", len(items)).Msg("Items retrieved")
	c.JSON(http.StatusOK, items)
}

func (h *Handler) UpdateItem(c *ginext.Context) {

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		zlog.Logger.Error().Err(err).Msg("Invalid item ID")
		c.JSON(http.StatusBadRequest, ginext.H{"error": "invalid id"})
		return
	}

	var item models.Item

	if err := c.ShouldBindJSON(&item); err != nil {
		zlog.Logger.Error().Err(err).Msg("Failed to bind JSON")
		c.JSON(http.StatusBadRequest, ginext.H{"error": err.Error()})
		return
	}

	if err := h.storage.UpdateItem(c.Request.Context(), id, &item); err != nil {
		zlog.Logger.Error().Err(err).Msg("Failed to update item")
		c.JSON(http.StatusInternalServerError, ginext.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, item)
}

func (h *Handler) DeleteItem(c *ginext.Context) {
	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		zlog.Logger.Error().Err(err).Msg("Invalid item ID")
		c.JSON(http.StatusBadRequest, ginext.H{"error": "invalid id"})
		return
	}

	if err := h.storage.DeleteItem(c.Request.Context(), id); err != nil {
		zlog.Logger.Error().Err(err).Msg("Failed to delete item")
		c.JSON(http.StatusInternalServerError, ginext.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusOK)
}

func (h *Handler) GetHistory(c *ginext.Context) {

	history, err := h.storage.GetHistory(c.Request.Context())

	if err != nil {
		zlog.Logger.Error().Err(err).Msg("Failed to get history")
		c.JSON(http.StatusInternalServerError, ginext.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, history)
}

func (h *Handler) ItemHistory(c *ginext.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		zlog.Logger.Error().Err(err).Msg("Invalid item ID")
		c.JSON(http.StatusBadRequest, ginext.H{"error": "invalid item id"})
		return
	}

	history, err := h.storage.GetItemHistory(c.Request.Context(), id)
	if err != nil {
		zlog.Logger.Error().Err(err).Msg("Failed to get item history")
		c.JSON(http.StatusInternalServerError, ginext.H{"error": "failed to get history"})
		return
	}

	c.JSON(http.StatusOK, history)
}

func (h *Handler) Login(c *ginext.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, ginext.H{"error": err.Error()})
		return
	}

	dbUser, err := h.storage.GetUser(c.Request.Context(), user.Username, user.Password)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusUnauthorized, ginext.H{"error": "Invalid credentials"})
		} else {
			c.JSON(http.StatusInternalServerError, ginext.H{"error": "Database error"})
		}
		return
	}

	claims := &middleware.Claims{
		Username: dbUser.Username,
		Role:     dbUser.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			Issuer:    "warehouse",
		},
	}
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := t.SignedString(h.jwtSecret)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ginext.H{"error": "sign jwt failed"})
		return
	}
	c.JSON(http.StatusOK, ginext.H{"token": signed, "role": dbUser.Role})
}
