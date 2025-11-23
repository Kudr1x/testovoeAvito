package postgres

import (
	"context"
	"errors"
	"testovoeAvito/internal/domain"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
)

func (r *Repository) SavePR(ctx context.Context, pr *domain.PullRequest) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	_, err = tx.Exec(ctx, `
		INSERT INTO pull_requests (id, name, author_id, status, created_at, merged_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		ON CONFLICT (id) DO UPDATE 
		SET status = EXCLUDED.status, merged_at = EXCLUDED.merged_at
	`, pr.ID, pr.Name, pr.AuthorID, pr.Status, pr.CreatedAt, pr.MergedAt)

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
			return domain.ErrPRExists
		}
		return mapError(err)
	}

	_, err = tx.Exec(ctx, "DELETE FROM pr_reviewers WHERE pr_id = $1", pr.ID)
	if err != nil {
		return err
	}

	for _, revID := range pr.Reviewers {
		_, err = tx.Exec(ctx, "INSERT INTO pr_reviewers (pr_id, reviewer_id) VALUES ($1, $2)", pr.ID, revID)
		if err != nil {
			return err
		}
	}

	return tx.Commit(ctx)
}

func (r *Repository) GetPR(ctx context.Context, id string) (*domain.PullRequest, error) {
	pr := &domain.PullRequest{}

	err := r.pool.QueryRow(ctx, `
		SELECT id, name, author_id, status, created_at, merged_at 
		FROM pull_requests WHERE id = $1
	`, id).Scan(&pr.ID, &pr.Name, &pr.AuthorID, &pr.Status, &pr.CreatedAt, &pr.MergedAt)

	if err != nil {
		return nil, mapError(err)
	}

	rows, err := r.pool.Query(ctx, "SELECT reviewer_id FROM pr_reviewers WHERE pr_id = $1", id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var revID string
		if err := rows.Scan(&revID); err != nil {
			return nil, err
		}
		pr.Reviewers = append(pr.Reviewers, revID)
	}

	return pr, nil
}

func (r *Repository) GetUserPRs(ctx context.Context, userID string) ([]*domain.PullRequest, error) {
	query := `
		SELECT pr.id, pr.name, pr.author_id, pr.status 
		FROM pull_requests pr
		JOIN pr_reviewers rev ON pr.id = rev.pr_id
		WHERE rev.reviewer_id = $1
	`
	rows, err := r.pool.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var prs []*domain.PullRequest
	for rows.Next() {
		pr := &domain.PullRequest{}
		if err := rows.Scan(&pr.ID, &pr.Name, &pr.AuthorID, &pr.Status); err != nil {
			return nil, err
		}
		prs = append(prs, pr)
	}
	return prs, nil
}
