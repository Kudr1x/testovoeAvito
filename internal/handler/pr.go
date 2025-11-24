package handler

import (
	"encoding/json"
	"net/http"
)

type createPRRequest struct {
	ID       string `json:"pull_request_id"`
	Name     string `json:"pull_request_name"`
	AuthorID string `json:"author_id"`
}

func (h *Handler) createPR(w http.ResponseWriter, r *http.Request) {
	var req createPRRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	pr, err := h.svc.CreatePR(r.Context(), req.ID, req.Name, req.AuthorID)
	if err != nil {
		respondWithError(w, err)
		return
	}

	respondJSON(w, http.StatusCreated, map[string]interface{}{
		"pr": pr,
	})
}

type mergePRRequest struct {
	ID string `json:"pull_request_id"`
}

func (h *Handler) mergePR(w http.ResponseWriter, r *http.Request) {
	var req mergePRRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	pr, err := h.svc.MergePR(r.Context(), req.ID)
	if err != nil {
		respondWithError(w, err)
		return
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"pr": pr,
	})
}

type reassignRequest struct {
	ID        string `json:"pull_request_id"`
	OldUserID string `json:"old_user_id"`
}

func (h *Handler) reassignReviewer(w http.ResponseWriter, r *http.Request) {
	var req reassignRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	pr, newReviewer, err := h.svc.ReassignReviewer(r.Context(), req.ID, req.OldUserID)
	if err != nil {
		respondWithError(w, err)
		return
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"pr":          pr,
		"replaced_by": newReviewer,
	})
}
