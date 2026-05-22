# Plan: Sistema de Reservas de Inventario

## Objetivo

Implementar el challenge con una solución simple y verificable: backend Go con PostgreSQL para evitar oversell, frontend React + Vite + TypeScript para operar reservas, contrato OpenAPI y tests que prueben concurrencia e idempotencia.

## Artefactos Base

- Historias: `user_histories/`
- Specs: `specs/001-inventory-reservation-system/`
- Skills activos: `skills/`
- Próximos artefactos: implementación, OpenAPI, README y chat history

## Decisiones Técnicas

### Backend

- Lenguaje: Go.
- Router HTTP: Gin.
- ORM/acceso DB: GORM.
- Base de datos: PostgreSQL.
- Migraciones: SQL versionadas en `db/migrations/`.
- Seed data: SQL en `db/seeds/`.
- Estructura simple:

```text
backend/
  cmd/api/
  internal/config/
  internal/http/
  internal/store/
  internal/items/
  internal/reservations/
```

No se usará una arquitectura hexagonal/CQRS formal. Para este challenge es suficiente separar HTTP, store y lógica de reservas, manteniendo el código testeable.

### Frontend

- Framework: React.
- Build tool: Vite.
- Lenguaje: TypeScript.
- Estilos: CSS simple con componentes propios.
- Estructura simple:

```text
frontend/
  src/
    api/
    components/
    hooks/
    lib/
    types/
```

No se usará Next.js porque el PDF pide React + Vite + TypeScript.

### API

Endpoints previstos:

- `GET /items`
- `POST /reservations`
- `GET /reservations`
- `DELETE /reservations/{id}`

El contrato final vivirá en `openapi/openapi.yaml` y debe coincidir con la implementación.

## Estrategia de Concurrencia

PostgreSQL será la fuente de verdad. La creación de reservas correrá dentro de una transacción.

Estrategia elegida:

1. Mantener en `items` un contador `reserved_stock`.
2. Crear reserva mediante update condicional:

```sql
UPDATE items
SET reserved_stock = reserved_stock + ?
WHERE id = ?
  AND total_stock - reserved_stock >= ?
RETURNING id, total_stock, reserved_stock;
```

3. Si el update no devuelve fila, responder `409 insufficient_stock`.
4. Insertar la reserva activa en la misma transacción.

Motivo: es simple, atómico, testeable y evita oversell sin locks en memoria.

## Estrategia de Idempotencia

`POST /reservations` requiere `Idempotency-Key`.

Tabla prevista: `idempotency_keys`.

Campos mínimos:

- `scope`
- `key`
- `request_hash`
- `status`
- `response_code`
- `response_body`
- `reservation_id`
- timestamps

Reglas:

- `scope + key` es único.
- Misma key y mismo payload devuelve el mismo outcome.
- Misma key y payload distinto devuelve `409 idempotency_key_conflict`.
- Requests paralelos con la misma key se serializan usando transacción y row lock sobre el registro de idempotencia.
- Errores transitorios `5xx` no se guardan como outcome exitoso.

## Estrategia de TTL y Release

- Cada reserva activa tiene `expires_at = created_at + 60 segundos`.
- El backend/base de datos es la fuente de tiempo.
- Las reservas expiradas no aparecen como activas.
- Expiración se manejará con cleanup lazy:
  - antes de listar inventario/reservas
  - antes de mutaciones de reserva
- La expiración usa update transaccional para cambiar `active -> expired` y devolver stock una sola vez.
- `DELETE /reservations/{id}` cambia `active -> released` y devuelve stock una sola vez.
- Si la reserva ya está `released` o `expired`, `DELETE` responde no-op exitoso.

Motivo: evita sumar workers o cron para un challenge corto, y sigue cumpliendo TTL mientras el sistema recibe tráfico.

## Estado Frontend

- El backend es canónico.
- El frontend hace refetch después de:
  - reserve
  - release
  - timer en cero
- El timer es visual; no decide estado final.
- Cada intento de reserva genera una `Idempotency-Key`.
- El mismo intento reutiliza la key durante retries.
- Errores de stock insuficiente, cantidad inválida y conflictos se muestran en pantalla.

## Tests

### Backend

Tests obligatorios:

- 50+ requests concurrentes por última unidad: 1 éxito.
- 100 requests concurrentes por 10 unidades: 10 éxitos, 90 rechazos.
- Misma `Idempotency-Key` en paralelo: una reserva, un decremento.
- Release doble: stock devuelto una vez.

Tests adicionales si el tiempo alcanza:

- Payload distinto con misma idempotency key.
- Expiración dos veces no devuelve stock dos veces.
- Carrera release vs expiración.

### Frontend

Tests obligatorios:

- Unit test de lógica de timer.
- Component test del happy path de reserva.
- Component test de error por stock insuficiente.

## Entrega

Archivos finales esperados:

- `backend/`
- `frontend/`
- `db/migrations/`
- `db/seeds/`
- `openapi/openapi.yaml`
- `README.md`
- `docs/chat-history.md` o ruta documentada del historial

## Riesgos y Mitigaciones

- Riesgo: TTL lazy cleanup no expira reservas si no hay tráfico.
  - Mitigación: documentarlo como decisión de challenge; opcionalmente agregar endpoint/función invocable en tests.
- Riesgo: GORM puede ocultar detalles de locking.
  - Mitigación: usar SQL crudo para updates condicionales críticos dentro de transacciones GORM.
- Riesgo: frontend queda stale.
  - Mitigación: refetch después de mutations y cuando timer llega a cero.
- Riesgo: demasiada arquitectura para 8-10 horas.
  - Mitigación: mantener estructura simple y enfocada en tests requeridos.

## Trazabilidad

- Dashboard: `user_histories/01_inventory_dashboard.feature`
- Reservas atómicas: `user_histories/02_atomic_reservations.feature`
- TTL: `user_histories/03_reservation_ttl.feature`
- Release manual: `user_histories/04_manual_release.feature`
- Idempotencia: `user_histories/05_idempotency.feature`
- UI/estado: `user_histories/06_ui_feedback_and_state.feature`
- OpenAPI: `user_histories/07_openapi_contract.feature`
- Entrega SDD: `user_histories/08_sdd_traceability.feature`
