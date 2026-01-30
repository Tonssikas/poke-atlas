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
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// pokemon
	query := `INSERT OR IGNORE INTO pokemons (id, name, height, weight) VALUES (?, ?, ?, ?)`

	_, err = tx.ExecContext(ctx, query, pokemon.ID, pokemon.Name, pokemon.Height, pokemon.Weight)
	if err != nil {
		return err
	}

	// types and pokemon_types
	stmtType, _ := tx.PrepareContext(ctx, `INSERT OR IGNORE INTO types (name) VALUES (?)`)
	defer stmtType.Close()

	stmtPokemonType, _ := tx.PrepareContext(ctx, `INSERT OR IGNORE INTO pokemon_types (pokemon_id, type_name, slot) VALUES (?, ?, ?)`)
	defer stmtPokemonType.Close()

	for _, t := range pokemon.Types {
		if _, err := stmtType.ExecContext(ctx, t.Type.Name); err != nil {
			return err
		}
		if _, err := stmtPokemonType.ExecContext(ctx, pokemon.ID, t.Type.Name, t.Slot); err != nil {
			return err
		}
	}

	// abilities and pokemon_ability

	stmtAbility, _ := tx.PrepareContext(ctx, `INSERT OR IGNORE INTO abilities (name) VALUES (?)`)
	defer stmtAbility.Close()

	stmtPokemonAbility, _ := tx.PrepareContext(ctx, `INSERT OR IGNORE INTO pokemon_ability (pokemon_id, ability_name, is_hidden) VALUES (?, ?, ?)`)
	defer stmtPokemonAbility.Close()
	for _, a := range pokemon.Abilities {
		if _, err := stmtAbility.ExecContext(ctx, a.Ability.Name); err != nil {
			return err
		}

		if _, err := stmtPokemonAbility.ExecContext(ctx, pokemon.ID, a.Ability.Name, a.IsHidden); err != nil {
			return err
		}
	}

	// moves, move_learn_methods and version_group

	stmtMoves, _ := tx.PrepareContext(ctx, `INSERT OR IGNORE INTO moves (name) VALUES (?)`)
	defer stmtMoves.Close()
	stmtMoveLearnMethods, _ := tx.PrepareContext(ctx, `INSERT OR IGNORE INTO move_learn_methods (learn_method) VALUES (?)`)
	defer stmtMoveLearnMethods.Close()
	stmtVersionGroup, _ := tx.PrepareContext(ctx, `INSERT OR IGNORE INTO version_groups (version_name) VALUES (?)`)
	defer stmtVersionGroup.Close()
	stmtPokemonMoves, _ := tx.PrepareContext(ctx, `INSERT OR IGNORE INTO pokemon_moves (move_name, pokemon_id, version_group, move_learn_method, level_learned_at, move_order) VALUES (?, ?, ?, ?, ?, ?)`)
	defer stmtPokemonMoves.Close()

	for _, m := range pokemon.Moves {
		if _, err := stmtMoves.ExecContext(ctx, m.Move.Name); err != nil {
			return err
		}

		for _, d := range m.VersionGroupDetails {
			if _, err := stmtMoveLearnMethods.ExecContext(ctx, d.MoveLearnMethod.Name); err != nil {
				return err
			}
			if _, err := stmtVersionGroup.ExecContext(ctx, d.VersionGroup.Name); err != nil {
				return err
			}
			if _, err := stmtPokemonMoves.ExecContext(ctx, m.Move.Name, pokemon.ID, d.VersionGroup.Name, d.MoveLearnMethod.Name, d.LevelLearnedAt, d.Order); err != nil {
				return err
			}
		}
	}

	// pokemon stats

	stmtStats, _ := tx.PrepareContext(ctx, `INSERT OR IGNORE INTO stats (name) VALUES (?)`)
	defer stmtStats.Close()
	stmtPokemonStats, _ := tx.PrepareContext(ctx, `INSERT OR IGNORE INTO pokemon_stats (pokemon_id, stat_name, effort, base_stat) VALUES (?, ?, ?, ?)`)
	defer stmtPokemonStats.Close()

	for _, s := range pokemon.Stats {
		if _, err := stmtStats.ExecContext(ctx, s.Stat.Name); err != nil {
			return err
		}
		if _, err := stmtPokemonStats.ExecContext(ctx, pokemon.ID, s.Stat.Name, s.Effort, s.BaseStat); err != nil {
			return err
		}
	}

	return tx.Commit()
}

// TODO: More queries based on frontend needs
// Return a brief summary of pokemon for now
func (s *sqliteDatabase) GetPokemon(ctx context.Context, name string) (model.Pokemon_summary, error) {
	query := `
	SELECT pokemons.id, pokemons.name, pokemons.weight, pokemons.height, json_group_array(pokemon_types.type_name) FROM pokemons
	JOIN pokemon_types ON pokemon_types.pokemon_id = pokemons.id
	WHERE pokemons.name = ?
	GROUP BY pokemons.id, pokemons.name, pokemons.weight, pokemons.height
	`

	var pokemon model.Pokemon_summary

	// Temporary variable to store types in for unmarshaling
	var typesJSON []byte

	err := s.db.QueryRowContext(ctx, query, name).Scan(
		&pokemon.ID,
		&pokemon.Name,
		&pokemon.Weight,
		&pokemon.Height,
		&typesJSON,
	)

	if err == sql.ErrNoRows {
		return model.Pokemon_summary{}, fmt.Errorf("pokemon not found")
	}
	if err != nil {
		return model.Pokemon_summary{}, err
	}

	if err := json.Unmarshal(typesJSON, &pokemon.Types); err != nil {
		return model.Pokemon_summary{}, err
	}

	return pokemon, nil
}

