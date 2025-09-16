package handlers

import "warehouse/internal/db"

type Handler struct {
	storage *db.DB
}

func NewHandler(db *db.DB) *Handler {
	return &Handler{
		storage: db,
	}
}

// Ты реализовать API ручки
