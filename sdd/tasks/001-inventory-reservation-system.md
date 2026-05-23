# Tasks: Sistema de Reservas de Inventario

## Plan Asociado

- `sdd/plans/001-inventory-reservation-system.md`

## Convenciones

- Cada tarea debe terminar con archivos modificados o una validación concreta.
- Mantener trazabilidad hacia `sdd/plans/001-inventory-reservation-system.md`, `sdd/user_histories/` y `sdd/specs/001-inventory-reservation-system/`.
- Mantener `sdd/TRACEABILITY.md` actualizado cuando se agregue, elimine o cambie cobertura de tareas.
- Marcar tareas completadas cambiando `[ ]` por `[x]`.

## Fase 0 - SDD

- [x] T-001 Validar coherencia inicial de SDD.
  - Owner: raíz del repo.
  - Referencias: `sdd/plans/001-inventory-reservation-system.md`, `sdd/user_histories/08_sdd_traceability.feature`.
  - Resultado: `sdd/plans.md`, `sdd/tasks.md`, `sdd/plans/`, `sdd/tasks/`, `sdd/user_histories/` y `sdd/specs/` quedan alineados antes de implementar.

- [x] T-002 Registrar decisión de simplificar arquitectura.
  - Owner: `sdd/plans/001-inventory-reservation-system.md`.
  - Referencias: `sdd/plans/001-inventory-reservation-system.md#decisiones-técnicas`.
  - Resultado: nota explícita de que no se usará Next.js, Tailwind, CQRS formal ni arquitectura hexagonal formal.

## Fase 1 - Estructura Base

- [x] T-003 Crear estructura de carpetas del proyecto.
  - Owner: raíz del repo.
  - Archivos/carpetas: `backend/`, `frontend/`, `db/migrations/`, `db/seeds/`, `openapi/`, `docs/`.
  - Referencias: `sdd/plans/001-inventory-reservation-system.md#entrega`.
  - Resultado: carpetas base versionadas con `.gitkeep`.

- [x] T-004 Crear `docker-compose.yml` para PostgreSQL local.
  - Owner: `docker-compose.yml`.
  - Referencias: `sdd/plans/001-inventory-reservation-system.md#entrega`, `skills/delivery-artifacts/SKILL.md`.
  - Validación: `docker compose config`.
  - Resultado: servicio PostgreSQL 16 local con volumen persistente, healthcheck y carpetas DB montadas; validado con `docker-compose config` porque el binario disponible localmente es `docker-compose`.

## Fase 2 - Base de Datos

- [x] T-005 Crear migración inicial PostgreSQL.
  - Owner: `db/migrations/`.
  - Tablas mínimas: `items`, `reservations`, `idempotency_keys`.
  - Referencias: `sdd/specs/001-inventory-reservation-system/domain-model.md`, `sdd/plans/001-inventory-reservation-system.md#estrategia-de-idempotencia`.
  - Criterios: constraints para cantidad positiva, stock no negativo, status conocido, unique index para `scope + key`.
  - Resultado: `db/migrations/001_initial_schema.sql` crea schema inicial con constraints de stock, reservas, idempotencia, FKs e índices base.

- [x] T-006 Crear seed data de revisión.
  - Owner: `db/seeds/`.
  - Referencias: `skills/delivery-artifacts/SKILL.md`, `sdd/specs/001-inventory-reservation-system/acceptance-criteria.md#inventario`.
  - Criterios: al menos un item con stock suficiente y un item low-stock.
  - Resultado: `db/seeds/001_seed_items.sql` crea seed determinístico con un item de stock holgado y uno low-stock.

## Fase 3 - Backend

- [x] T-007 Inicializar módulo Go backend.
  - Owner: `backend/`.
  - Stack: Go, Gin, GORM, PostgreSQL driver.
  - Referencias: `sdd/plans/001-inventory-reservation-system.md#backend`, `skills/backend-reservation-system/SKILL.md`.
  - Validación: `go test ./...` desde `backend/`.
  - Resultado: módulo Go inicializado con Gin, GORM, driver PostgreSQL, UUID y estructura backend acordada.

- [x] T-008 Implementar configuración y conexión DB.
  - Owner: `backend/internal/config/`, `backend/internal/store/`.
  - Referencias: `sdd/plans/001-inventory-reservation-system.md#backend`.
  - Criterios: conexión por env vars y helper transaccional reutilizable.
  - Resultado: config por env vars, conexión GORM y helper transaccional reutilizable.

