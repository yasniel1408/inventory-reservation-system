---
name: backend-reservation-system
description: Usar para implementar, planificar o revisar el backend Go del sistema de reservas con Gin, GORM y PostgreSQL, incluyendo migraciones versionadas, schema, handlers REST, transacciones, prevención de oversell, TTL de 60 segundos, release idempotente, Idempotency-Key, seed compatibility y tests de concurrencia.
---

# Backend Reservation System

Usar este skill para todo el backend. Mantenerlo simple: el challenge evalúa concurrencia, idempotencia y trazabilidad, no una arquitectura grande.

## Stack

- Go.
- Gin para HTTP.
- GORM para acceso a PostgreSQL.
- Migraciones SQL versionadas; no depender de AutoMigrate como contrato final.

## Estructura Recomendada

```text
backend/
  cmd/api/
  internal/config/
  internal/http/
  internal/store/
  internal/reservations/
  internal/items/
```

## Reglas Críticas

- PostgreSQL es la fuente de verdad.
- Si el frontend corre en Vite u otro origen distinto al backend, el router debe exponer CORS y responder `OPTIONS` para `Content-Type`, `Idempotency-Key` y `X-Session-ID`.
- Toda mutación de reserva corre en transacción.
- Usar SQL crudo dentro de GORM cuando haga falta locking, conditional update o `RETURNING`.
- No usar locks en memoria como mecanismo de correctness.
- `POST /reservations` requiere `Idempotency-Key`.
- `DELETE /reservations/{id}` es idempotente.
- Expiración y release no pueden devolver stock más de una vez.
- El replay idempotente debe devolver el outcome almacenado, no reconstruirlo desde estado vivo que pudo cambiar por release o expiración.
- Si se persiste un fallo determinístico de idempotencia, la transacción debe commitear ese outcome y luego devolver el error al handler; no hacer rollback del registro idempotente.
- Los timestamps canónicos de reservas, expiración y release deben venir de PostgreSQL cuando afecten TTL o auditoría.
- Validar UUIDs de entrada antes de casts SQL para devolver errores `4xx` estables en lugar de `500`.

## Tests Backend Obligatorios

- 50+ requests concurrentes por la última unidad: exactamente un éxito.
- 100 requests concurrentes por 10 unidades: 10 éxitos y 90 rechazos.
- Misma `Idempotency-Key` en paralelo: una reserva y un decremento.
- Release doble: stock devuelto una sola vez.
- Los tests de integración PostgreSQL deben verificar estado final de DB, no solo conteo de errores/responses.
- Usar `TEST_DATABASE_URL` explícito para tests DB; no intentar levantar contenedores desde los tests.
- Si se usa Apple Container y `127.0.0.1:5432` está ocupado por otro servicio, usar la IP del contenedor o un puerto publicado no conflictivo.
