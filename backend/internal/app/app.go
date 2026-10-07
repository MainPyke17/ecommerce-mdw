// Package app is the HTTP composition root: it mounts every module's routes
// exactly once on a single router. The long-running binary (cmd/api) and the
// Vercel function (api/index.go) both serve this same handler, so the two
// entry points can never expose different routes.
package app

import (
	"net/http"

	chimw "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"gorm.io/gorm"

	"github.com/Unknowns24/ecommerce-mdw/config"
	apierr "github.com/Unknowns24/ecommerce-mdw/internal/adapters/rest/errors"
	"github.com/Unknowns24/ecommerce-mdw/internal/adapters/rest/handlers/catalogo"
	"github.com/Unknowns24/ecommerce-mdw/internal/adapters/rest/handlers/identidad"
	"github.com/Unknowns24/ecommerce-mdw/internal/adapters/rest/handlers/inventario"
	"github.com/Unknowns24/ecommerce-mdw/internal/adapters/rest/handlers/pagos"
	"github.com/Unknowns24/ecommerce-mdw/internal/adapters/rest/handlers/pedidos"
	"github.com/Unknowns24/ecommerce-mdw/internal/adapters/rest/middleware"
	"github.com/Unknowns24/ecommerce-mdw/internal/adapters/rest/router"
)

// Handler builds the complete API router.
func Handler(cfg config.Config, db *gorm.DB, version string) http.Handler {
	mw := middleware.Nuevo(cfg.AppSecretKey)
	r := router.NewRouter()
	r.Use(chimw.RequestID, chimw.RealIP, chimw.Recoverer, chimw.Logger)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PATCH", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	salud := func(w http.ResponseWriter, r *http.Request) {
		apierr.JSON(w, http.StatusOK, map[string]string{"estado": "ok", "version": version})
	}
	r.Get("/api/salud", salud)
	// La raíz responde algo útil a quien abre la URL pública en el navegador:
	// esta entrega es sólo la API (el frontend llega en el Parcial II).
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		apierr.JSON(w, http.StatusOK, map[string]string{
			"servicio": "BC Importados — API",
			"estado":   "ok",
			"version":  version,
			"salud":    "/api/salud",
		})
	})

	identidad.Montar(r, db, cfg, mw)
	catalogo.Montar(r, db, cfg, mw)
	inventario.Montar(r, db, cfg, mw)
	pagos.Montar(r, db, cfg, mw)
	pedidos.Montar(r, db, cfg, mw)

	return r
}
