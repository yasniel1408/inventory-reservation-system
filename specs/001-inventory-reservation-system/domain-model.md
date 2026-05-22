# Modelo de Dominio: Sistema de Reservas de Inventario

## Entidades

### Item

- `id`: identificador único estable.
- `name`: nombre visible.
- `total_stock`: entero no negativo.
- `created_at`: timestamp de creación.
- `updated_at`: timestamp de actualización.

Reglas:

- `total_stock` nunca debe ser negativo.
- El stock disponible se deriva de reservas activas salvo que la implementación elija un contador denormalizado protegido por transacciones.

### Reservation

- `id`: identificador único estable.
- `item_id`: item asociado.
- `user_id` o `session_id`: owner de la reserva.
- `quantity`: entero positivo.
- `status`: `active`, `released`, `expired` o `confirmed`.
- `expires_at`: timestamp de creación más 60 segundos.
- `created_at`: timestamp de creación.
- `released_at`: timestamp nullable.
- `expired_at`: timestamp nullable.

Reglas:

- Solo reservas `active` cuentan contra stock disponible.
- `released`, `expired` y `confirmed` no cuentan como holds activos.
- Una reserva puede pasar de `active` a `released`.
- Una reserva puede pasar de `active` a `expired`.
- Una reserva no debe pasar de `released` a `active`.
- Una reserva no debe pasar de `expired` a `active`.
- Los efectos de devolución de stock deben ocurrir una sola vez por reserva.

### Idempotency Record

- `key`: idempotency key provista por el cliente.
- `scope`: endpoint más user/session identifier.
- `request_hash`: hash canónico de método, ruta y payload.
- `status`: `in_progress`, `completed` o `failed`.
- `reservation_id`: referencia nullable a reserva.
- `response_code`: response code almacenado para replay.
- `response_body`: response body almacenado para replay.
- `created_at`: timestamp de creación.
- `updated_at`: timestamp de actualización.

Reglas:

- `scope` más `key` debe ser único.
- Misma key más mismo request hash debe devolver el outcome almacenado.
- Misma key más distinto request hash debe devolver conflicto de idempotencia.
- Duplicados en progreso deben esperar, bloquear o devolver una respuesta retryable determinística. El enfoque elegido debe documentarse en `plan.md`.
- Los registros de idempotencia deben tener política de retención documentada en `plan.md`.

## Transiciones de Estado

```text
active -> released
active -> expired
active -> confirmed
released -> released (no-op idempotente)
expired -> expired (no-op idempotente)
confirmed -> confirmed (no-op idempotente para release salvo que una futura semántica de confirmación requiera otra cosa)
```

## Invariantes de Stock

- `total_stock >= 0`
- `active_reserved_stock >= 0`
- `available_stock = total_stock - active_reserved_stock`
- `available_stock >= 0`
- Crear reserva no debe hacer que `active_reserved_stock > total_stock`.
- Release y expiración no deben reducir `active_reserved_stock` por debajo de cero.
- Retry de creación de reserva no debe crear reservas activas duplicadas.
- Retry de release no debe devolver stock más de una vez.

## Límite de Concurrencia

La transacción de base de datos es el límite de concurrencia. La aplicación no debe depender de locks en memoria como mecanismo primario de correctness porque el servicio puede escalar horizontalmente.

Estrategias PostgreSQL candidatas a evaluar en `plan.md`:

- Update atómico condicional sobre fila de item, protegido por `WHERE available >= requested`.
- Row-level lock sobre item con `SELECT ... FOR UPDATE`.
- Isolation level serializable con retry explícito ante serialization failure.

La implementación elegida debe ser testeable con tests Go concurrentes y debe soportar registros de idempotencia dentro del mismo límite de correctness.

## Límite de Tiempo

El reloj de backend/base de datos es autoritativo para `created_at`, `expires_at`, `released_at` y `expired_at`. Los timers frontend son solo display y deben reconciliar contra estado backend.

## Decisión de Persistencia

Las reservas expiradas pueden quedar persistidas con status `expired` para auditabilidad y manejo idempotente de no-op. No deben aparecer en queries de reservas activas ni ser interactuables como holds activos.

