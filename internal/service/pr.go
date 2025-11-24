package service

import (
	"context"
	"fmt"
	"testovoeAvito/internal/domain"
	"time"
)

func (s *Manager) CreatePR(ctx context.Context, prID, prName, authorID string) (*domain.PullRequest, error) {
	author, err := s.repo.GetUser(ctx, authorID)
	if err != nil {
		return nil, fmt.Errorf("failed to get author: %w", err)
	}

	team, err := s.repo.GetTeamByName(ctx, author.TeamName)
	if err != nil {
		return nil, fmt.Errorf("failed to get team: %w", err)
	}

	reviewers := s.selectReviewers(team.Members, authorID)

	pr := &domain.PullRequest{
		ID:        prID,
		Name:      prName,
		AuthorID:  authorID,
		Status:    domain.PRStatusOpen,
		Reviewers: reviewers,
		CreatedAt: time.Now(),
	}

	if err := s.repo.SavePR(ctx, pr); err != nil {
		return nil, err
	}

	return pr, nil
}

func (s *Manager) selectReviewers(members []*domain.User, authorID string) []string {
	candidates := make([]string, 0, len(members))

	for _, m := range members {
		if m.IsActive && m.ID != authorID {
			candidates = append(candidates, m.ID)
		}
	}

	s.rnd.Shuffle(len(candidates), func(i, j int) {
		candidates[i], candidates[j] = candidates[j], candidates[i]
	})

	limit := 2
	if len(candidates) < limit {
		limit = len(candidates)
	}

	return candidates[:limit]
}

func (s *Manager) MergePR(ctx context.Context, prID string) (*domain.PullRequest, error) {
	pr, err := s.repo.GetPR(ctx, prID)
	if err != nil {
		return nil, err
	}

	if pr.Status == domain.PRStatusMerged {
		return pr, nil
	}

	now := time.Now()
	pr.Status = domain.PRStatusMerged
	pr.MergedAt = &now

	if err := s.repo.SavePR(ctx, pr); err != nil {
		return nil, err
	}

	return pr, nil
}

func (s *Manager) ReassignReviewer(ctx context.Context, prID, oldReviewerID string) (*domain.PullRequest, string, error) {
	pr, err := s.repo.GetPR(ctx, prID)
	if err != nil {
		return nil, "", fmt.Errorf("failed to get PR: %w", err)
	}

	if pr.Status == domain.PRStatusMerged {
		return nil, "", domain.ErrPRMerged
	}

	reviewerIndex := -1
	for i, id := range pr.Reviewers {
		if id == oldReviewerID {
			reviewerIndex = i
			break
		}
	}
	if reviewerIndex == -1 {
		return nil, "", domain.ErrNotAssigned
	}

	oldReviewer, err := s.repo.GetUser(ctx, oldReviewerID)
	if err != nil {
		return nil, "", fmt.Errorf("failed to get old reviewer info: %w", err)
	}

	team, err := s.repo.GetTeamByName(ctx, oldReviewer.TeamName)
	if err != nil {
		return nil, "", fmt.Errorf("failed to get team: %w", err)
	}

	candidates := make([]string, 0)
	for _, member := range team.Members {
		if !member.IsActive {
			continue
		}
		if member.ID == pr.AuthorID {
			continue
		}
		if member.ID == oldReviewerID {
			continue
		}
		isAlreadyAssigned := false
		for _, existingID := range pr.Reviewers {
			if existingID == member.ID {
				isAlreadyAssigned = true
				break
			}
		}
		if isAlreadyAssigned {
			continue
		}

		candidates = append(candidates, member.ID)
	}

	if len(candidates) == 0 {
		return nil, "", domain.ErrNoCandidate
	}

	randIndex := s.rnd.Intn(len(candidates))
	newReviewerID := candidates[randIndex]

	pr.Reviewers[reviewerIndex] = newReviewerID

	if err := s.repo.SavePR(ctx, pr); err != nil {
		return nil, "", err
	}

	return pr, newReviewerID, nil
}
