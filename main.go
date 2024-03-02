package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"

	"cloud.google.com/go/cloudsqlconn"
	"github.com/go-sql-driver/mysql"
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

type Ability struct {
	ID     int    `json:"id"`
	Name   string `json:"name"`
	NameJa string `json:"name_ja"`
}

var allowedOrigins = []string{
	"http://localhost:5173",
	"https://one-milli.github.io/pokemon-builder/",
}

func main() {
	dbUser := os.Getenv("DB_USERNAME")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	instanceConnectionName := os.Getenv("INSTANCE_CONNECTION_NAME")

	d, err := cloudsqlconn.NewDialer(context.Background())
	if err != nil {
		log.Fatalf("cloudsqlconn.NewDialer: %v", err)
	}
	mysql.RegisterDialContext("cloudsqlconn", func(ctx context.Context, addr string) (net.Conn, error) {
		return d.Dial(ctx, instanceConnectionName)
	})

	dsn := fmt.Sprintf("%s:%s@cloudsqlconn/%s?parseTime=true", dbUser, dbPassword, dbName)
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("sql.Open: %v", err)
	}
	defer db.Close()

	r := mux.NewRouter()

	// root
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "Pokemon Builder API")
	})

	// health check
	r.Handle("/health", enableCORS(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := map[string]string{"status": "API is up and running"}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	})))

	r.Handle("/user-pokemons/{userId}", enableCORS(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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
	}))).Methods("GET")

	r.Handle("/abilities", enableCORS(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rows, err := db.Query("SELECT id, name, name_ja FROM abilities")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var abilities []Ability
		for rows.Next() {
			var ability Ability
			if err := rows.Scan(&ability.ID, &ability.Name, &ability.NameJa); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			abilities = append(abilities, ability)
		}

		if err := rows.Err(); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(abilities)
	}))).Methods("GET")

	fmt.Println("Server is running at http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}

func enableCORS(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")

		for _, allowedOrigin := range allowedOrigins {
			if origin == allowedOrigin {
				w.Header().Set("Access-Control-Allow-Origin", origin)
				break
			}
		}

		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == "OPTIONS" {
			return
		}

		next.ServeHTTP(w, r)
	})
}
