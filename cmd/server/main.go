package main

import (
	"context"
	"fmt"
	"log"
	"testovoeAvito/internal/domain"
	"testovoeAvito/internal/repository/postgres"
	"testovoeAvito/internal/service"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5/pgxpool"
)

const dbURL = "postgres://user:password@localhost:5435/reviewdb?sslmode=disable"

func main() {
	ctx := context.Background()

	m, _ := migrate.New("file://migrations", dbURL)
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatal(err)
	}

	pool, err := pgxpool.New(ctx, dbURL)
	if err != nil {
		log.Fatal(err)
	}

	repo := postgres.NewRepository(pool)
	svc := service.NewManager(repo)

	repo.SaveTeam(ctx, &domain.Team{
		Name: "GoTeam",
		Members: []*domain.User{
			{ID: "u1", Username: "Alice", IsActive: true, TeamName: "GoTeam"},
			{ID: "u2", Username: "Bob", IsActive: true, TeamName: "GoTeam"},
			{ID: "u3", Username: "Charlie", IsActive: true, TeamName: "GoTeam"},
			{ID: "u4", Username: "Dave", IsActive: true, TeamName: "GoTeam"},
		},
	})

	pr, err := svc.CreatePR(ctx, "pr-final", "Fix stuff", "u1")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("OK. Reviewers: %v\n", pr.Reviewers)

	if len(pr.Reviewers) > 0 {
		oldRev := pr.Reviewers[0]
		fmt.Printf("2. Reassigning %s... ", oldRev)

		_, newRev, err := svc.ReassignReviewer(ctx, pr.ID, oldRev)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("OK. New reviewer: %s\n", newRev)
	}
}
