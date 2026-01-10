package main

import (
	"context"
	"net/http"
	"poke-atlas/web-service/internal/pokeapi"
	"poke-atlas/web-service/internal/repository"

	"github.com/gin-gonic/gin"
)

func main() {

	repository := repository.NewRepository(pokeapi.NewPokeAPIClient(http.DefaultClient))

	router := gin.Default()

	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "it works!",
		})
	})

	router.GET("/poketest", func(c *gin.Context) {
		pokemon, err := repository.GetPokemon(context.Background(), "pikachu")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, pokemon)
	})

	router.Run("localhost:8080")
}
