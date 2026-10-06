package catalogo

import (
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"

	repo "github.com/Unknowns24/ecommerce-mdw/internal/adapters/db/postgresql/repositories/catalogo"
)

type Publico struct{ repo *repo.RepositorioPublico }

func NuevoPublico(db *gorm.DB) *Publico { return &Publico{repo: repo.NuevoRepositorioPublico(db)} }

type PaginaVariantes struct {
	Items     []repo.VariantePublica `json:"items"`
	Pagina    int                    `json:"pagina"`
	PorPagina int                    `json:"porPagina"`
	Total     int64                  `json:"total"`
}

func (s *Publico) Listar(ctx context.Context, f repo.Filtros, pagina, porPagina int) (PaginaVariantes, error) {
	items, total, err := s.repo.Listar(ctx, f, porPagina, (pagina-1)*porPagina)
	if err != nil {
		return PaginaVariantes{}, err
	}
	// TODO(Genaro): conectar disponibilidad en lote cuando inventario esté en develop.
	for i := range items {
		items[i].Disponible = true
	}
	return PaginaVariantes{Items: items, Pagina: pagina, PorPagina: porPagina, Total: total}, nil
}

type DetalleVariante struct {
	repo.VariantePublica
	Imagenes       []string               `json:"imagenes"`
	OtrasVariantes []repo.VariantePublica `json:"otrasVariantes"`
}

func (s *Publico) Detalle(ctx context.Context, id uuid.UUID) (DetalleVariante, error) {
	item, err := s.repo.PorID(ctx, id)
	if err != nil {
		return DetalleVariante{}, err
	}
	imagenes, err := s.repo.Imagenes(ctx, id)
	if err != nil {
		return DetalleVariante{}, err
	}
	otras, err := s.repo.Otras(ctx, item.ProductoID, id)
	if err != nil {
		return DetalleVariante{}, err
	}
	item.Disponible = true
	for i := range otras {
		otras[i].Disponible = true
	}
	return DetalleVariante{VariantePublica: item, Imagenes: imagenes, OtrasVariantes: otras}, nil
}

func (s *Publico) Marcas(ctx context.Context) ([]map[string]any, error) { return s.repo.Marcas(ctx) }
func (s *Publico) Categorias(ctx context.Context) ([]map[string]any, error) {
	return s.repo.Categorias(ctx)
}
