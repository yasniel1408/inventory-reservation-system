# language: es
Característica: Reservas atomicas de inventario
  Como usuario comprador
  Quiero reservar unidades de un articulo de forma segura
  Para que mi reserva no se pierda ni provoque sobreventa bajo alta concurrencia

  Antecedentes:
    Dado que existe un articulo con stock total de 10 unidades
    Y no existen reservas activas para ese articulo

  Escenario: Reservar unidades disponibles correctamente
    Cuando solicito reservar 2 unidades del articulo
    Entonces se crea una reserva activa por 2 unidades
    Y el identificador de la reserva es unico
    Y el stock reservado del articulo aumenta a 2
    Y el stock disponible del articulo baja a 8

  Escenario: Rechazar reserva por stock insuficiente
    Dado que el articulo tiene 1 unidad disponible
    Cuando solicito reservar 2 unidades del articulo
    Entonces la reserva es rechazada
    Y recibo un error claro de stock insuficiente
    Y el stock reservado no cambia
    Y el stock disponible no queda en negativo

  Escenario: Rechazar cantidad invalida
    Cuando solicito reservar 0 unidades del articulo
    Entonces la reserva es rechazada
    Y recibo un error claro de cantidad invalida
    Y el stock reservado no cambia

  Escenario: Evitar sobreventa con 50 solicitudes simultaneas para la ultima unidad
    Dado que el articulo tiene exactamente 1 unidad disponible
    Cuando 50 usuarios solicitan reservar 1 unidad simultaneamente
    Entonces exactamente 1 solicitud crea una reserva activa
    Y las 49 solicitudes restantes son rechazadas por conflicto o stock insuficiente
    Y el stock disponible final es 0
    Y el stock reservado final es 1
    Y nunca existe stock disponible negativo

  Escenario: Limitar 100 solicitudes concurrentes a 10 reservas exitosas
    Dado que el articulo tiene exactamente 10 unidades disponibles
    Cuando 100 usuarios solicitan reservar 1 unidad simultaneamente
    Entonces exactamente 10 solicitudes crean reservas activas
    Y exactamente 90 solicitudes son rechazadas por conflicto o stock insuficiente
    Y el stock disponible final es 0
    Y el stock reservado final es 10
    Y nunca existe stock disponible negativo

