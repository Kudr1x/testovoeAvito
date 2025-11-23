package main

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/golang-migrate/migrate/v4"

	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"

	"log"
	"testovoeAvito/internal/domain"
	"testovoeAvito/internal/repository/postgres"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	dbURL = "postgres://user:password@localhost:5435/reviewdb?sslmode=disable"
)

func main() {
	if err := runMigrations(); err != nil {
		log.Fatalln(err)
	}

	ctx := context.Background()
	pool, err := pgxpool.New(ctx, dbURL)
	if err != nil {
		log.Fatalln(err)
	}
	defer pool.Close()

	repo := postgres.NewRepository(pool)

	team := &domain.Team{
		Name: "Backend-Go",
		Members: []*domain.User{
			{ID: "u1", Username: "Alice", IsActive: true},
			{ID: "u2", Username: "Bob", IsActive: true},
			{ID: "u3", Username: "Charlie", IsActive: false},
		},
	}

	err = repo.SaveTeam(ctx, team)
	if err != nil {
		log.Println(err)
	} else {
		log.Println("Team saved")
	}

	savedTeam, err := repo.GetTeamByName(ctx, "Backend-Go")
	if err != nil {
		log.Fatalln(err)
	}
	printJSON(savedTeam)

	pr := &domain.PullRequest{
		ID:        "pr-1",
		Name:      "Fix login bug",
		AuthorID:  "u1",
		Status:    domain.PRStatusOpen,
		Reviewers: []string{"u2"},
		CreatedAt: time.Now(),
	}

	err = repo.SavePR(ctx, pr)
	if err != nil {
		log.Println(err)
	} else {
		log.Println("ok")
	}

	savedPR, err := repo.GetPR(ctx, "pr-1")
	if err != nil {
		log.Fatalln(err)
	}
	printJSON(savedPR)
}

func runMigrations() error {
	m, err := migrate.New(
		"file://migrations",
		dbURL,
	)
	if err != nil {
		return err
	}
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return err
	}
	return nil
}

func printJSON(v interface{}) {
	b, _ := json.MarshalIndent(v, "", "  ")
	fmt.Println(string(b))
}
