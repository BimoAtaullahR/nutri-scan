package httpapi

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"time"

	"github.com/BimoAtaullahR/nutri-scan/services/backend/internal/nudge"
	"github.com/BimoAtaullahR/nutri-scan/services/backend/internal/platform/config"
	"github.com/BimoAtaullahR/nutri-scan/services/backend/internal/platform/database"
	"github.com/BimoAtaullahR/nutri-scan/services/backend/internal/scan"
	"github.com/BimoAtaullahR/nutri-scan/services/backend/internal/trend"
	"github.com/BimoAtaullahR/nutri-scan/services/backend/internal/user"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func NewRouter(cfg config.Config, db *database.DB, logger *slog.Logger) http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Recoverer)
	r.Use(middleware.CleanPath)
	r.Use(middleware.Timeout(30 * time.Second))

	r.Get("/healthz", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusOK, map[string]string{
			"status": "ok",
		})
	})

	user.NewHandler(user.NewPostgresStore(db.Pool), logger).RegisterRoutes(r)
	scan.NewHandler(cfg.AIInferenceURL, logger).RegisterRoutes(r)
	trend.NewHandler(logger).RegisterRoutes(r)
	nudge.NewHandler(logger).RegisterRoutes(r)

	return r
}

func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}
