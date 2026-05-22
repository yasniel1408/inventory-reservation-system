# language: es
Característica: Feedback y estado de la interfaz
  Como usuario del frontend
  Quiero recibir feedback claro durante reservas, errores y liberaciones
  Para entender el estado real de mis acciones sin refrescar manualmente

  Antecedentes:
    Dado que estoy usando la aplicacion web de reservas
    Y el frontend consume el backend REST del sistema

  Escenario: Reserva exitosa desde la lista de inventario
    Dado que un articulo tiene stock disponible
    Cuando ingreso una cantidad valida
    Y ejecuto la accion de reservar
    Entonces veo un feedback visible de exito
    Y la reserva aparece en mi vista de reservas activas
    Y el stock disponible visible disminuye de acuerdo con la reserva

  Escenario: Error visible por stock tomado por otro usuario
    Dado que el dashboard muestra 1 unidad disponible de un articulo
    Y otro usuario reserva esa ultima unidad antes que yo
    Cuando intento reservar 1 unidad del mismo articulo
    Entonces veo un error claro indicando que el stock ya no esta disponible
    Y el dashboard actualiza el stock disponible
    Y no aparece una reserva nueva en mi lista activa

  Escenario: Error visible por cantidad invalida
    Cuando ingreso una cantidad invalida para reservar
    Entonces veo un mensaje de validacion visible
    Y la solicitud de reserva no se envia al backend

  Escenario: Estado de carga durante accion asincrona
    Cuando ejecuto una reserva o una liberacion
    Entonces veo un estado de carga asociado a esa accion
    Y se evita duplicar la accion accidentalmente mientras la solicitud esta pendiente

  Escenario: Vista de reservas activas con temporizador y boton de liberacion
    Dado que tengo reservas activas
    Cuando abro la vista de mis reservas
    Entonces cada reserva muestra el articulo reservado
    Y cada reserva muestra la cantidad reservada
    Y cada reserva muestra el tiempo restante antes de expirar
    Y cada reserva tiene una accion para liberarla manualmente

  Escenario: Mantener la interfaz sincronizada con el backend
    Dado que cambia el estado de stock por una reserva, liberacion o expiracion
    Cuando el frontend recibe la actualizacion por polling, refetch o mecanismo equivalente
    Entonces el dashboard refleja el stock actualizado
    Y la vista de reservas activas refleja solo reservas operables

