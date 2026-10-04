// Package inventario mounts the stock administration endpoints: providers,
// lot intake, availability, movement history and adjustments.
package inventario

import (
	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"

	"github.com/Unknowns24/ecommerce-mdw/config"
	"github.com/Unknowns24/ecommerce-mdw/internal/adapters/rest/middleware"
)

// Montar attaches this module's routes to the shared router. main.go calls
// it exactly once. Lo completa Agustín en el Bloque 1.
func Montar(r chi.Router, db *gorm.DB, cfg config.Config, mw *middleware.Middleware) {}
