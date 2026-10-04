package inventario

import (
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/Unknowns24/ecommerce-mdw/internal/core/domain"
)

type RepositorioMovimientos struct{ db *gorm.DB }

func NuevoRepositorioMovimientos(db *gorm.DB) *RepositorioMovimientos {
	return &RepositorioMovimientos{db: db}
}

func (r *RepositorioMovimientos) Crear(ctx context.Context, movimiento *domain.MovimientoStock) error {
	return r.db.WithContext(ctx).Create(movimiento).Error
}

// PorVariante pagina el historial de movimientos de una variante, más
// recientes primero.
func (r *RepositorioMovimientos) PorVariante(ctx context.Context, varianteID uuid.UUID, limit, offset int) ([]domain.MovimientoStock, error) {
	var movimientos []domain.MovimientoStock
	err := r.db.WithContext(ctx).
		Joins("JOIN lote ON lote.id = movimiento_stock.lote_id").
		Where("lote.variante_id = ?", varianteID).
		Order("movimiento_stock.creado_en DESC").
		Limit(limit).
		Offset(offset).
		Find(&movimientos).Error
	return movimientos, err
}
