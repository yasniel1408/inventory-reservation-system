# language: es
Característica: Expiracion automatica de reservas
  Como sistema de inventario
  Quiero expirar reservas no confirmadas despues de 60 segundos
  Para devolver stock al inventario disponible sin intervencion manual

  Antecedentes:
    Dado que existe un articulo con stock total de 5 unidades
    Y existe una reserva activa por 2 unidades del articulo
    Y la reserva tiene una ventana de expiracion de 60 segundos

  Escenario: Mantener reserva activa antes del TTL
    Cuando han pasado 59 segundos desde la creacion de la reserva
    Entonces la reserva sigue activa
    Y el stock reservado del articulo sigue incluyendo esas 2 unidades
    Y el stock disponible no incluye esas 2 unidades

  Escenario: Expirar reserva automaticamente al cumplirse el TTL
    Cuando han pasado 60 segundos desde la creacion de la reserva
    Entonces la reserva se elimina permanentemente como reserva interactuable
    Y el stock reservado del articulo disminuye en 2 unidades
    Y el stock disponible del articulo aumenta en 2 unidades

  Escenario: Impedir interaccion con una reserva expirada
    Dado que la reserva ya expiro automaticamente
    Cuando intento consultar acciones disponibles sobre esa reserva
    Entonces la reserva no aparece como activa
    Y cualquier accion dependiente de reserva activa devuelve una respuesta bien definida de no encontrada o no operable

  Escenario: Sincronizar expiracion visible en la interfaz
    Dado que tengo una reserva activa visible en la interfaz
    Cuando el temporizador llega a cero
    Entonces la reserva deja de mostrarse como activa
    Y el inventario visible se actualiza para reflejar el stock devuelto

