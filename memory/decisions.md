# Decisiones Estables

## AGENTS.md como Fuente Canonica

- Decision: `AGENTS.md` es la fuente canonica tool-agnostic del harness.
- Motivo: Codex, Claude Code, OpenCode u otras herramientas deben compartir la misma forma de trabajo.
- Implicacion: wrappers como `CLAUDE.md` u `opencode.json` no deben duplicar reglas; deben apuntar a las rutas canonicas.

## HARNESS.md como Mapa Operativo

- Decision: `HARNESS.md` explica como se conectan `AGENTS.md`, `skills/`, `.agents/`, `memory/`, wrappers y SDD.
- Motivo: el harness debe ser entendible para cualquier herramienta o developer sin leer todos los archivos primero.
- Implicacion: `HARNESS.md` no reemplaza `AGENTS.md`; solo documenta el mapa.

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

## Stack del Challenge

- Decision: backend Go + Gin + GORM + PostgreSQL; frontend React + Vite + TypeScript.
- Motivo: alinear implementacion con el PDF del challenge y decisiones ya documentadas.
- Implicacion: no usar Next.js, Tailwind, CQRS formal ni arquitectura hexagonal formal en este challenge.