- [x] T-009 Implementar lectura de inventario.
  - Owner: `backend/internal/items/`, `backend/internal/http/`.
  - Endpoint: `GET /items`.
  - Referencias: `sdd/user_histories/01_inventory_dashboard.feature`, `sdd/specs/001-inventory-reservation-system/api-spec.md#get-items`.
  - Criterios: devuelve nombre, stock total, reserved stock y available stock.
  - Resultado: handler y servicio de items devuelven stock total, reservado y disponible; ejecutan expiración lazy antes de leer.

- [x] T-010 Implementar creación de reserva atómica.
  - Owner: `backend/internal/reservations/`, `backend/internal/http/`.
  - Endpoint: `POST /reservations`.
  - Referencias: `sdd/user_histories/02_atomic_reservations.feature`, `sdd/plans/001-inventory-reservation-system.md#estrategia-de-concurrencia`.
  - Criterios: update condicional en transacción, sin oversell, validación de cantidad e item.
  - Resultado: creación usa transacción PostgreSQL y `UPDATE items ... WHERE total_stock - reserved_stock >= quantity RETURNING`.

- [x] T-011 Implementar idempotencia para `POST /reservations`.
  - Owner: `backend/internal/reservations/`, `backend/internal/store/`.
  - Referencias: `sdd/user_histories/05_idempotency.feature`, `sdd/plans/001-inventory-reservation-system.md#estrategia-de-idempotencia`.
  - Criterios: misma key/payload devuelve mismo outcome; misma key/payload distinto devuelve `409`.
  - Resultado: `Idempotency-Key` obligatorio, scope por session, request hash, `FOR UPDATE`, replay desde outcome almacenado y conflicto por payload distinto.

- [x] T-012 Implementar listado de reservas activas.
  - Owner: `backend/internal/reservations/`, `backend/internal/http/`.
  - Endpoint: `GET /reservations`.
  - Referencias: `sdd/user_histories/06_ui_feedback_and_state.feature`, `sdd/specs/001-inventory-reservation-system/api-spec.md#get-reservations`.
  - Criterios: no devuelve reservas `released` ni `expired`.
  - Resultado: lista reservas activas por `X-Session-ID`/session anónima después de expiración lazy.

- [x] T-013 Implementar release idempotente.
  - Owner: `backend/internal/reservations/`, `backend/internal/http/`.
  - Endpoint: `DELETE /reservations/{id}`.
  - Referencias: `sdd/user_histories/04_manual_release.feature`, `sdd/plans/001-inventory-reservation-system.md#estrategia-de-ttl-y-release`.
  - Criterios: release activo devuelve stock una vez; release repetido es no-op exitoso.
  - Resultado: release usa lock de fila, cambia `active -> released`, devuelve stock una vez y trata terminales como no-op.

- [x] T-014 Implementar expiración lazy de reservas.
  - Owner: `backend/internal/reservations/`.
  - Referencias: `sdd/user_histories/03_reservation_ttl.feature`, `sdd/plans/001-inventory-reservation-system.md#estrategia-de-ttl-y-release`.
  - Criterios: `active -> expired` devuelve stock una vez; se ejecuta antes de lecturas/mutaciones relevantes.
  - Resultado: expiración lazy usa CTE transaccional con `NOW()` antes de lecturas y mutaciones relevantes.

- [x] T-015 Implementar manejo consistente de errores API.
  - Owner: `backend/internal/http/`.
  - Referencias: `sdd/specs/001-inventory-reservation-system/api-spec.md#shape-de-error`.
  - Criterios: error shape estable y status codes documentados.
  - Resultado: handlers devuelven shape `{ "error": { "code", "message", "details" } }`, mapean errores de dominio y validan UUIDs antes de casts SQL.

## Fase 4 - Tests Backend

- [x] T-016 Test de concurrencia para última unidad.
  - Owner: `backend/internal/reservations/`.
  - Referencias: `sdd/specs/001-inventory-reservation-system/test-spec.md#concurrencia-última-unidad`.
  - Validación: exactamente 1 éxito, stock final 0 disponible.
  - Resultado: `TestPostgresCreateReservationLastUnitConcurrency` valida 50 goroutines, 1 éxito, `reserved_stock=1`, `available=0`.

- [x] T-017 Test de 100 requests para 10 unidades.
  - Owner: `backend/internal/reservations/`.
  - Referencias: `sdd/specs/001-inventory-reservation-system/test-spec.md#concurrencia-100-requests-por-10-unidades`.
  - Validación: 10 éxitos, 90 rechazos, sin stock negativo.
  - Resultado: `TestPostgresCreateReservationTenUnitsConcurrency` valida 100 goroutines, 10 éxitos, 90 rechazos, stock final correcto.

