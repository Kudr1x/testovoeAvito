package service

import (
	"context"
	"testovoeAvito/internal/domain"
)

func (s *Manager) SetUserActive(ctx context.Context, userID string, isActive bool) (*domain.User, error) {
	return s.repo.UpdateUserActivity(ctx, userID, isActive)
}

func (s *Manager) GetUserPendingReviews(ctx context.Context, userID string) ([]*domain.PullRequest, error) {
	return s.repo.GetUserPRs(ctx, userID)
}
