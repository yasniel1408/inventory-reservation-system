# language: es
Característica: Contrato OpenAPI de la API REST
  Como consumidor de la API
  Quiero contar con un contrato OpenAPI completo
  Para integrar y validar el backend de reservas de forma consistente

  Escenario: Documentar consulta de inventario
    Cuando reviso el contrato OpenAPI
    Entonces existe una operación para listar artículos de inventario
    Y la respuesta documenta nombre, stock total, stock reservado y stock disponible
    Y se documentan respuestas de error comunes

  Escenario: Documentar creación de reserva idempotente
    Cuando reviso el contrato OpenAPI
    Entonces existe la operación "POST /reservations"
    Y se documenta el header obligatorio "Idempotency-Key"
    Y se documenta el payload con artículo y cantidad
    Y se documenta la respuesta exitosa con identificador de reserva, estado y expiración
    Y se documentan errores por cantidad inválida, stock insuficiente y conflicto de idempotencia

  Escenario: Documentar liberación idempotente de reserva
    Cuando reviso el contrato OpenAPI
    Entonces existe la operación "DELETE /reservations/{id}"
    Y se documenta que la operación es segura ante reintentos
    Y se documenta la respuesta para liberación efectiva
    Y se documenta la respuesta no-op para reservas ya liberadas o expiradas

  Escenario: Documentar consulta de reservas activas del usuario
    Cuando reviso el contrato OpenAPI
    Entonces existe una operación para listar reservas activas del usuario
    Y la respuesta incluye artículo, cantidad, estado y expiración

