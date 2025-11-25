package tests

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
)

var testPool *pgxpool.Pool

const defaultTestDB = "postgres://user:password@localhost:5435/reviewdb?sslmode=disable"

func TestMain(m *testing.M) {
	dbURL := os.Getenv("TEST_DB_URL")
	if dbURL == "" {
		dbURL = defaultTestDB
	}

	var err error
	ctx := context.Background()
	testPool, err = pgxpool.New(ctx, dbURL)
	if err != nil {
		log.Fatalf("Could not connect to database: %v", err)
	}
	defer testPool.Close()

	if err := testPool.Ping(ctx); err != nil {
		log.Fatalf("Could not ping database: %v", err)
	}

	code := m.Run()

	os.Exit(code)
}

func cleanupDB(t *testing.T) {
	ctx := context.Background()
	queries := []string{
		"TRUNCATE TABLE pr_reviewers CASCADE",
		"TRUNCATE TABLE pull_requests CASCADE",
		"TRUNCATE TABLE users CASCADE",
		"TRUNCATE TABLE teams CASCADE",
	}

	for _, q := range queries {
		_, err := testPool.Exec(ctx, q)
		if err != nil {
			t.Fatalf("Failed to cleanup DB: %v", err)
		}
	}
}
