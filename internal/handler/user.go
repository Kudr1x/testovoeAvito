package handler

import (
	"encoding/json"
	"net/http"
)

type setActiveRequest struct {
	UserID   string `json:"user_id"`
	IsActive bool   `json:"is_active"`
}

func (h *Handler) setUserActive(w http.ResponseWriter, r *http.Request) {
	var req setActiveRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	user, err := h.svc.SetUserActive(r.Context(), req.UserID, req.IsActive)
	if err != nil {
		respondWithError(w, err)
		return
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"user": user,
	})
}

func (h *Handler) getUserReviews(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("user_id")
	if userID == "" {
		http.Error(w, "missing user_id", http.StatusBadRequest)
		return
	}

	prs, err := h.svc.GetUserPendingReviews(r.Context(), userID)
	if err != nil {
		respondWithError(w, err)
		return
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"user_id":       userID,
		"pull_requests": prs,
	})
}
