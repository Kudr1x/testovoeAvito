package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"testovoeAvito/internal/handler"
	"testovoeAvito/internal/repository/postgres"
	"testovoeAvito/internal/service"
	"time"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5/pgxpool"
)

const serverPort = ":8080"

func main() {
	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		dbURL = "postgres://user:password@localhost:5435/reviewdb?sslmode=disable"
	}

	m, err := migrate.New("file://migrations", dbURL)
	if err != nil {
		log.Fatalf("Migration init failed: %v", err)
	}
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatalf("Migration failed: %v", err)
	}

	ctx := context.Background()
	pool, err := pgxpool.New(ctx, dbURL)
	if err != nil {
		log.Fatalf("DB connection failed: %v", err)
	}
	defer pool.Close()

	repo := postgres.NewRepository(pool)
	svc := service.NewManager(repo)
	h := handler.NewHandler(svc)

	router := h.InitRoutes()

	srv := &http.Server{
		Addr:         serverPort,
		Handler:      router,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	log.Printf("Server started on port %s", serverPort)
	if err := srv.ListenAndServe(); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
