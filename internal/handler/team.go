package handler

import (
	"encoding/json"
	"net/http"
	"testovoeAvito/internal/domain"
)

type createTeamRequest struct {
	TeamName string         `json:"team_name"`
	Members  []*domain.User `json:"members"`
}

func (h *Handler) createTeam(w http.ResponseWriter, r *http.Request) {
	var req createTeamRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, domain.ErrInvalidParams)
		return
	}

	team := &domain.Team{
		Name:    req.TeamName,
		Members: req.Members,
	}

	if _, err := h.svc.CreateTeam(r.Context(), team); err != nil {
		respondWithError(w, err)
		return
	}

	respondJSON(w, http.StatusCreated, map[string]interface{}{
		"team": team,
	})
}

func (h *Handler) getTeam(w http.ResponseWriter, r *http.Request) {
	teamName := r.URL.Query().Get("team_name")
	if teamName == "" {
		respondWithError(w, domain.ErrInvalidParams)
		return
	}

	team, err := h.svc.GetTeam(r.Context(), teamName)
	if err != nil {
		respondWithError(w, err)
		return
	}

	respondJSON(w, http.StatusOK, team)
}
