# Especificación: Sistema de Reservas de Inventario

## Estado

- Fase: Especificación
- Fuente: `Beeyond Media FS - Code Challenge (1).pdf`
- Historias de usuario: `user_histories/`
- Estado de `skills/` local: presente y obligatorio para ejecutar tareas del repo

## Objetivo

Construir un sistema de reservas de inventario para un contexto de flash sale con alto tráfico, donde usuarios puedan reservar stock temporalmente sin vender dos veces el mismo inventario bajo carga concurrente.

El sistema debe mostrar estado de stock, crear reservas de forma atómica, expirar reservas automáticamente después de 60 segundos, permitir release manual y exponer un contrato REST documentado.

## En Alcance

- Dashboard de inventario con nombre de item, stock total, stock reservado y stock disponible.
- Reserva de N unidades de un item por usuario/session.
- Control de concurrencia respaldado por PostgreSQL para impedir over-reservation.
- TTL de reserva de 60 segundos para reservas no confirmadas.
- Release manual de reservas activas, liberadas o ya expiradas con comportamiento idempotente.
- Creación idempotente de reservas usando header `Idempotency-Key`.
- Feedback de conflicto y validación en UI React.
- Vista de reservas activas con acción release y countdown.
- Estados loading, empty y error en frontend.
- Contrato OpenAPI para la REST API implementada.
- Seed data para reviewers.
- Tests backend y frontend requeridos por el challenge.
- Artefactos Architecture First: spec, plan, tasks, notes y trazabilidad.

## Fuera de Alcance

- Procesamiento de pagos.
- Checkout final o flujo de confirmación más allá del modelo necesario para reservar, expirar y liberar.
- Autenticación de usuarios más allá de un concepto mínimo user/session para asociar reservas activas con el usuario actual.
- WebSocket como único mecanismo de sincronización. Polling o refetch explícito es aceptable si mantiene UI sincronizada.
- Inventario multi-warehouse.
- Gestión admin de inventario más allá de seed/reviewer data.

## Actores

- Shopper: reserva y libera inventario desde el frontend.
- Frontend client: envía requests API, reintenta de forma segura y renderiza feedback.
- Reservation API: valida comandos, aplica reglas de concurrencia y devuelve estado canónico.
- Expiration worker: expira reservas después del TTL y devuelve stock una sola vez.
- Evaluador: valida artefactos, cobertura de tests y trazabilidad.

## Conceptos Core

- Item: registro de inventario vendible con identificador inmutable, nombre y stock total.
- Reservation: hold temporal por una cantidad de item owned por user/session.
- Stock disponible: `total_stock - active_reserved_stock`.
- Reserva activa: reserva no liberada, expirada ni confirmada.
- Reserva expirada: reserva cuyo TTL de 60 segundos pasó y cuyo stock fue devuelto exactamente una vez.
- Reserva liberada: reserva cancelada manualmente por el usuario y cuyo stock fue devuelto exactamente una vez.
- Registro de idempotencia: registro durable que vincula un `Idempotency-Key` al fingerprint original del request y su outcome.

## Requisitos Funcionales

### Dashboard de Inventario

- El sistema debe listar items de inventario.
- Cada item debe mostrar:
  - nombre del item
  - stock total
  - stock reservado por reservas activas
  - stock disponible calculado desde el estado backend canónico actual
- Si el inventario está vacío, la UI debe mostrar un empty state claro.
- Si falla la carga de inventario, la UI debe mostrar un error recuperable.
- Las acciones de reserva deben estar deshabilitadas o protegidas mientras los datos de inventario no estén listos.

### Reservas Atómicas

- Un usuario puede solicitar reserva por una cantidad entera positiva de un item.
- Backend debe rechazar cantidades cero, negativas, faltantes, no enteras o inválidas.
- Backend debe rechazar requests para items inexistentes.
- Backend debe rechazar requests cuando el stock disponible es menor a la cantidad solicitada.
- La creación de reservas debe ser atómica bajo requests concurrentes.
- Bajo ninguna condición el stock disponible puede quedar negativo.
- La actualización canónica de stock debe ocurrir en PostgreSQL dentro de una transacción.
- API debe devolver la reserva creada con identificador, item, cantidad, status, timestamp de creación y timestamp de expiración.

