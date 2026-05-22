# Spec Kit Notes: Sistema de Reservas de Inventario

## Estado

- Fecha: 2026-05-22.
- Fase actual: listo para iniciar implementacion.
- Fuente primaria: `Beeyond Media FS - Code Challenge (1).pdf`.
- Artefactos base:
  - `user_histories/`
  - `specs/001-inventory-reservation-system/`
  - `plan.md`
  - `tasks.md`
  - `skills/`

## Comandos y Revisiones Ejecutadas

- `find skills -type f -maxdepth 3 -print | sort`
- `sed -n '1,220p' skills/SELECTING_SKILLS.md`
- `sed -n '1,220p' skills/development-flow/SKILL.md`
- `sed -n '1,220p' skills/spec-kit-architecture/SKILL.md`
- `sed -n '1,240p' skills/backend-reservation-system/SKILL.md`
- `sed -n '1,220p' skills/frontend-reservation-app/SKILL.md`
- `sed -n '1,220p' skills/delivery-artifacts/SKILL.md`
- `sed -n '1,240p' plan.md`
- `sed -n '1,240p' tasks.md`
- `sed -n '1,260p' specs/001-inventory-reservation-system/spec.md`
- `sed -n '1,220p' specs/001-inventory-reservation-system/test-spec.md`
- `rg -n "Next\\.js|Next|Tailwind|CQRS|hexagonal|agent-|Team Leader|Analista|Reviewer|route groups|App Router|sub-agente|subagente" .`
- `python3 /Users/yasnielfajardo/.codex/skills/.system/skill-creator/scripts/quick_validate.py "$d"` para cada skill local con `SKILL.md`.

## Supuestos

- No hay autenticacion real requerida para el challenge; se usara un identificador minimo de usuario o sesion para scopear reservas e idempotencia.
- PostgreSQL es la fuente de verdad para stock, reservas, expiracion e idempotencia.
- Los timestamps expuestos por API se manejaran como UTC ISO 8601.
- `Idempotency-Key` es obligatorio para `POST /reservations`.
- `DELETE /reservations/{id}` no necesita `Idempotency-Key` porque su comportamiento debe ser naturalmente idempotente.
- La expiracion de 60 segundos debe cumplirse con tiempo de backend/base de datos, no con reloj del frontend.
- Las reservas expiradas o liberadas pueden quedar persistidas para auditoria, pero no deben aparecer como activas.

## Decisiones

### Stack

- Backend: Go + Gin + GORM + PostgreSQL.
- Frontend: React + Vite + TypeScript.
- Migraciones: SQL versionadas en `db/migrations/`.
- Seeds: SQL en `db/seeds/`.
- OpenAPI final: `openapi/openapi.yaml`.

### Simplificacion de Arquitectura

Se descarta usar Next.js, Tailwind, CQRS formal y arquitectura hexagonal formal para la implementacion de este challenge.

Motivo:

- El PDF pide React + Vite + TypeScript para frontend.
- El objetivo principal evaluable es correctness bajo concurrencia, TTL, release, idempotencia y trazabilidad.
- Una arquitectura grande agregaria costo sin mejorar la validacion central del challenge.

La estructura backend queda simple y testeable:

```text
backend/
  cmd/api/
  internal/config/
  internal/http/
  internal/store/
  internal/items/
  internal/reservations/
```

La estructura frontend queda simple:

```text
frontend/
  src/
    api/
    components/
    hooks/
    lib/
    types/
```

## Estrategia de Concurrencia

La reserva se implementara con una transaccion PostgreSQL y un update condicional sobre `items.reserved_stock`.

Forma esperada:

```sql
UPDATE items
SET reserved_stock = reserved_stock + ?
WHERE id = ?
  AND total_stock - reserved_stock >= ?
RETURNING id, total_stock, reserved_stock;
```

Si no hay fila retornada, la API debe responder `409 insufficient_stock`.

No se usaran locks en memoria como mecanismo de correctness.

## Estrategia de Idempotencia

Se usara tabla `idempotency_keys` con unicidad por `scope + key`.

Reglas:

- Misma key y mismo payload devuelve el mismo outcome.
- Misma key y payload distinto devuelve `409 idempotency_key_conflict`.
- Requests paralelos con la misma key se serializan dentro de transaccion.
- El outcome replayable debe guardar status code y response body.
- Errores transitorios `5xx` no deben guardarse como outcome exitoso final.

## Estrategia de TTL y Release

- Cada reserva activa tendra `expires_at = created_at + 60 segundos`.
- La expiracion se implementara como cleanup lazy antes de lecturas y mutaciones relevantes.
- `active -> expired` debe devolver stock una sola vez.
- `active -> released` debe devolver stock una sola vez.
- Reintentos de release o expiracion deben ser no-op seguros.

Riesgo aceptado:

- Si no hay trafico, una reserva puede no expirar fisicamente hasta la siguiente lectura o mutacion.

Mitigacion:

- Documentarlo en README.
- Exponer la logica como funcion testeable para probar expiracion sin worker permanente.

## Refinamientos y Pivots

- Inicialmente se considero una arquitectura mas amplia con skills muy atomicos y agentes de proceso.
- Se redujo a pocos skills operativos para evitar ceremonia y mantener foco en el challenge.
- Se regreso de Next.js/Tailwind a React + Vite + TypeScript por alineacion directa con el PDF.
- Se descarto CQRS/hexagonal formal y se preservo una separacion practica por HTTP, store, items y reservations.
- Se mantiene Gin/GORM porque fue una decision tecnica aceptada por el usuario y no contradice el PDF, usando SQL crudo donde GORM no sea suficiente para concurrencia.

## Trazabilidad

- Dashboard de inventario: `user_histories/01_inventory_dashboard.feature`, `specs/001-inventory-reservation-system/spec.md`.
- Reservas atomicas: `user_histories/02_atomic_reservations.feature`, `plan.md#estrategia-de-concurrencia`.
- TTL: `user_histories/03_reservation_ttl.feature`, `plan.md#estrategia-de-ttl-y-release`.
- Release manual: `user_histories/04_manual_release.feature`.
- Idempotencia: `user_histories/05_idempotency.feature`, `plan.md#estrategia-de-idempotencia`.
- UI y estado: `user_histories/06_ui_feedback_and_state.feature`.
- OpenAPI: `user_histories/07_openapi_contract.feature`, `specs/001-inventory-reservation-system/api-spec.md`.
- Entrega Spec Kit: `user_histories/08_spec_kit_traceability.feature`.

## Validaciones Pendientes

Estas validaciones deben completarse durante o al final de la implementacion:

- `docker compose config`
- tests backend desde `backend/`
- tests frontend desde `frontend/`
- build frontend desde `frontend/`
- validacion del contrato OpenAPI si hay herramienta disponible

## Checklist de Implementacion

- Crear estructura base del repo.
- Crear `docker-compose.yml` para PostgreSQL.
- Crear migracion inicial.
- Crear seeds.
- Implementar backend.
- Implementar tests backend requeridos.
- Implementar frontend.
- Implementar tests frontend requeridos.
- Crear `openapi/openapi.yaml`.
- Crear `README.md`.
- Documentar chat history o ruta del historial.
- Actualizar este archivo con comandos reales finales y resultados.
