# Especificación API: Sistema de Reservas de Inventario

Este documento define la superficie REST prevista. El archivo OpenAPI final debe coincidir con la API implementada.

## Reglas Generales

- Todos los timestamps son strings UTC ISO 8601.
- Los requests se asocian al usuario/session actual. Para el challenge es aceptable una session anónima sembrada o header simple si queda documentado en `plan.md`.
- `Idempotency-Key` está scopeado por endpoint y user/session.
- Los error codes deben ser estables y documentados en OpenAPI.

## Shape de Error

Todos los errores API deben usar un shape consistente:

```json
{
  "error": {
    "code": "insufficient_stock",
    "message": "Not enough stock available.",
    "details": {}
  }
}
```

## `GET /items`

Lista items de inventario con estado de stock calculado.

Response `200`:

```json
{
  "items": [
    {
      "id": "item_123",
      "name": "Campaign Slot A",
      "totalStock": 10,
      "reservedStock": 4,
      "availableStock": 6
    }
  ]
}
```

Errores:

- `500 internal_error`

## `POST /reservations`

Crea una reserva idempotente.

Headers requeridos:

- `Idempotency-Key`: key única generada por el cliente para este intento de reserva.

Request:

```json
{
  "itemId": "item_123",
  "quantity": 2
}
```

Response `201` para reserva recién creada:

```json
{
  "reservation": {
    "id": "res_123",
    "itemId": "item_123",
    "quantity": 2,
    "status": "active",
    "expiresAt": "2026-05-22T12:01:00Z",
    "createdAt": "2026-05-22T12:00:00Z"
  }
}
```

Response `200` para replay idempotente de un outcome exitoso existente:

```json
{
  "reservation": {
    "id": "res_123",
    "itemId": "item_123",
    "quantity": 2,
    "status": "active",
    "expiresAt": "2026-05-22T12:01:00Z",
    "createdAt": "2026-05-22T12:00:00Z"
  },
  "idempotentReplay": true
}
```

Errores:

- `400 invalid_quantity`
- `400 missing_idempotency_key`
- `404 item_not_found`
- `409 insufficient_stock`
- `409 idempotency_key_conflict`
- `425 idempotency_in_progress` si la implementación elegida devuelve respuesta retryable para duplicado en progreso en vez de esperar
- `500 internal_error`

## `GET /reservations`

Lista reservas activas del usuario/session actual.

Response `200`:

```json
{
  "reservations": [
    {
      "id": "res_123",
      "itemId": "item_123",
      "itemName": "Campaign Slot A",
      "quantity": 2,
      "status": "active",
      "expiresAt": "2026-05-22T12:01:00Z",
      "createdAt": "2026-05-22T12:00:00Z"
    }
  ]
}
```

Errores:

- `500 internal_error`

## `DELETE /reservations/{id}`

Libera una reserva de forma idempotente.

Response `200` para release efectivo:

```json
{
  "reservation": {
    "id": "res_123",
    "status": "released",
    "releasedAt": "2026-05-22T12:00:30Z"
  },
  "stockReturned": true
}
```

Response `200` para release no-op:

```json
{
  "reservation": {
    "id": "res_123",
    "status": "released"
  },
  "stockReturned": false,
  "noop": true
}
```

Errores:

- `404 reservation_not_found`
- `500 internal_error`

## Expectativas de Sincronización

- Después de `POST /reservations`, frontend debe refrescar items y reservas activas o aplicar la response canónica de forma segura.
- Después de `DELETE /reservations/{id}`, frontend debe refrescar items y reservas activas o aplicar la response canónica de forma segura.
- La expiración del timer en frontend debe disparar refresh/reconciliación con estado backend.

## Schemas OpenAPI a Proveer

- `InventoryItem`
- `Reservation`
- `CreateReservationRequest`
- `CreateReservationResponse`
- `ReleaseReservationResponse`
- `ErrorResponse`
- `IdempotencyConflictError`

