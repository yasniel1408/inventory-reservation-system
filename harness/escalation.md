# Reglas de Escalamiento

Un sub-agente debe parar y devolver al `team-leader` cuando se cumpla cualquiera de estas condiciones.

## Escalar Por Contexto

- El brief es insuficiente o ambiguo.
- Faltan fuentes necesarias.
- La tarea no tiene owner claro.
- La validacion esperada no es concreta.

## Escalar Por Scope

- El scope crece.
- Hay que tocar archivos fuera del owner.
- Aparece una dependencia no prevista.
- Se necesita cambiar una task, spec, plan o contrato.

## Escalar Por Riesgo

- Cambia el risk level a alto o critico.
- Toca concurrencia, idempotencia, TTL, release, migraciones, contratos, seguridad o datos.
- Aparece un bug raro, no reproducible o con causa incierta.
- Falla una validacion clave.

## Escalar Por Modelo

- El sub-agente esta usando modelo rapido/chico y aparece ambiguedad.
- Se requiere decision arquitectonica.
- Se requiere razonamiento global entre backend, frontend, DB y entrega.

## Formato de Escalamiento

```md
## Escalamiento

- Para: team-leader
- Motivo:
- Risk level anterior:
- Risk level nuevo:
- Archivos afectados:
- Fuentes revisadas:
- Validacion fallida:
- Decision requerida:
- Recomendacion:
```
