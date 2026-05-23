# Checklist del Harness

Este checklist valida que el harness siga siendo portable, descubrible y util sin agregar ceremonia innecesaria.

## Fuente Canonica

- [ ] `AGENTS.md` existe y define reglas tool-agnostic.
- [ ] `HARNESS.md` existe y explica el mapa operativo.
- [ ] Wrappers como `CLAUDE.md` u `opencode.json` no duplican reglas globales.
- [ ] Si hay conflicto entre wrappers y `AGENTS.md`, gana `AGENTS.md`.

## Skills

- [ ] `skills/SELECTING_SKILLS.md` lista todas las skills activas.
- [ ] Cada carpeta `skills/<skill>/` tiene `SKILL.md`.
- [ ] Cada `SKILL.md` tiene frontmatter con `name` y `description`.
- [ ] Cada skill operativa tiene `agents/openai.yaml`.
- [ ] Las descripciones explican cuando usar la skill.
- [ ] No existen skills redundantes para la misma responsabilidad.

## Agentes

- [ ] `.agents/README.md` lista el orden del flujo.
- [ ] Cada `.agents/*.md` tiene `name`, `description`, `model_profile` y `reasoning`.
- [ ] Cada agente define `## Objetivo`, `## Modelo`, `## Entradas`, `## Proceso`, `## Salida Esperada`, `## Handoff` y `## Reglas`.
- [ ] El flujo incluye `delivery-manager` entre `tester` y `skills-expert`.
- [ ] Cada agente declara skills seleccionadas y perfil/modelo usado.
- [ ] Cada sub-agente define `## Contexto Aislado`.
- [ ] `team-leader` define `## Brief de Sub-agente`.

## SDD

- [ ] `sdd/user_histories/` existe.
- [ ] `sdd/specs/001-inventory-reservation-system/` existe.
- [ ] `sdd/plans.md` existe como indice.
- [ ] `sdd/plans/001-inventory-reservation-system.md` existe como plan activo.
- [ ] `sdd/tasks.md` existe como indice.
- [ ] `sdd/tasks/001-inventory-reservation-system.md` existe como tablero activo.
- [ ] `sdd/TRACEABILITY.md` conecta historias, specs, plan, tasks y validaciones.

## Memoria

- [ ] `memory/decisions.md` contiene decisiones estables.
- [ ] `memory/learnings.md` contiene aprendizajes reutilizables.
- [ ] `memory/progress.md` es snapshot actual y no log historico.
- [ ] `memory/current-task.md` es snapshot de una tarea tecnica activa y no backlog.
- [ ] Si una regla guia ejecucion futura, vive en `skills/`.
- [ ] Si un dato solo da contexto historico, puede vivir en `memory/`.

## Validacion

- [ ] `scripts/validate-harness.sh` existe.
- [ ] `scripts/validate-harness.sh` corre sin errores.
- [ ] `opencode.json` parsea como JSON valido.
- [ ] No quedan referencias a `plan.md` raiz.
- [ ] No quedan referencias a nombres de procesos eliminados.

## Cierre de Flujo

- [ ] El resumen final reporta archivos cambiados.
- [ ] El resumen final reporta skills aplicadas.
- [ ] El resumen final reporta agentes usados.
- [ ] El resumen final reporta perfiles/modelos usados o fallback.
- [ ] El resumen final reporta validaciones ejecutadas.
- [ ] El resumen final reporta riesgos o pendientes.
