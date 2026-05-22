# Tasks: Sistema de Reservas de Inventario

## Plan Asociado

- `plans/001-inventory-reservation-system.md`

## Convenciones

- Cada tarea debe terminar con archivos modificados o una validación concreta.
- Mantener trazabilidad hacia `plans/001-inventory-reservation-system.md`, `user_histories/` y `specs/001-inventory-reservation-system/`.
- Marcar tareas completadas cambiando `[ ]` por `[x]`.

## Fase 0 - SDD

- [x] T-001 Validar coherencia inicial de SDD.
  - Owner: raíz del repo.
  - Referencias: `plans/001-inventory-reservation-system.md`, `user_histories/08_sdd_traceability.feature`.
  - Resultado: `plans.md`, `tasks.md`, `plans/`, `tasks/`, `user_histories/` y `specs/` quedan alineados antes de implementar.

- [x] T-002 Registrar decisión de simplificar arquitectura.
  - Owner: `plans/001-inventory-reservation-system.md`.
  - Referencias: `plans/001-inventory-reservation-system.md#decisiones-técnicas`.
  - Resultado: nota explícita de que no se usará Next.js, Tailwind, CQRS formal ni arquitectura hexagonal formal.

## Fase 1 - Estructura Base

- [ ] T-003 Crear estructura de carpetas del proyecto.
  - Owner: raíz del repo.
  - Archivos/carpetas: `backend/`, `frontend/`, `db/migrations/`, `db/seeds/`, `openapi/`, `docs/`.
  - Referencias: `plans/001-inventory-reservation-system.md#entrega`.

- [ ] T-004 Crear `docker-compose.yml` para PostgreSQL local.
  - Owner: `docker-compose.yml`.
  - Referencias: `plans/001-inventory-reservation-system.md#entrega`, `skills/delivery-artifacts/SKILL.md`.
  - Validación: `docker compose config`.

## Fase 2 - Base de Datos

- [ ] T-005 Crear migración inicial PostgreSQL.
  - Owner: `db/migrations/`.
  - Tablas mínimas: `items`, `reservations`, `idempotency_keys`.
  - Referencias: `specs/001-inventory-reservation-system/domain-model.md`, `plans/001-inventory-reservation-system.md#estrategia-de-idempotencia`.
  - Criterios: constraints para cantidad positiva, stock no negativo, status conocido, unique index para `scope + key`.

- [ ] T-006 Crear seed data de revisión.
  - Owner: `db/seeds/`.
  - Referencias: `skills/delivery-artifacts/SKILL.md`, `specs/001-inventory-reservation-system/acceptance-criteria.md#inventario`.
  - Criterios: al menos un item con stock suficiente y un item low-stock.

## Fase 3 - Backend

- [ ] T-007 Inicializar módulo Go backend.
  - Owner: `backend/`.
  - Stack: Go, Gin, GORM, PostgreSQL driver.
  - Referencias: `plans/001-inventory-reservation-system.md#backend`, `skills/backend-reservation-system/SKILL.md`.
  - Validación: `go test ./...` desde `backend/`.

- [ ] T-008 Implementar configuración y conexión DB.
  - Owner: `backend/internal/config/`, `backend/internal/store/`.
  - Referencias: `plans/001-inventory-reservation-system.md#backend`.
  - Criterios: conexión por env vars y helper transaccional reutilizable.

- [ ] T-009 Implementar lectura de inventario.
  - Owner: `backend/internal/items/`, `backend/internal/http/`.
  - Endpoint: `GET /items`.
  - Referencias: `user_histories/01_inventory_dashboard.feature`, `specs/001-inventory-reservation-system/api-spec.md#get-items`.
  - Criterios: devuelve nombre, stock total, reserved stock y available stock.

- [ ] T-010 Implementar creación de reserva atómica.
  - Owner: `backend/internal/reservations/`, `backend/internal/http/`.
  - Endpoint: `POST /reservations`.
  - Referencias: `user_histories/02_atomic_reservations.feature`, `plans/001-inventory-reservation-system.md#estrategia-de-concurrencia`.
  - Criterios: update condicional en transacción, sin oversell, validación de cantidad e item.

- [ ] T-011 Implementar idempotencia para `POST /reservations`.
  - Owner: `backend/internal/reservations/`, `backend/internal/store/`.
  - Referencias: `user_histories/05_idempotency.feature`, `plans/001-inventory-reservation-system.md#estrategia-de-idempotencia`.
  - Criterios: misma key/payload devuelve mismo outcome; misma key/payload distinto devuelve `409`.

- [ ] T-012 Implementar listado de reservas activas.
  - Owner: `backend/internal/reservations/`, `backend/internal/http/`.
  - Endpoint: `GET /reservations`.
  - Referencias: `user_histories/06_ui_feedback_and_state.feature`, `specs/001-inventory-reservation-system/api-spec.md#get-reservations`.
  - Criterios: no devuelve reservas `released` ni `expired`.

- [ ] T-013 Implementar release idempotente.
  - Owner: `backend/internal/reservations/`, `backend/internal/http/`.
  - Endpoint: `DELETE /reservations/{id}`.
  - Referencias: `user_histories/04_manual_release.feature`, `plans/001-inventory-reservation-system.md#estrategia-de-ttl-y-release`.
  - Criterios: release activo devuelve stock una vez; release repetido es no-op exitoso.

- [ ] T-014 Implementar expiración lazy de reservas.
  - Owner: `backend/internal/reservations/`.
  - Referencias: `user_histories/03_reservation_ttl.feature`, `plans/001-inventory-reservation-system.md#estrategia-de-ttl-y-release`.
  - Criterios: `active -> expired` devuelve stock una vez; se ejecuta antes de lecturas/mutaciones relevantes.

