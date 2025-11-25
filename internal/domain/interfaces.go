package domain

import "context"

type TeamRepository interface {
	SaveTeam(ctx context.Context, team *Team) error
	GetTeamByName(ctx context.Context, name string) (*Team, error)
}

type UserRepository interface {
	UpdateUserActivity(ctx context.Context, userID string, isActive bool) (*User, error)
	GetUser(ctx context.Context, userID string) (*User, error)
}

type PullRequestRepository interface {
	SavePR(ctx context.Context, pr *PullRequest) error
	GetPR(ctx context.Context, id string) (*PullRequest, error)
	GetUserPRs(ctx context.Context, userID string) ([]*PullRequest, error)
}

type Repository interface {
	TeamRepository
	UserRepository
	PullRequestRepository
	StatsRepository
}

type TeamService interface {
	CreateTeam(ctx context.Context, team *Team) (*Team, error)
	GetTeam(ctx context.Context, name string) (*Team, error)
}

type UserService interface {
	SetUserActive(ctx context.Context, userID string, isActive bool) (*User, error)
	GetUserPendingReviews(ctx context.Context, userID string) ([]*PullRequest, error)
}

type PullRequestService interface {
	CreatePR(ctx context.Context, prID, prName, authorID string) (*PullRequest, error)
	MergePR(ctx context.Context, prID string) (*PullRequest, error)
	ReassignReviewer(ctx context.Context, prID, oldReviewerID string) (*PullRequest, string, error)
}

type Service interface {
	TeamService
	UserService
	PullRequestService
	StatsService
}

type StatsRepository interface {
	GetGlobalStats(ctx context.Context) (*Stats, error)
}

type StatsService interface {
	GetStatistics(ctx context.Context) (*Stats, error)
}
