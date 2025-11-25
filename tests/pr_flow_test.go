package tests

import (
	"context"
	"testing"
	"testovoeAvito/internal/domain"
	"testovoeAvito/internal/repository/postgres"
	"testovoeAvito/internal/service"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreatePR_AssignsReviewers(t *testing.T) {
	cleanupDB(t)
	ctx := context.Background()

	repo := postgres.NewRepository(testPool)
	svc := service.NewManager(repo)

	teamName := "backend"
	users := []*domain.User{
		{ID: "u1", Username: "Author", IsActive: true},
		{ID: "u2", Username: "Reviewer1", IsActive: true},
		{ID: "u3", Username: "Reviewer2", IsActive: true},
		{ID: "u4", Username: "Reviewer3", IsActive: false},
	}
	team := &domain.Team{Name: teamName, Members: users}

	_, err := svc.CreateTeam(ctx, team)
	require.NoError(t, err)

	prID := "pr-1"
	pr, err := svc.CreatePR(ctx, prID, "Feature X", "u1")
	require.NoError(t, err)

	assert.Equal(t, prID, pr.ID)
	assert.Equal(t, domain.PRStatusOpen, pr.Status)
	assert.Equal(t, "u1", pr.AuthorID)

	assert.Len(t, pr.Reviewers, 2)
	assert.NotContains(t, pr.Reviewers, "u1", "Author should not be a reviewer")
	assert.NotContains(t, pr.Reviewers, "u4", "Inactive user should not be a reviewer")
	assert.Contains(t, pr.Reviewers, "u2")
	assert.Contains(t, pr.Reviewers, "u3")
}
