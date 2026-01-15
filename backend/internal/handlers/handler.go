package handlers

import (
	"poke-atlas/web-service/internal/repository"
)

type Handler struct {
	repo repository.Repository
}

func NewHandler(repository repository.Repository) *Handler {
	return &Handler{
		repo: repository,
	}
}
