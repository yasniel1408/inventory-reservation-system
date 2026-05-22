# language: es
Característica: Expiración automática de reservas
  Como sistema de inventario
  Quiero expirar reservas no confirmadas después de 60 segundos
  Para devolver stock al inventario disponible sin intervención manual

  Antecedentes:
    Dado que existe un artículo con stock total de 5 unidades
    Y existe una reserva activa por 2 unidades del artículo
    Y la reserva tiene una ventana de expiración de 60 segundos

  Escenario: Mantener reserva activa antes del TTL
    Cuando han pasado 59 segundos desde la creación de la reserva
    Entonces la reserva sigue activa
    Y el stock reservado del artículo sigue incluyendo esas 2 unidades
    Y el stock disponible no incluye esas 2 unidades

  Escenario: Expirar reserva automáticamente al cumplirse el TTL
    Cuando han pasado 60 segundos desde la creación de la reserva
    Entonces la reserva se elimina permanentemente como reserva interactuable
    Y el stock reservado del artículo disminuye en 2 unidades
    Y el stock disponible del artículo aumenta en 2 unidades

  Escenario: Impedir interacción con una reserva expirada
    Dado que la reserva ya expiró automáticamente
    Cuando intento consultar acciones disponibles sobre esa reserva
    Entonces la reserva no aparece como activa
    Y cualquier acción dependiente de reserva activa devuelve una respuesta bien definida de no encontrada o no operable

  Escenario: Sincronizar expiración visible en la interfaz
    Dado que tengo una reserva activa visible en la interfaz
    Cuando el temporizador llega a cero
    Entonces la reserva deja de mostrarse como activa
    Y el inventario visible se actualiza para reflejar el stock devuelto

