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

// LectorPedidos son las consultas de lectura que el comprador puede hacer. Las
// de datos propios llevan el id del comprador: no existe una forma de pedir "el
// pedido X" sin decir de quién es.
type LectorPedidos interface {
	PorIDYUsuario(ctx context.Context, id, usuarioID uuid.UUID) (domain.Pedido, error)
	PorToken(ctx context.Context, token string) (domain.PedidoPublico, error)
	ListarDeUsuario(ctx context.Context, usuarioID uuid.UUID, limit, offset int) ([]domain.Pedido, int64, error)
}
