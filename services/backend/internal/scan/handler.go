package scan

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type Handler struct {
	aiInferenceURL string
	logger         *slog.Logger
}

func NewHandler(aiInferenceURL string, logger *slog.Logger) *Handler {
	return &Handler{
		aiInferenceURL: aiInferenceURL,
		logger:         logger,
	}
}

func (h *Handler) RegisterRoutes(r chi.Router) {
	r.Post("/scans", h.createScan)
	r.Get("/scans/{scanID}", h.getScan)
}

func (h *Handler) createScan(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusNotImplemented, map[string]string{
		"error":   "not_implemented",
		"message": "scan creation is not implemented yet",
	})
}

func (h *Handler) getScan(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusNotImplemented, map[string]string{
		"error":   "not_implemented",
		"message": "scan retrieval is not implemented yet",
		"scanId":  chi.URLParam(r, "scanID"),
	})
}

func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}
