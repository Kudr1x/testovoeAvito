package postgres

import (
	"context"
	"errors"
	"testovoeAvito/internal/domain"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
)

func (r *Repository) SaveTeam(ctx context.Context, team *domain.Team) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	_, err = tx.Exec(ctx, "INSERT INTO teams (name) VALUES ($1)", team.Name)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
			return domain.ErrTeamExists
		}
		return err
	}

	for _, member := range team.Members {
		_, err = tx.Exec(ctx, `
			INSERT INTO users (id, username, team_name, is_active)
			VALUES ($1, $2, $3, $4)
			ON CONFLICT (id) DO UPDATE
			SET username = EXCLUDED.username, 
			    team_name = EXCLUDED.team_name, 
			    is_active = EXCLUDED.is_active
		`, member.ID, member.Username, team.Name, member.IsActive)
		if err != nil {
			return err
		}
	}

	return tx.Commit(ctx)
}

func (r *Repository) GetTeamByName(ctx context.Context, name string) (*domain.Team, error) {
	var teamName string
	err := r.pool.QueryRow(ctx, "SELECT name FROM teams WHERE name = $1", name).Scan(&teamName)
	if err != nil {
		return nil, mapError(err)
	}

	rows, err := r.pool.Query(ctx, "SELECT id, username, team_name, is_active FROM users WHERE team_name = $1", name)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var members []*domain.User
	for rows.Next() {
		u := &domain.User{}
		if err := rows.Scan(&u.ID, &u.Username, &u.TeamName, &u.IsActive); err != nil {
			return nil, err
		}
		members = append(members, u)
	}

	return &domain.Team{
		Name:    teamName,
		Members: members,
	}, nil
}
