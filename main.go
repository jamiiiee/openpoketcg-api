package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/jackc/pgx/v5"
	"github.com/jamiiiee/openpoketcg-api/middleware"
	"github.com/joho/godotenv"
)

type App struct {
	DB *pgx.Conn
}

type Card struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	conn, err := pgx.Connect(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}
	defer conn.Close(context.Background())

	app := &App{DB: conn}

	http.HandleFunc("/cards", middleware.WithCacheControl(3600, middleware.WithETag(app.cardsHandler)))

	port := "8080"
	log.Printf("Server is running on port %s", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func (app *App) cardsHandler(w http.ResponseWriter, r *http.Request) {
	rows, err := app.DB.Query(r.Context(), "SELECT id, name FROM cards LIMIT 10")
	if err != nil {
		http.Error(w, "Failed to query cards", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var cards []Card
	for rows.Next() {
		var card Card
		if err := rows.Scan(&card.ID, &card.Name); err != nil {
			http.Error(w, "Failed to scan card", http.StatusInternalServerError)
			return
		}
		cards = append(cards, card)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(cards)
}
