---
name: development-flow
description: "Usar siempre al inicio de toda tarea de desarrollo del repo para activar el flujo de agentes: leer skills, usar analyst para plan aprobado por el usuario, luego team-leader, developers, reviewer, tester y skills-expert al cierre."
---

# Flujo de Desarrollo

Este repo debe priorizar entrega clara sobre ceremonia. Usar este skill como punto de entrada obligatorio para tareas de desarrollo.

## Reglas

- Leer todos los archivos bajo `skills/` antes de ejecutar una tarea.
- Usar `sdd-architecture` antes de escribir código de aplicación.
- Para desarrollo, usar siempre el flujo de agentes definido en `.agents/`.
- `analyst` siempre debe producir un plan y esperar aprobacion explicita del usuario antes de continuar.
- `team-leader` solo puede continuar despues de que el usuario apruebe el plan del `analyst`.
- `team-leader` debe mantener status continuo: etapa actual, agente activo, tarea en curso, bloqueos y siguiente paso.
- `skills-expert` debe ejecutarse al final para revisar si `skills/` debe actualizarse, fusionarse o limpiarse.
- Cada agente debe seleccionar y declarar las skills necesarias para su tarea actual antes de ejecutar.
- Cada agente debe usar el perfil de modelo recomendado en `.agents/` cuando la herramienta lo permita.
- Todo `developer` que escriba codigo productivo debe usar `tdd-development`.
- Si se detecta o corrige un bug, regresion o fallo de validacion con causa reusable, documentar la regla preventiva en `skills/`.
- Dividir trabajo solo cuando haya scopes independientes reales.
- Preferir pocos artefactos buenos sobre muchos documentos repetidos.

## Flujo Práctico

1. Verificar que el pedido sea coherente con el PDF y las specs.
2. Activar `analyst` para generar un plan de negocio/alcance.
3. Presentar el plan al usuario y esperar aprobacion explicita.
4. Con plan aprobado, activar `team-leader` para traducirlo a tareas tecnicas.
5. Hacer que `team-leader` reporte status al iniciar cada etapa y al cambiar de agente activo.
6. Cada agente declara skills seleccionadas y perfil/modelo usado antes de ejecutar.
7. Todo `developer` agrega `tdd-development` si implementa codigo productivo.
8. Instanciar uno o mas `developer` solo si hay ownership claro y scopes independientes.
9. Pasar cambios por `reviewer`.
10. Usar `tester` para crear, corregir o eliminar tests cuando el cambio sea testeable.
11. Activar `skills-expert` para mantener `skills/` actualizadas y descubribles, incluyendo aprendizajes de bugs y validaciones fallidas.
12. Ejecutar validacion posible.
13. Reportar resumen final con cambios, skills aplicadas, agentes usados, perfiles/modelos usados, camino tomado, validacion y riesgos restantes.

## Bucle del Agente

Cada agente debe operar en este ciclo:

```text
leer contexto
  -> entender tarea
  -> planificar siguiente paso
  -> ejecutar
  -> validar
  -> reportar status
  -> aprender si aplica
  -> decidir si termina o repite
```

Reglas:

- No ejecutar sin contexto suficiente.
- No pasar al siguiente agente si falta aprobacion, validacion u ownership.
- Si hay bloqueo, reportarlo y replanificar.
- Si hay aprendizaje reusable, marcarlo para `skills-expert`.
- Si el trabajo puede pausarse o cambiar de herramienta, actualizar `memory/progress.md`.

## Status Obligatorio

Durante una tarea de desarrollo, mostrar status recurrente al usuario con:

- etapa actual del flujo
- agente activo
- tarea en curso
- skills seleccionadas
- perfil/modelo usado
- decision o bloqueo relevante
- siguiente paso

## Resumen Final Obligatorio

Al cerrar una tarea, reportar:

- que se hizo
- archivos cambiados
- skills aplicadas
- agentes usados
- perfiles/modelos usados o fallback aplicado
- camino que tomo el flujo
- validaciones ejecutadas
- riesgos o pendientes restantes

## Gate de Implementación

Antes de código deben existir:

- `user_histories/`
- `specs/001-inventory-reservation-system/`
- `plans.md` como indice de planes
- `plans/001-inventory-reservation-system.md` como plan activo
- `tasks.md` como indice de tareas
- `tasks/001-inventory-reservation-system.md` como tablero activo

Si falta alguno y la tarea es implementación, crear primero el artefacto faltante.
