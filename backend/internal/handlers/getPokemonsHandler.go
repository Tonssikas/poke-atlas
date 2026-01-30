package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (h *Handler) GetPokemonsHandler(c *gin.Context) {
	limitStr := c.DefaultQuery("limit", "20")
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "limit must be a valid integer"})
		return
	}

	if limit <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "limit must be a greater than 0"})
		return
	}

	offsetStr := c.Param("offset")
	offset, err := strconv.Atoi(offsetStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "offset must be a valid integer"})
		return
	}

	// Validate offset is not negative
	if offset < 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "offset cannot be negative"})
		return
	}

	pokemons, err := h.repo.GetPokemons(c.Request.Context(), offset, limit)

	c.JSON(http.StatusOK, pokemons)
}
