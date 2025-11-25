package domain

import "time"

type PRStatus string

const (
	PRStatusOpen   PRStatus = "OPEN"
	PRStatusMerged PRStatus = "MERGED"
)

type User struct {
	ID       string `json:"user_id"`
	Username string `json:"username"`
	TeamName string `json:"team_name,omitempty"`
	IsActive bool   `json:"is_active"`
}

type Team struct {
	Name    string  `json:"team_name"`
	Members []*User `json:"members"`
}

type PullRequest struct {
	ID        string     `json:"pull_request_id"`
	Name      string     `json:"pull_request_name"`
	AuthorID  string     `json:"author_id"`
	Status    PRStatus   `json:"status"`
	Reviewers []string   `json:"assigned_reviewers"`
	CreatedAt time.Time  `json:"created_at"`
	MergedAt  *time.Time `json:"merged_at"`
}

type Stats struct {
	TotalTeams       int     `json:"total_teams"`
	TotalUsers       int     `json:"total_users"`
	ActiveUsers      int     `json:"active_users"`
	TotalPRs         int     `json:"total_prs"`
	OpenPRs          int     `json:"open_prs"`
	MergedPRs        int     `json:"merged_prs"`
	AverageReviewers float64 `json:"avg_reviewers_per_pr"`
}
