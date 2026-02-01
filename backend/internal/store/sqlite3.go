package store

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"poke-atlas/web-service/internal/model"
	"strconv"
	"strings"

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

	// basic pokemon information and sprite
	query := `INSERT OR IGNORE INTO pokemons (id, name, height, weight, sprite_url) VALUES (?, ?, ?, ?, ?)`

	_, err = tx.ExecContext(ctx, query, pokemon.ID, pokemon.Name, pokemon.Height, pokemon.Weight, pokemon.Sprites.FrontDefault)
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

// Return a brief summary of pokemon for now
func (s *sqliteDatabase) GetPokemon(ctx context.Context, name string) (model.Pokemon_summary, error) {
	query := `
	SELECT pokemons.id, pokemons.name, pokemons.weight, pokemons.height, pokemons.sprite_url, json_group_array(pokemon_types.type_name)
	FROM pokemons
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
		&pokemon.SpriteUrl,
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
		err := rows.Scan(&pokemon.ID, &pokemon.Name, &pokemon.Weight, &pokemon.Height, &pokemon.SpriteUrl, &typesJSON)
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

// TODO: GetPokemonDetailed [name,id,height,weight,abilities,moves,evolution chain,games?]
func (s *sqliteDatabase) GetPokemonDetailed(ctx context.Context, id int) (model.Pokemon_details, error) {
	query := `
		WITH RECURSIVE full_chain AS (
            -- Start from the queried pokemon and go backwards to find the root
            SELECT pokemon_id, evolves_to_id, min_level, trigger_name, 0 as depth
            FROM evolution_chains
            WHERE evolves_to_id = ?
            
            UNION ALL
            
            -- Keep going backwards
            SELECT ec.pokemon_id, ec.evolves_to_id, ec.min_level, ec.trigger_name, fc.depth - 1
            FROM evolution_chains ec
            INNER JOIN full_chain fc ON ec.evolves_to_id = fc.pokemon_id
        ),
        root_pokemon AS (
            -- Find the root (pokemon with no prior evolution)
            SELECT COALESCE(
                (SELECT pokemon_id FROM full_chain ORDER BY depth LIMIT 1),
                ?
            ) as root_id
        ),
        complete_chain AS (
            -- Now traverse forward from the root to get all evolutions
            SELECT ec.pokemon_id, ec.evolves_to_id, ec.min_level, ec.trigger_name
            FROM evolution_chains ec, root_pokemon rp
            WHERE ec.pokemon_id = rp.root_id
            
            UNION ALL
            
            -- Recursively get all forward evolutions
            SELECT ec.pokemon_id, ec.evolves_to_id, ec.min_level, ec.trigger_name
            FROM evolution_chains ec
            INNER JOIN complete_chain cc ON ec.pokemon_id = cc.evolves_to_id
        )
        SELECT 
        pokemons.id,
        pokemons.name,
        pokemons.height,
        pokemons.weight,
		pokemons.sprite_url,
        (
            SELECT json_group_array(
                json_object(
                    'stat_name', pokemon_stats.stat_name,
                    'effort', pokemon_stats.effort,
                    'base_stat', pokemon_stats.base_stat
                )
            )
            FROM pokemon_stats
            WHERE pokemon_stats.pokemon_id = pokemons.id
        ) as stats,
        (
            SELECT json_group_array(pokemon_types.type_name)
            FROM pokemon_types
            WHERE pokemon_types.pokemon_id = pokemons.id
            ORDER BY pokemon_types.slot
        ) as types,
        (
            SELECT json_group_array(
                json_object(
                    'pokemon_id', cc.pokemon_id,
                    'pokemon_name', p1.name,
                    'evolves_to_id', cc.evolves_to_id,
                    'evolves_to_name', p2.name,
                    'min_level', COALESCE(cc.min_level, 0),
                    'trigger_name', cc.trigger_name
                )
            )
            FROM complete_chain cc
            JOIN pokemons p1 ON cc.pokemon_id = p1.id
            JOIN pokemons p2 ON cc.evolves_to_id = p2.id
        ) as evolution_chain
        FROM pokemons
        WHERE pokemons.id = ?
	`

	var pokemon model.Pokemon_details
	var statsJSON, typesJSON, evolutionJSON []byte

	err := s.db.QueryRowContext(ctx, query, id, id, id).Scan(
		&pokemon.ID,
		&pokemon.Name,
		&pokemon.Height,
		&pokemon.Weight,
		&pokemon.SpriteUrl,
		&statsJSON,
		&typesJSON,
		&evolutionJSON,
	)
	if err == sql.ErrNoRows {
		return model.Pokemon_details{}, sql.ErrNoRows
	}
	if err != nil {
		return model.Pokemon_details{}, err
	}

	json.Unmarshal(statsJSON, &pokemon.Stats)
	json.Unmarshal(typesJSON, &pokemon.Types)
	json.Unmarshal(evolutionJSON, &pokemon.EvolutionChain)

	return pokemon, nil
}

// TODO: Evolution chains

func (s *sqliteDatabase) AddEvolutionChain(ctx context.Context, chain model.Evolution_chain) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	query := `
	INSERT OR IGNORE INTO evolution_chains (pokemon_id, evolves_to_id, min_level, trigger_name) VALUES (?, ?, ?, ?)
	`

	stmtChain, err := tx.PrepareContext(ctx, query)
	if err != nil {
		return err
	}
	defer stmtChain.Close()

	// Recursive function to process chain links
	var processChainLink func(link model.ChainLink) error
	processChainLink = func(link model.ChainLink) error {
		// Get pokemon ID from species URL (extract ID from URL like ".../pokemon-species/21/")
		fromID := extractIDFromURL(link.Species.URL)

		// Process each evolution
		for _, evolvesTo := range link.EvolvesTo {
			toID := extractIDFromURL(evolvesTo.Species.URL)

			// Get evolution details (use first one if multiple exist)
			var minLevel *int
			var triggerName string

			if len(evolvesTo.EvolutionDetails) > 0 {
				detail := evolvesTo.EvolutionDetails[0]
				minLevel = detail.MinLevel
				triggerName = detail.Trigger.Name
			}

			// Insert evolution link
			_, err := stmtChain.ExecContext(ctx, fromID, toID, minLevel, triggerName)
			if err != nil {
				return err
			}

			// Recursively process next evolution stage
			if err := processChainLink(evolvesTo); err != nil {
				return err
			}
		}

		return nil
	}

	// Start processing from the root of the chain
	if err := processChainLink(chain.Chain); err != nil {
		return err
	}

	return tx.Commit()
}

// TODO: Effect entries for moves?
// Ability descriptions?

func (s *sqliteDatabase) InitDB() error {
	query := `
	PRAGMA foreign_keys = ON;
	CREATE TABLE IF NOT EXISTS pokemons (
	id INTEGER PRIMARY KEY,
	name TEXT UNIQUE NOT NULL,
	weight INTEGER,
	height INTEGER,
	sprite_url
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

	CREATE TABLE IF NOT EXISTS evolution_chains (
	pokemon_id INTEGER NOT NULL,
	evolves_to_id INTEGER NOT NULL,
	min_level INTEGER,
	trigger_name TEXT,

	PRIMARY KEY (pokemon_id, evolves_to_id),
	FOREIGN KEY (pokemon_id) REFERENCES pokemons(id),
	FOREIGN KEY (evolves_to_id) REFERENCES pokemons(id)
	);
	`

	_, err := s.db.Exec(query)

	return err
}

func (s *sqliteDatabase) Close() error {
	return s.db.Close()
}

// Helper function for extracting pokemon ID from pokeapi url
func extractIDFromURL(url string) int {
	parts := strings.Split(strings.TrimSuffix(url, "/"), "/")
	id, _ := strconv.Atoi(parts[len(parts)-1])
	return id
}