### TTL de Reserva

- Toda reserva exitosa expira 60 segundos después de su creación salvo que sea liberada antes o confirmada por un flujo futuro.
- Cuando el TTL pasa, la reserva ya no debe ser interactuable como hold activo.
- La expiración debe devolver la cantidad reservada al pool disponible exactamente una vez.
- La expiración debe ser durable y segura ante retries o múltiples ejecuciones de worker.
- La expiración debe usar tiempo de backend/base de datos como fuente de verdad, no el reloj frontend.
- La expiración puede implementarse mediante worker, scheduled job, lazy cleanup en lecturas/escrituras o combinación. El mecanismo elegido debe documentarse en `plan.md`.
- El timer frontend es solo orientativo; el estado backend es autoritativo.
- La UI debe dejar de mostrar reservas expiradas como activas después de sincronizar.

### Release Manual

- Un usuario puede liberar manualmente su reserva en cualquier momento.
- Liberar una reserva activa debe devolver stock exactamente una vez.
- Liberar una reserva ya liberada debe ser éxito o devolver un response no-op bien definido.
- Liberar una reserva ya expirada debe ser éxito o devolver un response no-op bien definido.
- Release manual no debe recrear, reactivar ni devolver stock dos veces.
- La UI debe mostrar feedback visible después del release.

### Idempotencia

- `POST /reservations` debe requerir header `Idempotency-Key`.
- Las idempotency keys se scopean por endpoint y user/session salvo que `plan.md` elija un scope global más estricto.
- La persistencia de idempotencia debe guardar un request hash canónico y el outcome final replayable.
- Dos requests de reserva con la misma key y mismo payload deben devolver el mismo outcome de reserva.
- Misma key y mismo payload no deben decrementar stock más de una vez.
- Dos requests con misma key y payload distinto deben rechazarse con error claro de conflicto de idempotencia.
- Requests paralelos con misma key y mismo payload deben converger a un outcome almacenado.
- Si llega un request duplicado mientras el primero sigue en progreso, API debe esperar, bloquear o devolver respuesta retryable determinística. La opción elegida debe documentarse en `plan.md`.
- Responses de validación y stock-conflict pueden almacenarse como outcomes idempotentes si son determinísticos. Errores `5xx` transitorios no deben cachearse como outcomes finales exitosos.
- `DELETE /reservations/{id}` debe ser seguro de llamar repetidamente.
- Retries de release no deben devolver stock más de una vez.

### Manejo de Conflictos

- Backend debe devolver error claro cuando el stock fue consumido por otro request antes de que el request actual commitee.
- Frontend debe mostrar errores de conflicto/stock insuficiente sin romperse.
- Después de un conflicto, frontend debe refrescar o reconciliar estado de inventario para no mostrar stock stale.

### Estado Frontend

- Frontend debe usar React, Vite y TypeScript.
- La primera pantalla útil debe ser la experiencia de inventario/reserva, no una landing page.
- La UI debe proveer:
  - lista de inventario
  - input/acción de cantidad a reservar por item
  - lista de reservas activas
  - botón release por reserva activa
  - timer por reserva activa
  - loading states
  - empty states
  - API error states
  - success feedback para reserve y release
- Frontend debe mantenerse sincronizado con backend usando polling, refetch después de mutations u otro mecanismo explícito de sincronización.
- UI debe tratar responses backend como canónicas. Optimistic UI es permitida solo si se reconcilia después del response del server.
- Double-clicks y requests lentos no deben disparar transiciones visuales duplicadas.

### OpenAPI

- La REST API implementada debe tener contrato OpenAPI.
- El contrato debe documentar request bodies, response bodies, headers, status codes y error shapes.
- El contrato debe incluir lista de inventario, creación de reserva, lista de reservas activas y release de reserva.

