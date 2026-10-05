package out

import (
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/Unknowns24/ecommerce-mdw/internal/core/domain"
)

// ItemReservaStock es una línea de reserva: cuántas unidades de una variante
// se comprometen para un pedido. Replica el contrato congelado de Agustín
// (usecases/inventario.ItemReserva).
type ItemReservaStock struct {
	VarianteID uuid.UUID
	Unidades   int
}

// ReservadorStock reserva unidades dentro de la transacción del checkout. Si
// no alcanzan, devuelve apierr.ErrRegla con Datos = mapa variante → unidades
// disponibles. Lo implementa ServicioStock (Agustín) a través de un adaptador.
type ReservadorStock interface {
	Reservar(ctx context.Context, tx *gorm.DB, pedidoID uuid.UUID, items []ItemReservaStock) error
}

// CreadorDePedidos guarda un pedido con su detalle dentro de la transacción tx.
type CreadorDePedidos interface {
	Crear(ctx context.Context, tx *gorm.DB, p *domain.Pedido, detalles []domain.DetallePedido) error
}

// LectorConfiguracion entrega la configuración vigente de la tienda.
type LectorConfiguracion interface {
	Actual(ctx context.Context) (domain.ConfiguracionTienda, error)
	TarifasOrdenadas(ctx context.Context) ([]domain.TarifaDistancia, error)
}

// Transaccionador ejecuta fn dentro de una transacción: si fn devuelve error,
// se revierte todo lo que hizo; si no, se confirma.
type Transaccionador interface {
	EnTransaccion(ctx context.Context, fn func(tx *gorm.DB) error) error
}
