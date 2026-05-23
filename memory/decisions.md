# Decisiones Estables

## AGENTS.md como Fuente Canonica

- Decision: `AGENTS.md` es la fuente canonica tool-agnostic del harness.
- Motivo: Codex, Claude Code, OpenCode u otras herramientas deben compartir la misma forma de trabajo.
- Implicacion: wrappers como `CLAUDE.md` u `opencode.json` no deben duplicar reglas; deben apuntar a las rutas canonicas.

## HARNESS.md como Mapa Operativo

- Decision: `HARNESS.md` explica como se conectan `AGENTS.md`, `skills/`, `.agents/`, `memory/`, wrappers y SDD.
- Motivo: el harness debe ser entendible para cualquier herramienta o developer sin leer todos los archivos primero.
- Implicacion: `HARNESS.md` no reemplaza `AGENTS.md`; solo documenta el mapa.

## Lenguaje SDD General

- Decision: usar SDD como lenguaje de especificacion general.
- Motivo: el harness apunta a SDD general y no debe depender de una marca o kit especifico.
- Implicacion: skills, tasks, histories y docs deben hablar de SDD, specs, plan y tasks.

## Skills como Reglas Operativas

- Decision: `skills/` contiene reglas tecnicas, arquitectura, criterios de ejecucion y aprendizajes accionables.
- Motivo: las skills deben ser descubribles por herramientas y agentes futuros.
- Implicacion: si un bug o decision genera una regla preventiva reusable, `skills-expert` debe evaluar actualizar una skill.

## Agents como Roles de Coordinacion

- Decision: `.agents/` describe roles de coordinacion, no reglas tecnicas principales.
- Motivo: separar roles de trabajo de instrucciones tecnicas evita duplicacion y drift.
- Implicacion: si hay conflicto entre `.agents/` y `skills/`, ganan `skills/`.

## Flujo con Aprobacion del Usuario

- Decision: `analyst` genera plan y el usuario debe aprobarlo antes de que `team-leader` continue.
- Motivo: evita implementaciones no alineadas con negocio o scope.
- Implicacion: tareas de desarrollo no pasan a implementacion sin aprobacion explicita.

## Status Continuo

- Decision: `team-leader` debe mantener status continuo y resumen final del flujo.
- Motivo: el usuario necesita ver etapa, agente activo, tarea, bloqueos y siguiente paso.
- Implicacion: el cierre debe incluir trabajo realizado, skills aplicadas, agentes usados, camino tomado, validaciones y riesgos.

## Progreso como Snapshot

- Decision: `memory/progress.md` guarda el estado actual del flujo y agentes como snapshot.
- Motivo: permite pausar y retomar con Codex, Claude Code, OpenCode u otra herramienta sin depender de memoria interna.
- Implicacion: no debe convertirse en log historico; se reemplaza el estado vigente.

## Tarea Actual Separada del Progreso

- Decision: `memory/current-task.md` guarda la tarea tecnica activa separada de `memory/progress.md`.
- Motivo: el flujo de agentes y la tarea tecnica cambian a ritmos distintos.
- Implicacion: `progress.md` muestra etapa/agente; `current-task.md` muestra task, owner, scope, validacion, bloqueo y siguiente paso.

## Bucle del Agente

- Decision: todo agente trabaja con un bucle explicito de contexto, plan, ejecucion, validacion, status y aprendizaje.
- Motivo: evita ejecucion a ciegas y hace portable el modo de trabajo entre herramientas.
- Implicacion: cada agente debe reportar status, validar antes de cerrar y marcar aprendizajes para `skills-expert`.

## Skills por Agente

- Decision: cada agente recolecta y declara las skills necesarias para su tarea actual.
- Motivo: reduce tokens y evita cargar contexto que no aplica sin perder reglas obligatorias.
- Implicacion: el status y el cierre deben mostrar skills seleccionadas/aplicadas.

## TDD Obligatorio para Developers

- Decision: todo agente `developer` que escriba codigo productivo debe usar `tdd-development`.
- Motivo: asegurar que los cambios de comportamiento tengan tests que fallan primero y evitar tests sesgados escritos despues.
- Implicacion: el resultado del developer debe incluir evidencia RED/GREEN o una excepcion aprobada.

## Perfiles de Modelo por Agente

- Decision: cada agente tiene un perfil de modelo recomendado; `team-leader` usa el modelo Codex maximo disponible, preferentemente `gpt-5.5-codex` si existe, con razonamiento `xhigh`.
- Motivo: concentrar mayor razonamiento en coordinacion global y ahorrar tokens en tareas acotadas.
- Implicacion: si una herramienta no soporta el modelo exacto, debe usar el modelo disponible mas cercano y reportar fallback.

## Delivery Manager Antes de Skills Expert

- Decision: agregar `delivery-manager` despues de `tester` y antes de `skills-expert`.
- Motivo: separar validacion de comportamiento de verificacion de entrega final.
- Implicacion: README, OpenAPI, seeds, setup, comandos y checklist final se revisan antes de cerrar y antes de actualizar skills.

## Harness Verificable

- Decision: agregar `harness/checklist.md`, `scripts/validate-harness.sh` y `sdd/TRACEABILITY.md`.
- Motivo: convertir reglas de harness y SDD en estructura verificable, no solo documentacion narrativa.
- Implicacion: cambios futuros del harness deben pasar por checklist, script y matriz de trazabilidad cuando apliquen.

## Handoff Entre Agentes

- Decision: cada agente entrega un bloque `## Handoff` al siguiente.
- Motivo: hacer portable el flujo entre Codex, Claude Code, OpenCode u otra herramienta.
- Implicacion: el siguiente agente recibe contexto, archivos, decisiones, validacion, riesgos y proximo paso sin reconstruir todo.

## Contexto Aislado de Sub-agentes

- Decision: los sub-agentes no heredan todo el contexto interno del `team-leader`.
- Motivo: reducir tokens, evitar contaminacion por dudas o ramas descartadas y mantener ownership claro.
- Implicacion: `team-leader` entrega un brief minimo y cada sub-agente lee por si mismo las fuentes requeridas.

## Stack del Challenge

- Decision: backend Go + Gin + GORM + PostgreSQL; frontend React + Vite + TypeScript.
- Motivo: alinear implementacion con el PDF del challenge y decisiones ya documentadas.
- Implicacion: no usar Next.js, Tailwind, CQRS formal ni arquitectura hexagonal formal en este challenge.
