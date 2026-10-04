package inventario

import (
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/Unknowns24/ecommerce-mdw/internal/core/domain"
)

type RepositorioReservas struct{ db *gorm.DB }

func NuevoRepositorioReservas(db *gorm.DB) *RepositorioReservas {
	return &RepositorioReservas{db: db}
}

func (r *RepositorioReservas) Crear(ctx context.Context, reserva *domain.ReservaStock) error {
	return r.db.WithContext(ctx).Create(reserva).Error
}

func (r *RepositorioReservas) ActivasPorPedido(ctx context.Context, pedidoID uuid.UUID) ([]domain.ReservaStock, error) {
	var reservas []domain.ReservaStock
	err := r.db.WithContext(ctx).
		Where("pedido_id = ? AND estado = ?", pedidoID, domain.ReservaActiva).
		Find(&reservas).Error
	return reservas, err
}

func (r *RepositorioReservas) ActivasPorLote(ctx context.Context, loteID uuid.UUID) ([]domain.ReservaStock, error) {
	var reservas []domain.ReservaStock
	err := r.db.WithContext(ctx).
		Where("lote_id = ? AND estado = ?", loteID, domain.ReservaActiva).
		Find(&reservas).Error
	return reservas, err
}

func (r *RepositorioReservas) CambiarEstado(ctx context.Context, id uuid.UUID, nuevo domain.EstadoReserva) error {
	return r.db.WithContext(ctx).
		Model(&domain.ReservaStock{}).
		Where("id = ?", id).
		Update("estado", nuevo).Error
}
