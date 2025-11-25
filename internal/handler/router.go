package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"testovoeAvito/internal/domain"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type Handler struct {
	svc domain.Service
}

func NewHandler(svc domain.Service) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) InitRoutes() http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json; charset=utf-8")
			next.ServeHTTP(w, r)
		})
	})

	r.Route("/team", func(r chi.Router) {
		r.Post("/add", h.createTeam)
		r.Get("/get", h.getTeam)
	})

	r.Route("/users", func(r chi.Router) {
		r.Post("/setIsActive", h.setUserActive)
		r.Get("/getReview", h.getUserReviews)
	})

	r.Route("/pullRequest", func(r chi.Router) {
		r.Post("/create", h.createPR)
		r.Post("/merge", h.mergePR)
		r.Post("/reassign", h.reassignReviewer)
	})

	r.Route("/stats", func(r chi.Router) {
		r.Get("/", h.getStatistics)
	})

	return r
}

func respondJSON(w http.ResponseWriter, status int, payload interface{}) {
	w.WriteHeader(status)
	if payload != nil {
		json.NewEncoder(w).Encode(payload)
	}
}

func respondWithError(w http.ResponseWriter, err error) {
	var statusCode int
	var errorCode string

	switch {
	case errors.Is(err, domain.ErrNotFound):
		statusCode = http.StatusNotFound
		errorCode = "NOT_FOUND"
	case errors.Is(err, domain.ErrTeamExists):
		statusCode = http.StatusBadRequest
		errorCode = "TEAM_EXISTS"
	case errors.Is(err, domain.ErrPRExists):
		statusCode = http.StatusConflict
		errorCode = "PR_EXISTS"
	case errors.Is(err, domain.ErrPRMerged):
		statusCode = http.StatusConflict
		errorCode = "PR_MERGED"
	case errors.Is(err, domain.ErrNotAssigned):
		statusCode = http.StatusConflict
		errorCode = "NOT_ASSIGNED"
	case errors.Is(err, domain.ErrNoCandidate):
		statusCode = http.StatusConflict
		errorCode = "NO_CANDIDATE"
	case errors.Is(err, domain.ErrInvalidParams):
		statusCode = http.StatusBadRequest
		errorCode = "INVALID_PARAMS"
	default:
		statusCode = http.StatusInternalServerError
		errorCode = "INTERNAL_ERROR"
	}

	respondJSON(w, statusCode, map[string]interface{}{
		"error": map[string]string{
			"code":    errorCode,
			"message": err.Error(),
		},
	})
}
