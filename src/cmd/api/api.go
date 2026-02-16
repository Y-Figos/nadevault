package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	nadevault "github.com/Y-Figos/nadevault/internal/db"
	internalhttp "github.com/Y-Figos/nadevault/internal/http"
	"github.com/Y-Figos/nadevault/internal/http/web"
	"github.com/Y-Figos/nadevault/internal/repository"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

func main() {

	godotenv.Load("../.env")
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	DATABASE_URL := os.Getenv("DATABASE_URL")
	log.Printf("Connecting to database at %s", DATABASE_URL)
	pool, err := pgxpool.New(ctx, DATABASE_URL)
	if err != nil {
		log.Printf("Failed to connect to database: %v", err)
		panic(err)
	}
	defer pool.Close()
	queries := nadevault.New(pool)

	repo := repository.NewPostgresNadeRepository(queries)
	renderer, err := web.NewRenderer("internal/http/web/templates/*.html")
	if err != nil{
		log.Printf("Failed to Render Glob: %v", err)
		panic(err)
	}
	router := internalhttp.NewRouter(repo, renderer)

	srv := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}
	go func() {
		log.Println("Starting server on port 8080")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Printf("Failed to start server: %v", err)
		}
	}()
	<-ctx.Done()

	log.Println("Shutting down server...")
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Printf("Failed to shutdown server gracefully: %v", err)
	}
	log.Println("Server stopped gracefully")
}