- [ ] T-015 Implementar manejo consistente de errores API.
  - Owner: `backend/internal/http/`.
  - Referencias: `specs/001-inventory-reservation-system/api-spec.md#shape-de-error`.
  - Criterios: error shape estable y status codes documentados.

## Fase 4 - Tests Backend

- [ ] T-016 Test de concurrencia para última unidad.
  - Owner: `backend/internal/reservations/`.
  - Referencias: `specs/001-inventory-reservation-system/test-spec.md#concurrencia-última-unidad`.
  - Validación: exactamente 1 éxito, stock final 0 disponible.

- [ ] T-017 Test de 100 requests para 10 unidades.
  - Owner: `backend/internal/reservations/`.
  - Referencias: `specs/001-inventory-reservation-system/test-spec.md#concurrencia-100-requests-por-10-unidades`.
  - Validación: 10 éxitos, 90 rechazos, sin stock negativo.

- [ ] T-018 Test de idempotencia de reserva en paralelo.
  - Owner: `backend/internal/reservations/`.
  - Referencias: `specs/001-inventory-reservation-system/test-spec.md#idempotencia-de-reserva`.
  - Validación: una reserva y un decremento.

- [ ] T-019 Test de release idempotente.
  - Owner: `backend/internal/reservations/`.
  - Referencias: `specs/001-inventory-reservation-system/test-spec.md#idempotencia-de-release`.
  - Validación: stock devuelto una sola vez.

- [ ] T-020 Tests backend adicionales si el tiempo alcanza.
  - Owner: `backend/internal/reservations/`.
  - Referencias: `plans/001-inventory-reservation-system.md#tests`.
  - Casos: payload distinto con misma key, expiración doble, carrera release vs expiración.

## Fase 5 - Frontend

- [ ] T-021 Inicializar frontend React + Vite + TypeScript.
  - Owner: `frontend/`.
  - Referencias: `plans/001-inventory-reservation-system.md#frontend`, `skills/frontend-reservation-app/SKILL.md`.
  - Validación: test/build inicial del frontend.

- [ ] T-022 Implementar cliente API tipado.
  - Owner: `frontend/src/api/`, `frontend/src/types/`.
  - Referencias: `specs/001-inventory-reservation-system/api-spec.md`.
  - Criterios: funciones para items, crear reserva, listar reservas y release.

- [ ] T-023 Implementar dashboard de inventario.
  - Owner: `frontend/src/components/`.
  - Referencias: `user_histories/01_inventory_dashboard.feature`.
  - Criterios: muestra item, total, reservado/disponible y acción de reserva.

- [ ] T-024 Implementar flujo de reserva con feedback.
  - Owner: `frontend/src/components/`, `frontend/src/hooks/`.
  - Referencias: `user_histories/06_ui_feedback_and_state.feature`.
  - Criterios: success, invalid quantity, insufficient stock, loading state.

- [ ] T-025 Implementar vista de reservas activas y release.
  - Owner: `frontend/src/components/`.
  - Referencias: `user_histories/04_manual_release.feature`, `user_histories/06_ui_feedback_and_state.feature`.
  - Criterios: lista reservas activas, botón release, refetch después de release.

- [ ] T-026 Implementar timer de expiración y reconciliación.
  - Owner: `frontend/src/hooks/`, `frontend/src/lib/`.
  - Referencias: `user_histories/03_reservation_ttl.feature`, `plans/001-inventory-reservation-system.md#estado-frontend`.
  - Criterios: timer llega a cero, no baja de cero, dispara refetch una vez.

## Fase 6 - Tests Frontend

- [ ] T-027 Unit test de timer.
  - Owner: `frontend/src/lib/` o `frontend/src/hooks/`.
  - Referencias: `specs/001-inventory-reservation-system/test-spec.md#unit-test-de-lógica-de-timer`.

- [ ] T-028 Component test de happy path de reserva.
  - Owner: `frontend/src/components/`.
  - Referencias: `specs/001-inventory-reservation-system/test-spec.md#component-test-de-happy-path-de-reserva`.

- [ ] T-029 Component test de error por stock insuficiente.
  - Owner: `frontend/src/components/`.
  - Referencias: `specs/001-inventory-reservation-system/test-spec.md#component-test-de-estado-de-error`.

## Fase 7 - Contrato y Entrega

- [ ] T-030 Crear `openapi/openapi.yaml`.
  - Owner: `openapi/openapi.yaml`.
  - Referencias: `user_histories/07_openapi_contract.feature`, `specs/001-inventory-reservation-system/api-spec.md`.
  - Criterios: documenta endpoints, schemas, headers, status codes y errores.

- [ ] T-031 Crear `README.md`.
  - Owner: `README.md`.
  - Referencias: `skills/delivery-artifacts/SKILL.md`, `plans/001-inventory-reservation-system.md#entrega`.
  - Criterios: explica setup, tests, concurrencia, TTL, idempotencia y LLM usado.

- [ ] T-032 Crear o documentar chat history.
  - Owner: `docs/chat-history.md` o README.
  - Referencias: `user_histories/08_sdd_traceability.feature`.
  - Criterios: el repo indica dónde está el historial completo de la conversación.

- [ ] T-033 Documentar comandos y decisiones finales.
  - Owner: `README.md`.
  - Referencias: `plans/001-inventory-reservation-system.md`, `user_histories/08_sdd_traceability.feature`.
  - Criterios: comandos usados, supuestos relevantes y decisiones finales.

- [ ] T-034 Ejecutar validación final.
  - Owner: repo completo.
  - Comandos esperados:
    - backend tests
    - frontend tests
    - frontend build
    - validación OpenAPI si hay herramienta disponible
  - Criterios: resultados documentados en `README.md`.
