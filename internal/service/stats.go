package service

import (
	"context"
	"fmt"
	"testovoeAvito/internal/domain"
)

func (s *Manager) GetStatistics(ctx context.Context) (*domain.Stats, error) {
	stats, err := s.repo.GetGlobalStats(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to collect statistics: %w", err)
	}
	return stats, nil
}