- [x] T-018 Test de idempotencia de reserva en paralelo.
  - Owner: `backend/internal/reservations/`.
  - Referencias: `sdd/specs/001-inventory-reservation-system/test-spec.md#idempotencia-de-reserva`.
  - Validación: una reserva y un decremento.
  - Resultado: `TestPostgresCreateReservationParallelIdempotency` exige mismo reservation ID para todos los outcomes paralelos y un solo decremento.

- [x] T-019 Test de release idempotente.
  - Owner: `backend/internal/reservations/`.
  - Referencias: `sdd/specs/001-inventory-reservation-system/test-spec.md#idempotencia-de-release`.
  - Validación: stock devuelto una sola vez.
  - Resultado: `TestPostgresReleaseReservationDoubleCallReturnsStockOnce` valida release paralelo y devolución única de stock.

- [x] T-020 Tests backend adicionales si el tiempo alcanza.
  - Owner: `backend/internal/reservations/`.
  - Referencias: `sdd/plans/001-inventory-reservation-system.md#tests`.
  - Casos: payload distinto con misma key, expiración doble, carrera release vs expiración.
  - Resultado: se cubren conflicto de idempotencia, expiración doble y carrera release vs expiración con estado final de DB.

## Fase 5 - Frontend

- [x] T-021 Inicializar frontend React + Vite + TypeScript.
  - Owner: `frontend/`.
  - Referencias: `sdd/plans/001-inventory-reservation-system.md#frontend`, `skills/frontend-reservation-app/SKILL.md`.
  - Validación: test/build inicial del frontend.
  - Resultado: frontend inicializado con Vite, React, TypeScript, ESLint y Vitest/jsdom; validado con `npm test`, `npm run build` y `npm run lint`.

- [x] T-022 Implementar cliente API tipado.
  - Owner: `frontend/src/api/`, `frontend/src/types/`.
  - Referencias: `sdd/specs/001-inventory-reservation-system/api-spec.md`.
  - Criterios: funciones para items, crear reserva, listar reservas y release.
  - Resultado: `createApiClient` tipado cubre `GET /items`, `POST /reservations`, `GET /reservations` y `DELETE /reservations/{id}` con `Idempotency-Key` y errores API.

- [x] T-023 Implementar dashboard de inventario.
  - Owner: `frontend/src/components/`.
  - Referencias: `sdd/user_histories/01_inventory_dashboard.feature`.
  - Criterios: muestra item, total, reservado/disponible y acción de reserva.
  - Resultado: `InventoryDashboard` muestra inventario, métricas de stock, cantidad y acción de reserva por item.

- [x] T-024 Implementar flujo de reserva con feedback.
  - Owner: `frontend/src/components/`, `frontend/src/hooks/`.
  - Referencias: `sdd/user_histories/06_ui_feedback_and_state.feature`.
  - Criterios: success, invalid quantity, insufficient stock, loading state.
  - Resultado: reserva genera key idempotente por intento, muestra success, cantidad inválida, stock insuficiente, error genérico y refetch posterior.

- [x] T-025 Implementar vista de reservas activas y release.
  - Owner: `frontend/src/components/`.
  - Referencias: `sdd/user_histories/04_manual_release.feature`, `sdd/user_histories/06_ui_feedback_and_state.feature`.
  - Criterios: lista reservas activas, botón release, refetch después de release.
  - Resultado: `ActiveReservations` lista reservas activas con botón de release, feedback y refetch posterior.

- [x] T-026 Implementar timer de expiración y reconciliación.
  - Owner: `frontend/src/hooks/`, `frontend/src/lib/`.
  - Referencias: `sdd/user_histories/03_reservation_ttl.feature`, `sdd/plans/001-inventory-reservation-system.md#estado-frontend`.
  - Criterios: timer llega a cero, no baja de cero, dispara refetch una vez.
  - Resultado: helper de timer calcula segundos restantes sin bajar de cero y dispara reconciliación una vez al expirar.

## Fase 6 - Tests Frontend

- [x] T-027 Unit test de timer.
  - Owner: `frontend/src/lib/` o `frontend/src/hooks/`.
  - Referencias: `sdd/specs/001-inventory-reservation-system/test-spec.md#unit-test-de-lógica-de-timer`.
  - Resultado: `frontend/src/lib/timer.test.ts` cubre countdown, límite cero y notificación única de expiración.

