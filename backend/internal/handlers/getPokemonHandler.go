package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) GetPokemonHandler(c *gin.Context) {
	name := c.Param("name")

	// Validate name is not empty
	if name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "pokemon name is required"})
		return
	}

	pokemon, err := h.repo.GetPokemon(c.Request.Context(), name)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, pokemon)
}