## Requisitos No Funcionales

- Lenguaje backend: Go.
- Dependencias backend: Go, Gin y GORM.
- Base de datos: PostgreSQL.
- Backend: Go con Gin y GORM, manteniendo una estructura simple y testeable.
- Frontend: React + Vite + TypeScript.
- El proyecto debe ajustarse a un target de implementación de 8-10 horas enfocadas y parar a las 12 horas.
- El código debe priorizar correctness bajo escrituras concurrentes sobre optimización prematura.
- API debe ser determinística ante retries.
- Test coverage debe probar concurrencia e idempotencia.
- README debe explicar estrategia de concurrencia, ejecución de tests y uso de LLM.
- `spec-kit-notes.md` debe documentar comandos, supuestos, refinamientos y pivots.

## Criterios de Aceptación

- Dado 50 o más requests concurrentes por la última unidad disponible, exactamente una reserva es exitosa y el stock nunca queda negativo.
- Dado 100 requests concurrentes por 10 unidades disponibles, exactamente 10 reservas son exitosas y exactamente 90 son rechazadas.
- Dado dos requests paralelos con mismo `Idempotency-Key` y payload, ambos devuelven el mismo outcome y stock decrementa una vez.
- Dado dos requests de release para la misma reserva, stock se devuelve una vez.
- Dado una reserva con más de 60 segundos, se expira y ya no aparece como activa.
- Dado una reserva expirada, release manual es no-op seguro o éxito bien definido.
- Dado UI stale donde el stock acaba de ser tomado, el intento de reserva muestra conflicto/stock insuficiente claro y refresca estado.
- Dado input de cantidad inválida, UI bloquea o API rechaza con error de validación claro.
- Dada la implementación API, su contrato OpenAPI describe todos los endpoints públicos.

## Casos Límite

- Misma idempotency key con mismo payload llega antes de que el primer request termine.
- Misma idempotency key con payload distinto llega después de un request exitoso.
- Misma idempotency key con payload distinto llega mientras el request original está en progreso.
- Reserva expira mientras el usuario hace click en release.
- Release se llama dos veces en paralelo.
- Expiration worker procesa la misma reserva más de una vez.
- Lecturas de inventario ocurren mientras se crean o expiran reservas.
- Cantidades concurrentes mixtas, por ejemplo stock 10 con requests simultáneos por 7 y 5 unidades.
- Timer del usuario llega a cero antes de que backend procese expiración.
- Backend expira la reserva antes de que el timer frontend llegue a cero.
- Transacción DB falla después de registrar idempotency key pero antes de tener outcome de reserva.
- Cliente reintenta después de timeout de red, pero el primer request sí commiteó.

## Supuestos

- Un identificador mínimo anónimo o user/session sembrado alcanza para el challenge salvo que se agregue autenticación explícitamente.
- No se requiere endpoint de confirmación para el challenge, solo comportamiento de reservas no confirmadas.
- Reservas expiradas pueden permanecer en base para audit/debugging, pero deben eliminarse permanentemente de vistas activas/interactuables.
- PostgreSQL es la fuente de verdad para stock y estado de reservas.
- Timestamps en API responses son UTC ISO 8601.
- La referencia visual del PDF guía la dirección visual, pero la implementación debe ser pragmática y accesible.

## Trazabilidad

- Historias de usuario:
  - `user_histories/01_inventory_dashboard.feature`
  - `user_histories/02_atomic_reservations.feature`
  - `user_histories/03_reservation_ttl.feature`
  - `user_histories/04_manual_release.feature`
  - `user_histories/05_idempotency.feature`
  - `user_histories/06_ui_feedback_and_state.feature`
  - `user_histories/07_openapi_contract.feature`
  - `user_histories/08_spec_kit_traceability.feature`
- Próximos artefactos:
  - `plan.md`
  - `tasks.md`
  - `spec-kit-notes.md`
  - contrato OpenAPI