- [x] T-028 Component test de happy path de reserva.
  - Owner: `frontend/src/components/`.
  - Referencias: `sdd/specs/001-inventory-reservation-system/test-spec.md#component-test-de-happy-path-de-reserva`.
  - Resultado: `frontend/src/App.test.tsx` cubre reserva exitosa, feedback y refetch.

- [x] T-029 Component test de error por stock insuficiente.
  - Owner: `frontend/src/components/`.
  - Referencias: `sdd/specs/001-inventory-reservation-system/test-spec.md#component-test-de-estado-de-error`.
  - Resultado: `frontend/src/App.test.tsx` cubre error de stock insuficiente sin crear reserva falsa.

## Fase 7 - Contrato y Entrega

- [x] T-030 Crear `openapi/openapi.yaml`.
  - Owner: `openapi/openapi.yaml`.
  - Referencias: `sdd/user_histories/07_openapi_contract.feature`, `sdd/specs/001-inventory-reservation-system/api-spec.md`.
  - Criterios: documenta endpoints, schemas, headers, status codes y errores.
  - Resultado: `openapi/openapi.yaml` documenta endpoints implementados, `Idempotency-Key`, `X-Session-ID`, schemas, status codes y errores reales.

- [x] T-031 Crear `README.md`.
  - Owner: `README.md`.
  - Referencias: `skills/delivery-artifacts/SKILL.md`, `sdd/plans/001-inventory-reservation-system.md#entrega`.
  - Criterios: explica setup, tests, concurrencia, TTL, idempotencia y LLM usado.
  - Resultado: `README.md` explica stack, setup DB/backend/frontend, seeds, tests, concurrencia, TTL, idempotencia, OpenAPI y LLM usado.

- [x] T-032 Crear o documentar chat history.
  - Owner: `docs/chat-history.md` o README.
  - Referencias: `sdd/user_histories/08_sdd_traceability.feature`.
  - Criterios: el repo indica dónde está el historial completo de la conversación.
  - Resultado: `docs/chat-history.md` resume el flujo de conversación y `README.md` apunta al historial completo en Codex.

- [x] T-033 Documentar comandos y decisiones finales.
  - Owner: `README.md`.
  - Referencias: `sdd/plans/001-inventory-reservation-system.md`, `sdd/user_histories/08_sdd_traceability.feature`.
  - Criterios: comandos usados, supuestos relevantes y decisiones finales.
  - Resultado: `README.md` documenta comandos, decisiones de concurrencia, idempotencia, TTL, PostgreSQL init y validación final.

- [x] T-034 Ejecutar validación final.
  - Owner: repo completo.
  - Comandos esperados:
    - backend tests
    - frontend tests
    - frontend build
    - validación OpenAPI si hay herramienta disponible
  - Criterios: resultados documentados en `README.md`.
  - Resultado: validado con `ruby -e "require 'yaml'; YAML.load_file('openapi/openapi.yaml')"`, `go test ./...`, `npm test`, `npm run lint`, `npm run build` y `./scripts/validate-harness.sh`.

## Fase 8 - Smoke End-to-End y Hardening Final

- [x] T-035 Validar CORS para frontend Vite.
  - Owner: `backend/internal/http/`.
  - Referencias: `skills/frontend-reservation-app/SKILL.md`, `skills/backend-reservation-system/SKILL.md`, `README.md#frontend`.
  - Criterios: `OPTIONS /reservations` permite `Content-Type`, `Idempotency-Key` y `X-Session-ID`; responses API exponen `Access-Control-Allow-Origin`.
  - Resultado: se agrego test RED/GREEN `TestCorsPreflightAllowsFrontendReservationHeaders` y middleware CORS en router Gin.

- [x] T-036 Ejecutar smoke UI + API + PostgreSQL real.
  - Owner: repo completo.
  - Referencias: `README.md`, `openapi/openapi.yaml`, `sdd/specs/001-inventory-reservation-system/acceptance-criteria.md`.
  - Criterios: cargar inventario seed, crear reserva desde navegador, ver reserva activa, liberar y reconciliar stock sin errores de consola.
  - Resultado: smoke validado con Apple Container PostgreSQL en `192.168.64.2`, backend local `8080`, frontend Vite `5173` y Playwright.

- [x] T-037 Hardening visual responsive final.
  - Owner: `frontend/src/App.css`.
  - Referencias: `skills/frontend-reservation-app/SKILL.md`.
  - Criterios: no comprimir texto de reserva activa ni desbordar layout mobile.
  - Resultado: se agrego `box-sizing` local de app y layout propio para `reservation-row`.
