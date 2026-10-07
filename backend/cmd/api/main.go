// Command api is the long-running entry point of the backend: it loads the
// configuration, opens the database connection, and serves the router that
// internal/app composes. With "migrate" as its argument it only applies the
// versioned migrations and exits.
package main

import (
	"log"
	"net/http"
	"os"

	"github.com/Unknowns24/ecommerce-mdw/config"
	"github.com/Unknowns24/ecommerce-mdw/internal/adapters/db/postgresql"
	database "github.com/Unknowns24/ecommerce-mdw/internal/adapters/db/postgresql/migrations"
	"github.com/Unknowns24/ecommerce-mdw/internal/app"
)

// version se completa en build con -ldflags "-X main.version=<git-sha-corto>".
var version = "dev"

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("configuración inválida: %v", err)
	}

	db, err := postgresql.Abrir(cfg.DatabaseDSN)
	if err != nil {
		log.Fatalf("no se pudo conectar a la base de datos: %v", err)
	}

	// Las migraciones no corren al arrancar el servidor: es un subcomando
	// explícito (./api migrate). Correrlas en el arranque haría que un
	// deploy con el server ya corriendo migre la base de la nube sin
	// querer. Ver ADR-002 y ADR-005.
	if len(os.Args) > 1 && os.Args[1] == "migrate" {
		migrador := database.NewMigrator(db, database.Todas())
		if err := migrador.Migrate(); err != nil {
			log.Fatalf("migración falló: %v", err)
		}
		log.Println("migraciones aplicadas correctamente")
		return
	}

	log.Printf("escuchando en :%s (version=%s)", cfg.AppPort, version)
	if err := http.ListenAndServe(":"+cfg.AppPort, app.Handler(cfg, db, version)); err != nil {
		log.Fatalf("el servidor se detuvo: %v", err)
	}
}