func (s *sqliteDatabase) GetPokemons(ctx context.Context, offset int, limit int) ([]model.Pokemon_summary, error) {
	query := `
	SELECT pokemons.id, pokemons.name, pokemons.weight, pokemons.height, json_group_array(pokemon_types.type_name) as types
	FROM pokemons
	JOIN pokemon_types ON pokemon_types.pokemon_id = pokemons.id 
	WHERE id > ? AND id <= ?
	GROUP BY pokemons.id, pokemons.name, pokemons.weight, pokemons.height
	ORDER BY id
	`
	// Old query

	/*
			`SELECT id, name, weight, height
		            FROM pokemons
					WHERE id > ? AND id <= ?
		            ORDER BY id`
	*/

	rows, err := s.db.QueryContext(ctx, query, offset, offset+limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var pokemons []model.Pokemon_summary
	var typesJSON []byte

	// Increment expected ID by 1 for each iteration
	expectedID := offset + 1

	for rows.Next() {
		var pokemon model.Pokemon_summary
		err := rows.Scan(&pokemon.ID, &pokemon.Name, &pokemon.Weight, &pokemon.Height, &typesJSON)
		if err != nil {
			return nil, err
		}
		// Check for gap in the ID sequence
		// This logic breaks when we hit regional / mega forms at some point after 1000, ID jumps up to 10000
		if pokemon.ID != expectedID {
			// Gap detected - return empty to trigger API fetch
			return []model.Pokemon_summary{}, nil
		}
		json.Unmarshal(typesJSON, &pokemon.Types)
		pokemons = append(pokemons, pokemon)
		expectedID++
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return pokemons, nil
}

func (s *sqliteDatabase) InitDB() error {
	query := `
	PRAGMA foreign_keys = ON;
	CREATE TABLE IF NOT EXISTS pokemons (
	id INTEGER PRIMARY KEY,
	name TEXT UNIQUE NOT NULL,
	weight INTEGER,
	height INTEGER
	);

	CREATE TABLE IF NOT EXISTS types (
	name TEXT PRIMARY KEY
	);

	CREATE TABLE IF NOT EXISTS pokemon_types (
	pokemon_id INTEGER NOT NULL,
	type_name  TEXT NOT NULL,
	slot INTEGER,

	PRIMARY KEY (pokemon_id, type_name),
	FOREIGN KEY (pokemon_id) REFERENCES pokemons(id),
	FOREIGN KEY (type_name) REFERENCES types(name)
	);

	CREATE TABLE IF NOT EXISTS abilities (
	name TEXT PRIMARY KEY
	);

	CREATE TABLE IF NOT EXISTS pokemon_ability (
	pokemon_id INTEGER NOT NULL,
	ability_name TEXT NOT NULL,
	is_hidden INTEGER CHECK (is_hidden IN (0, 1)),
	
	PRIMARY KEY (pokemon_id, ability_name),
	FOREIGN KEY (pokemon_id) REFERENCES pokemons(id),
	FOREIGN KEY (ability_name) REFERENCES abilities(name)
	);

	CREATE TABLE IF NOT EXISTS stats (
	name TEXT PRIMARY KEY
	);

	CREATE TABLE IF NOT EXISTS pokemon_stats (
	pokemon_id INTEGER NOT NULL,
	stat_name TEXT NOT NULL,
	effort INTEGER,
	base_stat INTEGER,

	PRIMARY KEY (pokemon_id, stat_name),
	FOREIGN KEY (pokemon_id) REFERENCES pokemons(id),
	FOREIGN KEY (stat_name) REFERENCES stats(name)
	);

	CREATE TABLE IF NOT EXISTS moves (
	name TEXT PRIMARY KEY
	);

	CREATE TABLE IF NOT EXISTS move_learn_methods (
	learn_method TEXT PRIMARY KEY
	);

	CREATE TABLE IF NOT EXISTS version_groups (
	version_name TEXT PRIMARY KEY
	);

	CREATE TABLE IF NOT EXISTS pokemon_moves (
	move_name TEXT NOT NULL,
	pokemon_id INTEGER NOT NULL,
	version_group TEXT NOT NULL,
	move_learn_method TEXT NOT NULL,
	level_learned_at INTEGER,
	move_order INTEGER,

	PRIMARY KEY (move_name, pokemon_id, version_group, move_learn_method),
	FOREIGN KEY (pokemon_id) REFERENCES pokemons(id),
	FOREIGN KEY (move_name) REFERENCES moves(name),
	FOREIGN KEY (version_group) REFERENCES version_groups(version_name),
	FOREIGN KEY (move_learn_method) REFERENCES move_learn_methods(learn_method)
	);
	`

	_, err := s.db.Exec(query)

	return err
}

func (s *sqliteDatabase) Close() error {
	return s.db.Close()
}
