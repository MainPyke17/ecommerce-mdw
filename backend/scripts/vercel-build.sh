#!/bin/sh
# Paso de build de Vercel. Aplica las migraciones versionadas contra la base
# de producción ANTES de que el deploy nuevo reciba tráfico: el mismo orden
# que el pipeline de la VPS (migrar y recién después levantar, ADR-002 y
# ADR-005). Si la migración falla, el build falla y sigue sirviendo el deploy
# anterior. La función (api/index.go) nunca migra.
set -eu

if [ "${VERCEL_ENV:-}" = "production" ]; then
	if ! command -v go >/dev/null 2>&1; then
		GO_VERSION=1.26.8
		case "$(uname -m)" in
		aarch64 | arm64) GO_ARCH=arm64 ;;
		*) GO_ARCH=amd64 ;;
		esac
		curl -fsSL "https://go.dev/dl/go${GO_VERSION}.linux-${GO_ARCH}.tar.gz" | tar -xz -C /tmp
		export PATH="/tmp/go/bin:$PATH"
	fi
	go version
	go run ./cmd/api migrate
	# Datos de demostración (catálogo, lotes y configuración de la tienda) para
	# que el flujo de compra se pueda recorrer en producción. Los tres seeds
	# son idempotentes: correrlos en cada deploy no duplica nada.
	go run ./cmd/seed-catalogo
	go run ./cmd/seed
	go run ./cmd/seed-pedidos
else
	echo "VERCEL_ENV=${VERCEL_ENV:-local}: sin migraciones (sólo producción tiene base configurada)"
fi

# Vercel exige un directorio de salida. Todo lo sirve la función de api/
# (vercel.json reescribe cada ruta hacia ella), así que queda vacío.
mkdir -p public
