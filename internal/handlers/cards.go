package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/jamiiiee/openpoketcg-api/internal/models"
)

func CardsHandler(app *models.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rows, err := app.DB.Query(r.Context(), "SELECT id, name FROM cards LIMIT 10")
		if err != nil {
			http.Error(w, "Failed to query cards", http.StatusInternalServerError)
			fmt.Print(err)
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
