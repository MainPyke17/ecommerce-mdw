// Package catalogo mounts the public and administrative catalog endpoints.
// TODO(Genaro): reemplazar con las rutas reales (catálogo público, CRUD de
// producto/variante, categorías, marcas, imágenes).
package catalogo

import (
	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"

	"github.com/Unknowns24/ecommerce-mdw/config"
	"github.com/Unknowns24/ecommerce-mdw/internal/adapters/rest/middleware"
)

// Montar attaches this module's routes to the shared router. main.go calls
// it exactly once.
func Montar(r chi.Router, db *gorm.DB, cfg config.Config, mw *middleware.Middleware) {}
