package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jamiiiee/openpoketcg-api/internal/models"
)

func CardsHandler(app *models.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rows, err := app.DB.Query(r.Context(), "SELECT id, name FROM cards LIMIT 10")
		if err != nil {
			http.Error(w, "Failed to query cards", http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var cards []models.Card
		for rows.Next() {
			var card models.Card
			if err := rows.Scan(&card.ID, &card.Name); err != nil {
				http.Error(w, "Failed to scan card", http.StatusInternalServerError)
				return
			}
			cards = append(cards, card)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(cards)
	}
}

func CardIDHandler(app *models.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		parts := strings.Split(r.URL.Path, "/")
		if len(parts) < 4 {
			http.Error(w, "Missing card ID in path", http.StatusBadRequest)
			return
		}

		idStr := parts[3]
		id, err := strconv.Atoi(idStr)
		if err != nil {
			http.Error(w, "Invalid card ID", http.StatusBadRequest)
			return
		}

		var card models.Card
		query := "SELECT id, name FROM cards WHERE id = $1"
		err = app.DB.QueryRow(r.Context(), query, id).Scan(&card.ID, &card.Name)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				http.Error(w, "Card not found", http.StatusNotFound)
			} else {
				http.Error(w, "Database error", http.StatusInternalServerError)
			}
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(card)
	}
}
