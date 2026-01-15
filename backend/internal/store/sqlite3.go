package store

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"poke-atlas/web-service/internal/model"

	_ "github.com/mattn/go-sqlite3"
)

type sqliteDatabase struct {
	db *sql.DB
}

func CreateSqliteDatabase() Database {
	db, err := sql.Open("sqlite3", "./pokedb.db")

	if err != nil {
		log.Fatal("Failed to open database", err.Error())
	}

	database := &sqliteDatabase{
		db: db,
	}
	return database
}

// Database interface implementation

func (s *sqliteDatabase) AddPokemon(ctx context.Context, pokemon model.Pokemon) error {
	jsonData, err := json.Marshal(pokemon)
	if err != nil {
		return fmt.Errorf("marshaling pokemon: %w", err)
	}

	query := `INSERT OR REPLACE INTO pokemon (id, name, data) VALUES (?, ?, ?)`
	_, err = s.db.ExecContext(ctx, query, pokemon.ID, pokemon.Name, string(jsonData))
	return err
}

func (s *sqliteDatabase) GetPokemon(ctx context.Context, name string) (model.Pokemon, error) {
	query := `SELECT data FROM pokemon WHERE name = ?`

	var jsonData string

	err := s.db.QueryRowContext(ctx, query, name).Scan(&jsonData)

	if err == sql.ErrNoRows {
		return model.Pokemon{}, fmt.Errorf("pokemon not found")
	}
	if err != nil {
		return model.Pokemon{}, err
	}

	var pokemon model.Pokemon

	err = json.Unmarshal([]byte(jsonData), &pokemon)
	return pokemon, err
}

func (s *sqliteDatabase) InitDB() error {
	query := `
	CREATE TABLE IF NOT EXISTS pokemon (
	id INTEGER PRIMARY KEY,
	name TEXT UNIQUE NOT NULL,
	data TEXT NOT NULL
	);
	
	CREATE INDEX IF NOT EXISTS idx_pokemon_name on pokemon(name);
	`

	_, err := s.db.Exec(query)

	return err
}

func (s *sqliteDatabase) Close() error {
	return s.db.Close()
}
