# language: es
Característica: Dashboard de inventario
  Como usuario del sistema de reservas
  Quiero ver el estado actual del inventario
  Para decidir que articulos puedo reservar durante una venta flash

  Antecedentes:
    Dado que existen articulos cargados en el inventario
    Y cada articulo tiene un stock total mayor o igual a cero

  Escenario: Visualizar articulos con stock total, reservado y disponible
    Cuando ingreso al dashboard de inventario
    Entonces veo una lista de articulos
    Y cada articulo muestra su nombre
    Y cada articulo muestra su stock total
    Y cada articulo muestra sus unidades reservadas
    Y cada articulo muestra sus unidades disponibles calculadas como stock total menos unidades reservadas activas

  Escenario: Mostrar inventario vacio sin romper la interfaz
    Dado que no existen articulos cargados
    Cuando ingreso al dashboard de inventario
    Entonces veo un estado vacio claro
    Y no veo acciones de reserva disponibles

  Escenario: Mostrar estado de carga al consultar inventario
    Dado que la consulta de inventario esta en progreso
    Cuando ingreso al dashboard de inventario
    Entonces veo un estado de carga
    Y la interfaz evita acciones que dependan de datos aun no cargados

  Escenario: Mostrar error si el inventario no puede cargarse
    Dado que el backend no responde correctamente al consultar inventario
    Cuando ingreso al dashboard de inventario
    Entonces veo un mensaje de error comprensible
    Y puedo reintentar la carga del inventario

