# Inventory Reservation System

Sistema de reservas de inventario para un escenario de alta concurrencia. La solucion prioriza correctness de stock, idempotencia, TTL de reservas y una entrega simple de revisar.

## Stack

- Backend: Go, Gin, GORM y PostgreSQL.
- Frontend: React, Vite y TypeScript.
- DB: migraciones SQL versionadas y seed SQL.
- Contrato: OpenAPI en `openapi/openapi.yaml`.

## Estructura

```text
backend/          API Go
frontend/         UI React + Vite
db/migrations/    schema PostgreSQL
db/seeds/         datos iniciales
db/init/          bootstrap para container PostgreSQL
openapi/          contrato OpenAPI
sdd/              specs, plan, tasks y trazabilidad
docs/             notas de entrega
skills/           reglas locales del harness
```

## PostgreSQL local

El repo incluye un `docker-compose.yml` con PostgreSQL 16. En esta maquina se valido con `docker-compose`, pero `docker compose` tambien deberia funcionar si tu runtime lo soporta.

```bash
docker-compose up -d postgres
```

Para validar la configuracion:

```bash
docker-compose config
```

La base por defecto es:

```text
postgres://inventory_user:inventory_password@localhost:5432/inventory_reservation?sslmode=disable
```

Seed inicial:

- `Standard Widget`, stock `100`.
- `Limited Edition Widget`, stock `2`.
- El seed es idempotente con `ON CONFLICT (id) DO UPDATE`.

La imagen oficial de PostgreSQL ejecuta scripts de `/docker-entrypoint-initdb.d` solo cuando el volumen se inicializa por primera vez. Este repo monta `db/init/001_run_db_scripts.sh`, que corre `db/migrations/` y `db/seeds/` en orden. Si ya existe el volumen y queres re-bootstrap desde cero:

```bash
docker-compose down -v
docker-compose up -d postgres
```

Si usas Apple Container y `localhost:5432` esta ocupado, publica otro puerto o usa la IP del container en `DATABASE_URL`/`TEST_DATABASE_URL`.

## Backend

Instalar dependencias y correr tests portables:

```bash
cd backend
go test ./...
```

Los tests de integracion PostgreSQL se activan solo con `TEST_DATABASE_URL`; si la variable no existe, se saltean.

```bash
cd backend
TEST_DATABASE_URL='postgres://inventory_user:inventory_password@localhost:5432/inventory_reservation?sslmode=disable' go test ./...
```

Correr la API:

```bash
cd backend
DATABASE_URL='postgres://inventory_user:inventory_password@localhost:5432/inventory_reservation?sslmode=disable' go run ./cmd/api
```

Variables:

- `PORT`: puerto HTTP, default `8080`.
- `DATABASE_URL`: DSN PostgreSQL, default apunta al compose local.

## Frontend

Instalar dependencias:

```bash
cd frontend
npm install
```

Correr tests, lint y build:

```bash
cd frontend
npm test
npm run lint
npm run build
```

Correr la UI:

```bash
cd frontend
VITE_API_BASE_URL='http://localhost:8080' npm run dev
```

## API

Contrato completo:

- `openapi/openapi.yaml`

Endpoints implementados:

- `GET /items`
- `POST /reservations`
- `GET /reservations`
- `DELETE /reservations/{id}`

Headers:

- `Idempotency-Key`: requerido en `POST /reservations`.
- `X-Session-ID`: opcional; si no se envia, se usa `anonymous-session`.

La API habilita CORS para el frontend local de Vite. El preflight de `POST /reservations` permite `Content-Type`, `Idempotency-Key` y `X-Session-ID`.

Los errores usan el shape estable:

```json
{
  "error": {
    "code": "insufficient_stock",
    "message": "Not enough stock available.",
    "details": {}
  }
}
```

## Estrategia de concurrencia

PostgreSQL es la fuente de verdad. La creacion de reservas corre en transaccion y usa un update condicional sobre `items.reserved_stock`:

```sql
UPDATE items
SET reserved_stock = reserved_stock + ?
WHERE id = ?
  AND total_stock - reserved_stock >= ?
RETURNING id;
```

Si el update no devuelve fila, no se crea reserva y se responde `409 insufficient_stock`. Esto evita oversell sin locks en memoria.

## Idempotencia

`POST /reservations` requiere `Idempotency-Key`.

La tabla `idempotency_keys` guarda `scope`, `key`, hash del payload, status y response body. Reglas:

- Misma key + mismo payload devuelve el mismo outcome.
- Misma key + payload distinto devuelve `409 idempotency_key_conflict`.
- Requests paralelos con la misma key se serializan con row lock.
- Los outcomes deterministas como stock insuficiente se persisten para replay estable.

## TTL y release

- Cada reserva activa expira a los 60 segundos.
- El tiempo canonico viene de PostgreSQL.
- La expiracion es lazy: se ejecuta antes de listar inventario/reservas y antes de mutaciones relevantes.
- La expiracion cambia `active -> expired` y devuelve stock una sola vez.
- `DELETE /reservations/{id}` cambia `active -> released` y devuelve stock una sola vez.
- Repetir `DELETE` sobre una reserva terminal responde `200` con `noop: true`.

## Validacion final

Comandos usados para validar esta entrega:

```bash
cd backend && go test ./...
cd frontend && npm test
cd frontend && npm run lint
cd frontend && npm run build
ruby -e "require 'yaml'; YAML.load_file('openapi/openapi.yaml'); puts 'OpenAPI YAML OK'"
./scripts/validate-harness.sh
```

Validacion CORS usada para el smoke navegador:

```bash
curl -i -H 'Origin: http://localhost:5173' http://localhost:8080/items
curl -i -X OPTIONS http://localhost:8080/reservations \
  -H 'Origin: http://localhost:5173' \
  -H 'Access-Control-Request-Method: POST' \
  -H 'Access-Control-Request-Headers: content-type,idempotency-key'
```

Smoke end-to-end ejecutado:

- PostgreSQL en Apple Container: `192.168.64.2:5432`.
- Backend local: `http://localhost:8080`.
- Frontend Vite: `http://127.0.0.1:5173`.
- Flujo validado en navegador: cargar seed, reservar `Standard Widget`, ver reserva activa, liberar y reconciliar stock sin errores de consola.

Para validacion de concurrencia con PostgreSQL real:

```bash
cd backend
TEST_DATABASE_URL='postgres://inventory_user:inventory_password@localhost:5432/inventory_reservation?sslmode=disable' go test ./...
```

## LLM y proceso

Se uso OpenAI Codex basado en GPT-5 como agente principal de planificacion, implementacion asistida, revision y generacion de artefactos de entrega. El trabajo siguio el harness local del repo:

- `skills/development-flow/SKILL.md`
- `skills/sdd-architecture/SKILL.md`
- `skills/backend-reservation-system/SKILL.md`
- `skills/frontend-reservation-app/SKILL.md`
- `skills/delivery-artifacts/SKILL.md`

El flujo aplicado fue: analyst, team-leader, developer, reviewer, tester, delivery-manager y skills-expert. Los sub-agentes se usaron para revisiones acotadas cuando el trabajo podia dividirse sin pisar archivos.

## Chat history

El historial completo vive en la conversacion de Codex usada para construir este repo. Para la entrega del repositorio, `docs/chat-history.md` resume el camino seguido y apunta a los artefactos verificables.
