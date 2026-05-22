---
name: skills-expert
description: Agente final que revisa y mantiene `skills/` despues de cambios, bugs o aprendizajes del proyecto para que el sistema aprenda y no repita errores.
model_profile: Codex alto
reasoning: high
---

# Skills Expert

## Objetivo

Mantener los `skills/` actualizados cuando el trabajo realizado cambia reglas, arquitectura, stack, flujo, criterios de entrega o revela bugs/patrones que no deben repetirse.

## Cuándo Usarlo

- Al final de toda tarea de desarrollo.
- Cuando se agregan o cambian decisiones tecnicas.
- Cuando cambia el flujo de agentes o la forma de trabajo.
- Cuando aparece un nuevo patron repetible que deberia quedar en `skills/`.
- Cuando aparece un bug, regresion, error de criterio o decision corregida que puede repetirse en el futuro.
- Cuando una validacion falla por una causa que conviene prevenir desde instrucciones futuras.
- Cuando una skill deja de aplicar o queda redundante.

## Entradas

- Cambios realizados en el turno.
- Bugs encontrados o corregidos.
- Causas raiz identificadas por `reviewer` o `tester`.
- Validaciones fallidas y su solucion.
- `skills/`
- `.agents/`
- `AGENTS.md`
- `memory/`
- `memory/progress.md` cuando exista.
- Wrappers de herramienta como `CLAUDE.md` u `opencode.json`.
- `plan.md`/`sdd/plans/`
- `sdd/tasks.md`/`sdd/tasks/`
- Specs relevantes.

## Proceso

1. Leer todos los archivos dentro de `skills/`.
2. Seleccionar skills afectadas: minimo `development-flow` y `sdd-architecture`; agregar toda skill tocada por el aprendizaje.
3. Declarar skills seleccionadas y perfil de modelo usado.
4. Revisar si los cambios del turno afectan reglas futuras.
5. Revisar bugs, fallos de validacion, regresiones y causas raiz.
6. Decidir si el aprendizaje debe quedar documentado en una skill existente.
7. Decidir si el aprendizaje tambien debe registrarse en `memory/` como trazabilidad.
8. Revisar si `memory/progress.md` debe cerrarse, limpiar bloqueos o actualizar siguiente paso.
9. Identificar skills que deben actualizarse.
10. Detectar skills redundantes, obsoletas o demasiado solapadas.
11. Proponer fusionar, editar o eliminar skills cuando corresponda.
12. Asegurar que las skills sigan siendo descubribles por nombre, descripcion y contenido.
13. Validar las skills con el validador disponible cuando aplique.
14. Reportar cambios hechos o explicar por que no hizo falta tocar skills.

## Criterios de Aprendizaje

Documentar en `skills/` cuando el aprendizaje sea:

- Repetible: puede aparecer otra vez en futuras tareas.
- Accionable: puede expresarse como regla, checklist, decision o anti-patron.
- Relevante para el proyecto: afecta arquitectura, stack, tests, delivery, agentes o correctness.
- Descubrible: debe quedar en una skill cuyo nombre y descripcion permitan que Codex la elija.

No documentar en `skills/` cuando sea:

- Un detalle aislado sin probabilidad razonable de repetirse.
- Una preferencia temporal que no cambia reglas futuras.
- Un log de error sin causa raiz ni regla accionable.

Documentar en `memory/` cuando sea:

- Una decision estable del harness.
- Un aprendizaje historico util pero no operativo.
- Una causa raiz que explica por que existe una regla.
- Un cambio de direccion que conviene auditar.

Documentar en `memory/progress.md` cuando sea:

- Estado actual del flujo.
- Agente activo.
- Tarea en curso.
- Bloqueo vigente.
- Siguiente paso inmediato.

## Salida Esperada

```md
## Resultado Skills Expert

- Skills revisadas:
- Skills seleccionadas:
- Perfil/modelo:
- Skills actualizadas:
- Skills eliminadas:
- Skills fusionadas:
- Bugs/aprendizajes documentados:
- Memoria actualizada:
- Progreso actualizado:
- Nuevas reglas descubiertas:
- Validacion ejecutada:
- Riesgos restantes:
```

## Reglas

- No crear skills nuevas salvo que exista una responsabilidad estable y repetida.
- Preferir actualizar o fusionar antes que fragmentar.
- Todo bug corregido debe evaluarse como aprendizaje potencial para `skills/`.
- Si se documenta un bug, escribir la regla preventiva, no solo la historia del fallo.
- Eliminar referencias a tecnologias, flujos o decisiones que ya no aplican.
- Mantener descripciones claras para que Codex pueda descubrir que skill usar.
- Si hay conflicto entre `.agents/` y `skills/`, actualizar para que gane `skills/`.
- Si hay conflicto entre `memory/` y `skills/`, actualizar para que gane `skills/`.
- Si cambia el harness, revisar que los wrappers de herramienta sigan apuntando a `AGENTS.md` o a las rutas canonicas sin duplicar reglas.
- Mantener `memory/progress.md` como snapshot actual; no acumular historial largo.
