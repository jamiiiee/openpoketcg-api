package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/jackc/pgx/v5"
	"github.com/jamiiiee/openpoketcg-api/internal/handlers"
	"github.com/jamiiiee/openpoketcg-api/internal/middleware"
	"github.com/joho/godotenv"
)

type App struct {
	DB *pgx.Conn
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

	app := &handlers.App{DB: conn}
	http.HandleFunc("/v0/cards", middleware.WithCacheControl(3600, middleware.WithETag(app.CardsHandler)))

	port := "8080"
	log.Printf("Server is running on port %s", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
