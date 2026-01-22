package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (h *Handler) GetPokemonsHandler(c *gin.Context) {
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

	pokemons, err := h.repo.GetPokemons(c.Request.Context(), offset)

	c.JSON(http.StatusOK, pokemons)

}
