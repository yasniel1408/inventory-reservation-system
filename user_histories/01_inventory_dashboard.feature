# language: es
Característica: Dashboard de inventario
  Como usuario del sistema de reservas
  Quiero ver el estado actual del inventario
  Para decidir qué artículos puedo reservar durante una venta flash

  Antecedentes:
    Dado que existen artículos cargados en el inventario
    Y cada artículo tiene un stock total mayor o igual a cero

  Escenario: Visualizar artículos con stock total, reservado y disponible
    Cuando ingreso al dashboard de inventario
    Entonces veo una lista de artículos
    Y cada artículo muestra su nombre
    Y cada artículo muestra su stock total
    Y cada artículo muestra sus unidades reservadas
    Y cada artículo muestra sus unidades disponibles calculadas como stock total menos unidades reservadas activas

  Escenario: Mostrar inventario vacío sin romper la interfaz
    Dado que no existen artículos cargados
    Cuando ingreso al dashboard de inventario
    Entonces veo un estado vacío claro
    Y no veo acciones de reserva disponibles

  Escenario: Mostrar estado de carga al consultar inventario
    Dado que la consulta de inventario está en progreso
    Cuando ingreso al dashboard de inventario
    Entonces veo un estado de carga
    Y la interfaz evita acciones que dependan de datos aún no cargados

  Escenario: Mostrar error si el inventario no puede cargarse
    Dado que el backend no responde correctamente al consultar inventario
    Cuando ingreso al dashboard de inventario
    Entonces veo un mensaje de error comprensible
    Y puedo reintentar la carga del inventario

