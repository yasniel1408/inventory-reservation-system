# language: es
Característica: Liberacion manual de reservas
  Como usuario comprador
  Quiero cancelar manualmente una reserva
  Para devolver stock cuando ya no deseo mantener los articulos reservados

  Antecedentes:
    Dado que existe una reserva creada por mi usuario
    Y la reserva pertenece a un articulo con stock controlado

  Escenario: Liberar una reserva activa
    Dado que la reserva esta activa
    Cuando libero manualmente la reserva
    Entonces la operacion finaliza correctamente
    Y la reserva deja de aparecer como activa
    Y el stock reservado del articulo disminuye por la cantidad de la reserva
    Y el stock disponible del articulo aumenta por la cantidad de la reserva

  Escenario: Liberar una reserva despues de que el TTL haya vencido
    Dado que el temporizador de la interfaz esta desincronizado
    Y la reserva ya vencio por TTL en el backend
    Cuando libero manualmente la reserva
    Entonces recibo una respuesta exitosa o no-op bien definida
    Y no se devuelve stock mas de una vez
    Y la reserva no vuelve a quedar activa

  Escenario: Mostrar feedback de liberacion en la interfaz
    Dado que tengo una reserva activa visible
    Cuando presiono el boton de liberar reserva
    Entonces veo confirmacion visible de liberacion
    Y la reserva desaparece de mi lista de reservas activas
    Y el inventario se actualiza sin requerir recarga manual

