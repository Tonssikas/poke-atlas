package main

import (
	"net/http"
	"poke-atlas/web-service/internal/handlers"
	"poke-atlas/web-service/internal/pokeapi"
	"poke-atlas/web-service/internal/repository"
	"poke-atlas/web-service/internal/store"

	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {

	err := godotenv.Load()

	if err != nil {
		log.Println("Error loading .env file, using default values")
	}

	// Initialize dependencies
	pokeAPIClient := pokeapi.NewPokeAPIClient(http.DefaultClient)
	database := store.CreateSqliteDatabase()
	defer database.Close()

	if err := database.InitDB(); err != nil {
		log.Fatal("Failed to initialize database:", err)
	}
	repository := repository.NewRepository(pokeAPIClient, database)
	handler := handlers.NewHandler(repository)

	router := gin.Default()
	router.SetTrustedProxies(nil)

	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "it works!",
		})
	})

	router.GET("/pokemon/:name", handler.GetPokemonHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	router.Run(fmt.Sprintf(":%s", port))
}
