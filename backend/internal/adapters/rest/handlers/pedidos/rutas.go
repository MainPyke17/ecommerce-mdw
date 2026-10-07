// Package pedidos mounts the cart, checkout and order endpoints.
package pedidos

import (
	"context"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/Unknowns24/ecommerce-mdw/config"
	repo "github.com/Unknowns24/ecommerce-mdw/internal/adapters/db/postgresql/repositories/pedidos"
	"github.com/Unknowns24/ecommerce-mdw/internal/adapters/rest/middleware"
	"github.com/Unknowns24/ecommerce-mdw/internal/core/ports/out"
	catalogo "github.com/Unknowns24/ecommerce-mdw/internal/usecases/catalogo"
	inventario "github.com/Unknowns24/ecommerce-mdw/internal/usecases/inventario"
	uc "github.com/Unknowns24/ecommerce-mdw/internal/usecases/pedidos"
)

// Montar attaches this module's routes to the shared router. main.go calls
// it exactly once.
//
// Acá se conectan los casos de uso con las implementaciones de los otros
// módulos: el catálogo de Genaro (precio vigente de cada variante) y el
// servicio de stock de Agustín (disponibilidad, reserva, confirmación y
// liberación). Los dos adaptadores de abajo sólo traducen tipos.
func Montar(r chi.Router, db *gorm.DB, cfg config.Config, mw *middleware.Middleware) {
	pedidos := repo.NuevoRepositorioPedidos(db)
	tx := repo.NuevasTransacciones(db)
	stock := inventario.NuevoServicioStock(db)
	lector := lectorCatalogo{catalogo.NuevoLector(db)}

	h := NuevoHandler(
		uc.NuevoCheckout(lector, reservadorStock{stock}, pedidos, repo.NuevoRepositorioConfiguracion(db), tx, cfg.PaymentReservationTTL, time.Now),
		uc.NuevaConsulta(pedidos),
		uc.NuevaGestion(pedidos, stock, tx, time.Now),
		uc.NuevaGestionCarrito(repo.NuevoRepositorioCarritos(db), lector, stock),
	)

	// Checkout: la sesión es opcional, el invitado compra igual.
	r.With(mw.SesionOpcional).Post("/api/pedidos", h.CrearPedido)
	// Enlace privado del invitado: público a propósito, el secreto es el token.
	r.Get("/api/pedidos/publico/{token}", h.ObtenerPedidoPublico)

	cuenta := r.With(mw.RequiereSesion)
	cuenta.Get("/api/pedidos", h.ListarMisPedidos)
	cuenta.Get("/api/pedidos/{id}", h.ObtenerPedido)
	cuenta.Get("/api/carrito", h.VerCarrito)
	cuenta.Delete("/api/carrito", h.VaciarCarrito)
	cuenta.Post("/api/carrito/items", h.AgregarAlCarrito)
	cuenta.Patch("/api/carrito/items/{id}", h.CambiarUnidadesDelCarrito)
	cuenta.Delete("/api/carrito/items/{id}", h.QuitarDelCarrito)

	admin := r.With(mw.RequiereSesion, mw.RequierePermiso("pedidos.gestionar"))
	admin.Get("/api/admin/pedidos", h.ListarPedidosAdmin)
	admin.Get("/api/admin/pedidos/{id}", h.ObtenerPedidoAdmin)
	admin.Post("/api/admin/pedidos/{id}/despacho", h.DespacharPedido)
	admin.Post("/api/admin/pedidos/{id}/entrega", h.EntregarPedido)
	admin.Post("/api/admin/pedidos/{id}/cobro-efectivo", h.CobrarPedidoEnEfectivo)
	admin.Post("/api/admin/pedidos/{id}/cancelacion", h.CancelarPedido)
}

// reservadorStock adapta ServicioStock.Reservar, que recibe
// []inventario.ItemReserva, al puerto out.ReservadorStock del checkout.
type reservadorStock struct{ stock *inventario.ServicioStock }

func (a reservadorStock) Reservar(ctx context.Context, tx *gorm.DB, pedidoID uuid.UUID, items []out.ItemReservaStock) error {
	convertidos := make([]inventario.ItemReserva, len(items))
	for i, it := range items {
		convertidos[i] = inventario.ItemReserva{VarianteID: it.VarianteID, Unidades: it.Unidades}
	}
	return a.stock.Reservar(ctx, tx, pedidoID, convertidos)
}

// lectorCatalogo adapta el Lector del catálogo, que devuelve
// catalogo.VarianteVendible, al puerto out.LectorCatalogo.
type lectorCatalogo struct{ lector *catalogo.Lector }

func (a lectorCatalogo) VariantePorID(ctx context.Context, id uuid.UUID) (out.VarianteVendible, error) {
	v, err := a.lector.VariantePorID(ctx, id)
	if err != nil {
		return out.VarianteVendible{}, err
	}
	return out.VarianteVendible{
		ID:                      v.ID,
		ProductoID:              v.ProductoID,
		Codigo:                  v.Codigo,
		Nombre:                  v.Nombre,
		PrecioMinoristaCentavos: v.PrecioMinoristaCentavos,
		Activa:                  v.Activa,
	}, nil
}

// Comprobaciones en compilación: si una firma deja de coincidir con su
// puerto, el proyecto no compila.
var (
	_ out.ReservadorStock = reservadorStock{}
	_ out.LectorCatalogo  = lectorCatalogo{}
	_ out.GestorStock     = (*inventario.ServicioStock)(nil)
	_ out.ConsultaStock   = (*inventario.ServicioStock)(nil)
)
