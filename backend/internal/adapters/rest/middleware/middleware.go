// Package middleware protects HTTP endpoints with session and permission
// checks. This file is a stub: it only exists so the other three modules can
// compile and mount their routes from minute zero.
//
// TODO(Yasmín): reemplazar con el middleware real de sesión y permisos
// (JWT con HS256, verificación explícita del algoritmo, claims con roles y
// permisos). Esta firma está congelada: Agustín, Genaro y Nicolás ya
// programan contra ella. No la cambies sin avisar por el grupo.
package middleware

import (
	"context"
	"net/http"

	"github.com/google/uuid"

	apierr "github.com/Unknowns24/ecommerce-mdw/internal/adapters/rest/errors"
)

// Usuario is the identity carried in the request context once a session is
// verified. It never comes from the body, a custom header, or the query.
type Usuario struct {
	ID       uuid.UUID
	Correo   string
	Roles    []string
	Permisos []string
}

// Middleware holds whatever the session verification needs (today: the app
// secret used to verify JWTs).
type Middleware struct {
	secreto string
}

// Nuevo builds the middleware with the application secret used to verify
// sessions.
func Nuevo(secreto string) *Middleware {
	return &Middleware{secreto: secreto}
}

// RequiereSesion stub: siempre responde 401. TODO(Yasmín): verificar el JWT
// del header Authorization o la cookie "sesion" y meter el Usuario en el
// contexto.
func (m *Middleware) RequiereSesion(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		apierr.Responder(w, r.Method+" "+r.URL.Path, apierr.ErrNoAutenticado)
	})
}

// RequierePermiso stub: siempre responde 401 (en la versión real corre
// después de RequiereSesion y responde 403 si falta el permiso).
func (m *Middleware) RequierePermiso(permiso string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			apierr.Responder(w, r.Method+" "+r.URL.Path, apierr.ErrNoAutenticado)
		})
	}
}

// SesionOpcional stub: deja pasar sin usuario en el contexto. TODO(Yasmín):
// si hay un token válido, ponerlo en el contexto; si no, dejar pasar igual.
func (m *Middleware) SesionOpcional(next http.Handler) http.Handler {
	return next
}

type usuarioContextKey struct{}

// UsuarioDeContexto lee el Usuario que RequiereSesion/SesionOpcional dejaron
// en el contexto de la request.
func UsuarioDeContexto(ctx context.Context) (Usuario, bool) {
	u, ok := ctx.Value(usuarioContextKey{}).(Usuario)
	return u, ok
}
