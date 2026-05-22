# language: es
Característica: Idempotencia en reservas y liberaciones
  Como cliente frontend o consumidor API
  Quiero que los reintentos de red sean seguros
  Para evitar reservas duplicadas o devoluciones de stock duplicadas

  Antecedentes:
    Dado que el backend expone una API REST para reservas
    Y las solicitudes de creación de reserva usan el header "Idempotency-Key"

  Escenario: Reintentar creación de reserva con la misma clave y el mismo payload
    Dado que existe un artículo con 5 unidades disponibles
    Cuando envío dos solicitudes "POST /reservations" con el mismo "Idempotency-Key" y el mismo payload
    Entonces ambas respuestas representan la misma reserva
    Y ambas respuestas tienen el mismo identificador de reserva
    Y el stock se decrementa una sola vez

  Escenario: Procesar dos creaciones simultáneas con la misma clave idempotente
    Dado que existe un artículo con 5 unidades disponibles
    Cuando envío dos solicitudes paralelas "POST /reservations" con el mismo "Idempotency-Key" y el mismo payload
    Entonces se crea una sola reserva
    Y ambas solicitudes reciben el mismo resultado final
    Y el stock reservado aumenta solo por la cantidad de una reserva

  Escenario: Rechazar misma clave idempotente con payload diferente
    Dado que ya existe una solicitud registrada con "Idempotency-Key" igual a "abc-123" para reservar 1 unidad del artículo A
    Cuando envío otra solicitud "POST /reservations" con "Idempotency-Key" igual a "abc-123" para reservar 2 unidades del artículo A
    Entonces la solicitud es rechazada
    Y recibo un error claro de conflicto de idempotencia
    Y no se crea una nueva reserva
    Y el stock no cambia

  Escenario: Liberar la misma reserva más de una vez
    Dado que existe una reserva activa por 2 unidades
    Cuando envío dos solicitudes "DELETE /reservations/:id" para la misma reserva
    Entonces ambas solicitudes finalizan exitosamente o con una respuesta no-op bien definida
    Y el stock se devuelve exactamente una vez
    Y la reserva no queda activa después de la primera liberación efectiva

  Escenario: Reintentar liberación de una reserva ya expirada
    Dado que una reserva expiró y su stock ya fue devuelto
    Cuando envío una solicitud "DELETE /reservations/:id" para esa reserva
    Entonces recibo una respuesta no-op bien definida
    Y el stock no cambia
    Y no se recrea la reserva

