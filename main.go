package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

type UserPokemon struct {
	UserID           int            `json:"userId"`
	MyPokemonID      int            `json:"myPokemonId"`
	PokemonID        int            `json:"pid"`
	Name             string         `json:"name"`
	Level            int            `json:"level"`
	Nature           string         `json:"nature"`
	EffortValues     string         `json:"evs"`
	IndividualValues string         `json:"ivs"`
	Item             sql.NullString `json:"item"`
	Moves            string         `json:"moves"`
	AbilityID        int            `json:"abilityId"`
	Notes            sql.NullString `json:"notes"`
}

func main() {
	dbUser := os.Getenv("DB_USERNAME")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", dbUser, dbPassword, dbHost, dbPort, dbName)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	r := mux.NewRouter()

	// health check
	r.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		response := map[string]string{"status": "API is up and running"}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	})

	r.HandleFunc("/user-pokemon/{userId}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		userID := vars["userId"]

		var pokemons []UserPokemon
		query := "SELECT userId, myPokemonId, pid, name, level, nature, evs, ivs, item, moves, abilityId, notes FROM user_pokemons WHERE userId = ?"
		rows, err := db.Query(query, userID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		for rows.Next() {
			var p UserPokemon
			if err := rows.Scan(&p.UserID, &p.MyPokemonID, &p.PokemonID, &p.Name, &p.Level, &p.Nature, &p.EffortValues, &p.IndividualValues, &p.Item, &p.Moves, &p.AbilityID, &p.Notes); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			pokemons = append(pokemons, p)
		}

		if err := rows.Err(); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(pokemons)
	}).Methods("GET")

	fmt.Println("Server is running at http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
