// Package handler is the Vercel entry point: Vercel's Go runtime serves the
// exported Handler as a serverless function, and vercel.json rewrites every
// path to it. It serves the same router as cmd/api (internal/app).
//
// Migrations do not run here: they run in the build step
// (scripts/vercel-build.sh), before the new deployment receives traffic,
// the same order the VPS pipeline uses (ADR-002, ADR-005).
package handler

import (
	"log"
	"net/http"
	"sync"

	"github.com/Unknowns24/ecommerce-mdw/config"
	"github.com/Unknowns24/ecommerce-mdw/internal/adapters/db/postgresql"
	apierr "github.com/Unknowns24/ecommerce-mdw/internal/adapters/rest/errors"
	"github.com/Unknowns24/ecommerce-mdw/internal/app"
)

var (
	mu  sync.Mutex
	api http.Handler
)

// Handler builds the router on the first request of each instance and reuses
// it (and its connection pool) afterwards. If the configuration or the
// database fails, it answers 503 and retries on the next request instead of
// caching the failure for the life of the instance.
func Handler(w http.ResponseWriter, r *http.Request) {
	h, err := obtener()
	if err != nil {
		log.Printf("api: no se pudo iniciar: %v", err)
		apierr.JSON(w, http.StatusServiceUnavailable, apierr.Cuerpo{Error: "Servicio no disponible, probá de nuevo en unos segundos"})
		return
	}
	h.ServeHTTP(w, r)
}

func obtener() (http.Handler, error) {
	mu.Lock()
	defer mu.Unlock()
	if api != nil {
		return api, nil
	}
	cfg, err := config.Load()
	if err != nil {
		return nil, err
	}
	db, err := postgresql.Abrir(cfg.DatabaseDSN)
	if err != nil {
		return nil, err
	}
	api = app.Handler(cfg, db, "vercel")
	return api, nil
}
