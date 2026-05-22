# Progreso Actual

Este archivo es un snapshot del estado actual del flujo. No es un log historico; se actualiza reemplazando el estado vigente.

## Estado del Flujo

- Fecha: 2026-05-22.
- Etapa: ajuste de harness de agentes.
- Agente activo: skills-expert.
- Plan aprobado: si.
- Tarea actual: perfiles de modelo explicitados dentro de cada `.agents/*.md`.
- Bloqueos: ninguno.
- Siguiente paso: continuar implementacion desde `sdd/tasks/001-inventory-reservation-system.md` T-003; developers deben usar `tdd-development` si escriben codigo productivo.

## Agentes

- analyst: completo; valido separar SDD como el que y mantener harness como el como.
- team-leader: completo; limito el cambio a artefactos SDD.
- developer: completo; movio historias, specs, planes y tasks a `sdd/` y actualizo referencias.
- reviewer: completo; reviso rutas canonicas y separacion de responsabilidades.
- tester: completo; valida estructura, referencias y configuracion.
- skills-expert: completo; mantuvo modelos en frontmatter y tambien visibles en el cuerpo de cada agente.

## Ultimo Resumen

- Cambios realizados: se agrego `## Modelo` a cada agente y se corrigio la referencia canonica a `sdd/plans.md`.
- Skills aplicadas: `development-flow`, `sdd-architecture`.
- Agentes usados: `team-leader`, `developer`, `reviewer`, `tester`, `skills-expert`.
- Validaciones: lectura de skills; revision de rutas; validacion de JSON/config.
- Riesgos: mantener root harness y `sdd/` sincronizados cuando se creen nuevos artefactos.
