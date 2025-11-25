package postgres

import (
	"context"
	"testovoeAvito/internal/domain"
)

func (r *Repository) GetGlobalStats(ctx context.Context) (*domain.Stats, error) {
	stats := &domain.Stats{}

	query := `
		SELECT
			(SELECT COUNT(*) FROM teams) as total_teams,
			(SELECT COUNT(*) FROM users) as total_users,
			(SELECT COUNT(*) FROM users WHERE is_active = true) as active_users,
			(SELECT COUNT(*) FROM pull_requests) as total_prs,
			(SELECT COUNT(*) FROM pull_requests WHERE status = 'OPEN') as open_prs,
			(SELECT COUNT(*) FROM pull_requests WHERE status = 'MERGED') as merged_prs,
			(SELECT COALESCE(AVG(cnt), 0) FROM (
				SELECT COUNT(*) as cnt FROM pr_reviewers GROUP BY pr_id
			 ) sub) as avg_reviewers
	`

	err := r.pool.QueryRow(ctx, query).Scan(
		&stats.TotalTeams,
		&stats.TotalUsers,
		&stats.ActiveUsers,
		&stats.TotalPRs,
		&stats.OpenPRs,
		&stats.MergedPRs,
		&stats.AverageReviewers,
	)

	if err != nil {
		return nil, mapError(err)
	}

	return stats, nil
}
