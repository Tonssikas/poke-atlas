#  Poke Atlas

A full-stack web application for browsing and exploring Pokémon data, built with Go and Svelte.

---

# DISCLAIMER
***This is a hobby project and is heavily in progress***



##  Tech Stack

**Backend:**
- Go with Gin web framework
- SQLite database
- RESTful API design
- Integration with [PokéAPI](https://pokeapi.co/)

**Frontend:**
- Svelte
- TypeScript
- CSS

##  Features

- Browse all Pokémon with pagination
- View detailed Pokémon information
- Fast and responsive user interface
- Local caching with SQLite for improved performance

##  API Endpoints

- `GET /pokemon/:name` - Get Pokémon by name
- `GET /pokemons/:offset` - Get paginated list of Pokémon
- `GET /pokemondetailed/:id` - Get detailed Pokémon information

##  Running Locally

### Backend
```bash
cd backend
go run cmd/server/main.go
```

### Frontend
```bash
cd frontend
pnpm install
pnpm run dev
```


##  Goals of the project

  - To explore full-stack development with Go and Svelte.
  - To learn more about CI/CD pipelines and testing
  - To help me quickly fetch information about pokemons when playing Pokemon games
