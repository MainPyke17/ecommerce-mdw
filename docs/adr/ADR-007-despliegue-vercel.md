# ADR-007 — Despliegue del Parcial I en Vercel con Neon

- Estado: aceptado
- Fecha: 2026-10-06
- Reemplaza parcialmente a: [ADR-005](ADR-005-despliegue.md) (la VPS sigue siendo la opción para el Parcial II)

## Contexto

La condición de admisión del Parcial I es una URL pública que responda y un login que funcione. Al 06/10 la VPS de ADR-005 no existía: el pipeline `desplegar` nunca corrió y `main` no tenía el código de la aplicación. Hacía falta una URL estable el mismo día, sin máquina propia que administrar.

## Decisión

Se despliega el mismo backend en Vercel, con la base PostgreSQL en Neon instalada desde el Marketplace de Vercel:

- **Una sola función** (`backend/api/index.go`) sirve el router completo; `vercel.json` reescribe todas las rutas hacia ella. El router se arma en `internal/app`, el mismo paquete que usa `cmd/api`, así que el binario de la VPS y la función de Vercel no pueden exponer rutas distintas.
- **Las migraciones corren en el build** (`backend/scripts/vercel-build.sh`), antes de que el deploy nuevo reciba tráfico, y sólo en producción. Si fallan, el build falla y sigue sirviendo el deploy anterior. La función nunca migra, igual que el servidor de la VPS (ADR-002).
- **La conexión** se lee de `DATABASE_DSN` o, si no está, de `DATABASE_URL`, que es la variable que publica la integración de Neon. Así el secreto no se copia a mano a otra variable.
- **Los secretos** (`APP_SECRET_KEY`, `ORDER_ACCESS_TOKEN_SECRET`, `BOOTSTRAP_MASTER_KEY`) viven en las variables de entorno de producción del proyecto de Vercel, nunca en el repositorio.

## Consecuencias

**Ganamos:** URL pública estable con HTTPS sin administrar una máquina, y *preview deployments* disponibles si se conecta el repositorio.

**Resignamos:** el arranque en frío de la función agrega latencia a la primera request de cada instancia, y cada instancia abre su propio pool de conexiones (se usa la conexión con *pooler* de Neon). El job `desplegar` del CI sigue apuntando a la VPS y no se usa para este despliegue: se despliega con `vercel deploy --prod` desde `backend/`.
