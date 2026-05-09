package user

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type Handler struct {
	logger *slog.Logger
}

func NewHandler(logger *slog.Logger) *Handler {
	return &Handler{logger: logger}
}

func (h *Handler) RegisterRoutes(r chi.Router) {
	r.Post("/anonymous-users", h.createAnonymousUser)
	r.Get("/me/profile", h.getProfile)
	r.Put("/me/profile", h.updateProfile)
}

func (h *Handler) createAnonymousUser(w http.ResponseWriter, r *http.Request) {
	userID := "anon_" + randomHex(16)
	token := "nutriscan_" + randomHex(32)

	writeJSON(w, http.StatusCreated, map[string]string{
		"anonymousUserId": userID,
		"accessToken":     token,
		"tokenType":       "Bearer",
	})
}

func (h *Handler) getProfile(w http.ResponseWriter, r *http.Request) {
	writePlaceholder(w, "profile retrieval is not implemented yet")
}

func (h *Handler) updateProfile(w http.ResponseWriter, r *http.Request) {
	writePlaceholder(w, "profile update is not implemented yet")
}

func randomHex(size int) string {
	buf := make([]byte, size)
	if _, err := rand.Read(buf); err != nil {
		return "unavailable"
	}

	return hex.EncodeToString(buf)
}

func writePlaceholder(w http.ResponseWriter, message string) {
	writeJSON(w, http.StatusNotImplemented, map[string]string{
		"error":   "not_implemented",
		"message": message,
	})
}

func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}
