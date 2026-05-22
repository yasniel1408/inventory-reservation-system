# language: es
Característica: Contrato OpenAPI de la API REST
  Como consumidor de la API
  Quiero contar con un contrato OpenAPI completo
  Para integrar y validar el backend de reservas de forma consistente

  Escenario: Documentar consulta de inventario
    Cuando reviso el contrato OpenAPI
    Entonces existe una operacion para listar articulos de inventario
    Y la respuesta documenta nombre, stock total, stock reservado y stock disponible
    Y se documentan respuestas de error comunes

  Escenario: Documentar creacion de reserva idempotente
    Cuando reviso el contrato OpenAPI
    Entonces existe la operacion "POST /reservations"
    Y se documenta el header obligatorio "Idempotency-Key"
    Y se documenta el payload con articulo y cantidad
    Y se documenta la respuesta exitosa con identificador de reserva, estado y expiracion
    Y se documentan errores por cantidad invalida, stock insuficiente y conflicto de idempotencia

  Escenario: Documentar liberacion idempotente de reserva
    Cuando reviso el contrato OpenAPI
    Entonces existe la operacion "DELETE /reservations/{id}"
    Y se documenta que la operacion es segura ante reintentos
    Y se documenta la respuesta para liberacion efectiva
    Y se documenta la respuesta no-op para reservas ya liberadas o expiradas

  Escenario: Documentar consulta de reservas activas del usuario
    Cuando reviso el contrato OpenAPI
    Entonces existe una operacion para listar reservas activas del usuario
    Y la respuesta incluye articulo, cantidad, estado y expiracion

