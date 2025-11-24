package service

import (
	"context"
	"testovoeAvito/internal/domain"
)

func (s *Manager) CreateTeam(ctx context.Context, team *domain.Team) (*domain.Team, error) {
	if err := s.repo.SaveTeam(ctx, team); err != nil {
		return nil, err
	}
	return team, nil
}

func (s *Manager) GetTeam(ctx context.Context, name string) (*domain.Team, error) {
	return s.repo.GetTeamByName(ctx, name)
}
