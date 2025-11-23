package postgres

import (
	"context"
	"testovoeAvito/internal/domain"
)

func (r *Repository) UpdateUserActivity(ctx context.Context, userID string, isActive bool) (*domain.User, error) {
	var u domain.User
	err := r.pool.QueryRow(ctx, `
		UPDATE users 
		SET is_active = $2 
		WHERE id = $1 
		RETURNING id, username, team_name, is_active
	`, userID, isActive).Scan(&u.ID, &u.Username, &u.TeamName, &u.IsActive)

	if err != nil {
		return nil, mapError(err)
	}
	return &u, nil
}
