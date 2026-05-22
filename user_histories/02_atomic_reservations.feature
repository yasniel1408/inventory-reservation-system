# language: es
Característica: Reservas atómicas de inventario
  Como usuario comprador
  Quiero reservar unidades de un artículo de forma segura
  Para que mi reserva no se pierda ni provoque sobreventa bajo alta concurrencia

  Antecedentes:
    Dado que existe un artículo con stock total de 10 unidades
    Y no existen reservas activas para ese artículo

  Escenario: Reservar unidades disponibles correctamente
    Cuando solicito reservar 2 unidades del artículo
    Entonces se crea una reserva activa por 2 unidades
    Y el identificador de la reserva es único
    Y el stock reservado del artículo aumenta a 2
    Y el stock disponible del artículo baja a 8

  Escenario: Rechazar reserva por stock insuficiente
    Dado que el artículo tiene 1 unidad disponible
    Cuando solicito reservar 2 unidades del artículo
    Entonces la reserva es rechazada
    Y recibo un error claro de stock insuficiente
    Y el stock reservado no cambia
    Y el stock disponible no queda en negativo

  Escenario: Rechazar cantidad inválida
    Cuando solicito reservar 0 unidades del artículo
    Entonces la reserva es rechazada
    Y recibo un error claro de cantidad inválida
    Y el stock reservado no cambia

  Escenario: Evitar sobreventa con 50 solicitudes simultáneas para la última unidad
    Dado que el artículo tiene exactamente 1 unidad disponible
    Cuando 50 usuarios solicitan reservar 1 unidad simultáneamente
    Entonces exactamente 1 solicitud crea una reserva activa
    Y las 49 solicitudes restantes son rechazadas por conflicto o stock insuficiente
    Y el stock disponible final es 0
    Y el stock reservado final es 1
    Y nunca existe stock disponible negativo

  Escenario: Limitar 100 solicitudes concurrentes a 10 reservas exitosas
    Dado que el artículo tiene exactamente 10 unidades disponibles
    Cuando 100 usuarios solicitan reservar 1 unidad simultáneamente
    Entonces exactamente 10 solicitudes crean reservas activas
    Y exactamente 90 solicitudes son rechazadas por conflicto o stock insuficiente
    Y el stock disponible final es 0
    Y el stock reservado final es 10
    Y nunca existe stock disponible negativo

