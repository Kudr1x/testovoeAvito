package service

import (
	"math/rand"
	"testovoeAvito/internal/domain"
	"time"
)

type Manager struct {
	repo domain.Repository
	rnd  *rand.Rand
}

func NewManager(repo domain.Repository) *Manager {
	source := rand.NewSource(time.Now().UnixNano())
	return &Manager{
		repo: repo,
		rnd:  rand.New(source),
	}
}
