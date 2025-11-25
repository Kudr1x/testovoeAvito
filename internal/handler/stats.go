package handler

import (
	"net/http"
)

func (h *Handler) getStatistics(w http.ResponseWriter, r *http.Request) {
	stats, err := h.svc.GetStatistics(r.Context())
	if err != nil {
		respondWithError(w, err)
		return
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"statistics": stats,
	})
}
