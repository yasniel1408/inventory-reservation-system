# Criterios de Aceptación

## Inventario

- El dashboard muestra todos los items sembrados.
- Cada item muestra nombre, stock total, stock reservado y stock disponible.
- Inventario vacío muestra un estado empty.
- Error al cargar inventario muestra error y opción de retry.

## Creación de Reservas

- Una cantidad positiva contra stock disponible crea una reserva activa.
- Cantidad `0` es rechazada.
- Cantidad negativa es rechazada.
- Cantidad faltante es rechazada.
- Item inexistente es rechazado.
- Stock insuficiente es rechazado.
- Requests concurrentes no pueden generar oversell.

## Concurrencia

- 50 o más requests concurrentes por la última unidad producen exactamente un éxito.
- 100 requests concurrentes por 10 unidades producen exactamente 10 éxitos y 90 rechazos.
- Ningún camino de test permite stock negativo.
- La suite backend debe validar estado final de base de datos, no solo conteos de response.

## Expiración TTL

- Una reserva permanece activa antes de 60 segundos.
- Una reserva expira a los 60 segundos o después si no fue confirmada ni liberada.
- Las reservas expiradas ya no aparecen en listas de reservas activas.
- La expiración devuelve stock exactamente una vez.
- Re-ejecutar lógica de expiración no cambia el stock de nuevo.

## Release Manual

- Liberar una reserva activa funciona correctamente.
- Liberar una reserva activa devuelve stock exactamente una vez.
- Liberar una reserva ya liberada es éxito no-op o response no-op documentado.
- Liberar una reserva ya expirada es éxito no-op o response no-op documentado.
- Liberar dos veces en concurrencia devuelve stock una sola vez.

## Idempotencia

- `POST /reservations` sin `Idempotency-Key` es rechazado.
- Misma idempotency key y mismo payload devuelve el mismo reservation ID y outcome.
- Misma idempotency key y mismo payload en paralelo decrementa stock una sola vez.
- Misma idempotency key y payload distinto es rechazado.
- Replay idempotente no crea una segunda reserva.

## UI

- Reserva exitosa muestra feedback visible de success.
- Cantidad inválida muestra feedback visible de validation.
- Stock insuficiente muestra feedback visible de error.
- Release exitoso muestra feedback visible.
- Reservas activas muestran item, cantidad, timer de expiración y acción release.
- Loading states se muestran para fetch de inventario y mutations.
- La UI actualiza stock después de reserve, release y expiración sin requerir refresh manual del browser.

## OpenAPI

- OpenAPI incluye todos los endpoints implementados.
- OpenAPI documenta `Idempotency-Key`.
- OpenAPI documenta responses exitosos y de error.
- Los error codes de OpenAPI coinciden con backend.

## Entrega

- El repositorio incluye source backend.
- El repositorio incluye source frontend.
- El repositorio incluye seed data PostgreSQL.
- El repositorio incluye README con estrategia de concurrencia y comandos de test.
- El repositorio incluye `spec-kit-notes.md` con comandos, supuestos, refinamientos y pivots.
- La implementación final mapea de vuelta a estos artefactos de especificación.

