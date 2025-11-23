package domain

import "errors"

var (
	ErrNotFound      = errors.New("resource not found")
	ErrTeamExists    = errors.New("team already exists")
	ErrPRExists      = errors.New("pull request already exists")
	ErrPRMerged      = errors.New("cannot edit merged pull request")
	ErrNotAssigned   = errors.New("reviewer is not assigned to this PR")
	ErrNoCandidate   = errors.New("no active replacement candidate in team")
	ErrInvalidParams = errors.New("invalid input parameters")
)
