# language: es
Característica: Feedback y estado de la interfaz
  Como usuario del frontend
  Quiero recibir feedback claro durante reservas, errores y liberaciones
  Para entender el estado real de mis acciones sin refrescar manualmente

  Antecedentes:
    Dado que estoy usando la aplicación web de reservas
    Y el frontend consume el backend REST del sistema

  Escenario: Reserva exitosa desde la lista de inventario
    Dado que un artículo tiene stock disponible
    Cuando ingreso una cantidad válida
    Y ejecuto la acción de reservar
    Entonces veo un feedback visible de éxito
    Y la reserva aparece en mi vista de reservas activas
    Y el stock disponible visible disminuye de acuerdo con la reserva

  Escenario: Error visible por stock tomado por otro usuario
    Dado que el dashboard muestra 1 unidad disponible de un artículo
    Y otro usuario reserva esa última unidad antes que yo
    Cuando intento reservar 1 unidad del mismo artículo
    Entonces veo un error claro indicando que el stock ya no está disponible
    Y el dashboard actualiza el stock disponible
    Y no aparece una reserva nueva en mi lista activa

  Escenario: Error visible por cantidad inválida
    Cuando ingreso una cantidad inválida para reservar
    Entonces veo un mensaje de validación visible
    Y la solicitud de reserva no se envía al backend

  Escenario: Estado de carga durante acción asíncrona
    Cuando ejecuto una reserva o una liberación
    Entonces veo un estado de carga asociado a esa acción
    Y se evita duplicar la acción accidentalmente mientras la solicitud está pendiente

  Escenario: Vista de reservas activas con temporizador y botón de liberación
    Dado que tengo reservas activas
    Cuando abro la vista de mis reservas
    Entonces cada reserva muestra el artículo reservado
    Y cada reserva muestra la cantidad reservada
    Y cada reserva muestra el tiempo restante antes de expirar
    Y cada reserva tiene una acción para liberarla manualmente

  Escenario: Mantener la interfaz sincronizada con el backend
    Dado que cambia el estado de stock por una reserva, liberación o expiración
    Cuando el frontend recibe la actualización por polling, refetch o mecanismo equivalente
    Entonces el dashboard refleja el stock actualizado
    Y la vista de reservas activas refleja solo reservas operables

