# Especificación de Tests

## Tests Backend

### Concurrencia: Última Unidad

Preparación:

- Sembrar un item con stock total `1`.
- Iniciar 50 o más goroutines.
- Cada goroutine intenta reservar cantidad `1`.

Esperado:

- Exactamente un request es exitoso.
- Todos los demás fallan con stock insuficiente o conflicto.
- Stock reservado activo final es `1`.
- Stock disponible final es `0`.
- No existe stock negativo.

### Concurrencia: 100 Requests por 10 Unidades

Preparación:

- Sembrar un item con stock total `10`.
- Iniciar 100 goroutines.
- Cada goroutine intenta reservar cantidad `1`.

Esperado:

- Exactamente 10 requests son exitosos.
- Exactamente 90 requests fallan con stock insuficiente o conflicto.
- Stock reservado activo final es `10`.
- Stock disponible final es `0`.
- No existe stock negativo.

### Concurrencia: Cantidades Mixtas

Preparación:

- Sembrar un item con stock total `10`.
- Iniciar al menos dos requests concurrentes donde uno reserva cantidad `7` y otro cantidad `5`.

Esperado:

- Como máximo un request es exitoso.
- Stock reservado activo final es `7` o `5`.
- Stock disponible final es `3` o `5`.
- Stock reservado activo final nunca supera `10`.
- No existe stock negativo.

### Idempotencia de Reserva

Preparación:

- Sembrar un item con stock suficiente.
- Enviar dos requests paralelos `POST /reservations` con el mismo `Idempotency-Key` y el mismo payload.

Esperado:

- Ambos requests devuelven el mismo reservation ID y el mismo outcome.
- Solo se crea una reserva.
- Stock se decrementa una sola vez.

### Conflicto de Idempotencia

Preparación:

- Enviar un request `POST /reservations` con un `Idempotency-Key`.
- Enviar otro request `POST /reservations` con la misma key pero payload distinto.

Esperado:

- El segundo request falla con conflicto de idempotencia.
- No se crea reserva adicional.
- El stock no cambia por el request conflictivo.

### Idempotencia de Release

Preparación:

- Crear una reserva activa.
- Llamar `DELETE /reservations/{id}` dos veces, opcionalmente en paralelo.

Esperado:

- Ambas llamadas son exitosas o no-ops bien definidos.
- Stock se devuelve exactamente una vez.
- La reserva no está activa después del release.

### Idempotencia de Expiración

Preparación:

- Crear una reserva con `expires_at` en el pasado.
- Ejecutar lógica de expiración dos veces.

Esperado:

- La reserva queda expirada.
- Stock se devuelve una vez.
- El segundo pase de expiración es no-op.

### Carrera Expiración y Release

Preparación:

- Crear una reserva que debe expirar.
- Ejecutar lógica de expiración y release manual en paralelo.

Esperado:

- El status final de la reserva es `expired` o `released`, según la primera transición commiteada.
- Stock se devuelve exactamente una vez.
- La reserva no queda activa.

## Tests Frontend

### Unit Test de Lógica de Timer

Esperado:

- Calcula segundos restantes desde `expiresAt`.
- Llega a cero sin ir a negativo.
- Dispara callback de expiración/reconciliación una sola vez cuando el tiempo llega a cero.

### Component Test de Happy Path de Reserva

Esperado:

- Renderiza item de inventario.
- Usuario ingresa cantidad válida y envía.
- Muestra loading mientras el request está pending.
- Muestra success después de reservar.
- Refresca o actualiza reservas activas y stock disponible.

### Component Test de Estado de Error

Esperado:

- Usuario intenta reservar stock no disponible.
- API devuelve stock insuficiente.
- UI muestra mensaje de error.
- UI no agrega una reserva activa falsa.
- UI reconcilia estado de inventario.
