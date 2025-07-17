package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jamiiiee/openpoketcg-api/internal/handlers"
	"github.com/jamiiiee/openpoketcg-api/internal/middleware"
	"github.com/jamiiiee/openpoketcg-api/internal/models"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	apiKey := os.Getenv("API_KEY")
	if apiKey == "" {
		log.Fatal("Missing required environment variable: API_KEY")
	}
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("Missing required environment variable: DATABASE_URL")
	}

	pool, err := pgxpool.New(context.Background(), dbURL)
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}
	app := &models.App{DB: pool}

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})

	http.HandleFunc("/v0/cards",
		middleware.WithCORS(
			middleware.WithCacheControl(3600,
				middleware.RequireAPIKey(apiKey,
					middleware.WithETag(handlers.CardsHandler(app))),
			),
		))

	http.HandleFunc("/v0/cards/",
		middleware.WithCORS(
			middleware.RequireAPIKey(apiKey,
				middleware.WithETag(handlers.CardIDHandler(app))),
		))

	port := "8080"
	srv := &http.Server{
		Addr:         ":" + port,
		Handler:      http.DefaultServeMux,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	go func() {
		log.Printf("Server is running on port %s", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed: %v", err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop

	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	pool.Close()
	log.Println("Server stopped cleanly")
}
