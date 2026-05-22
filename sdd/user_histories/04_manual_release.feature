# language: es
Característica: Liberación manual de reservas
  Como usuario comprador
  Quiero cancelar manualmente una reserva
  Para devolver stock cuando ya no deseo mantener los artículos reservados

  Antecedentes:
    Dado que existe una reserva creada por mi usuario
    Y la reserva pertenece a un artículo con stock controlado

  Escenario: Liberar una reserva activa
    Dado que la reserva está activa
    Cuando libero manualmente la reserva
    Entonces la operación finaliza correctamente
    Y la reserva deja de aparecer como activa
    Y el stock reservado del artículo disminuye por la cantidad de la reserva
    Y el stock disponible del artículo aumenta por la cantidad de la reserva

  Escenario: Liberar una reserva después de que el TTL haya vencido
    Dado que el temporizador de la interfaz está desincronizado
    Y la reserva ya venció por TTL en el backend
    Cuando libero manualmente la reserva
    Entonces recibo una respuesta exitosa o no-op bien definida
    Y no se devuelve stock más de una vez
    Y la reserva no vuelve a quedar activa

  Escenario: Mostrar feedback de liberación en la interfaz
    Dado que tengo una reserva activa visible
    Cuando presiono el botón de liberar reserva
    Entonces veo confirmación visible de liberación
    Y la reserva desaparece de mi lista de reservas activas
    Y el inventario se actualiza sin requerir recarga manual

