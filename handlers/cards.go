package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/jackc/pgx/v5"
)

type Card struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type App struct {
	DB *pgx.Conn
}

func (app *App) CardsHandler(w http.ResponseWriter, r *http.Request) {
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
