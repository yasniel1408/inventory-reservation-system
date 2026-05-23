# Niveles de Riesgo

Los niveles de riesgo definen modelo recomendado, paralelismo, revision y validacion.

## Bajo

Ejemplos:

- documentacion;
- estructura de carpetas;
- formato;
- indices;
- cambios mecanicos sin comportamiento.

Reglas:

- Puede usar modelo rapido/chico.
- Puede paralelizarse si no pisa archivos.
- Validacion simple suele alcanzar.

## Medio

Ejemplos:

- endpoints simples;
- UI sin logica critica;
- tests unitarios o componentes;
- configuracion reversible.

Reglas:

- Puede usar modelo rapido/estandar si el brief es claro.
- Requiere validacion concreta.
- Reviewer puede ser focalizado.

## Alto

Ejemplos:

- concurrencia;
- idempotencia;
- TTL;
- release de stock;
- migraciones;
- contratos API;
- cambios de datos.

Reglas:

- Usar modelo alto o escalar desde sub-agente.
- Requiere reviewer fuerte.
- Requiere tests o validacion reproducible.
- Evitar paralelismo sobre el mismo contrato o modulo.

## Critico

Ejemplos:

- seguridad;
- datos sensibles;
- cambios irreversibles;
- destruccion o migracion destructiva de datos;
- cambio de arquitectura base.

Reglas:

- Usar modelo maximo disponible.
- Requiere aprobacion explicita del usuario si cambia alcance o decision estable.
- No paralelizar sin ownership extremadamente claro.
- Requiere plan, validacion y rollback o mitigacion.
