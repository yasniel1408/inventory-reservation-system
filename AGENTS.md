# AGENTS.md

Instrucciones obligatorias para Codex y agentes en este repositorio.

## Contexto

- Este repo implementa el challenge de reservas de inventario.
- La fuente principal de requerimientos esta en:
  - `user_histories/`
  - `specs/001-inventory-reservation-system/`
  - `plan.md`
  - `tasks.md`
- Los `skills/` definen reglas tecnicas y de arquitectura.
- `.agents/` define roles de coordinacion.

## Reglas Obligatorias

- Leer siempre todos los archivos dentro de `skills/` antes de ejecutar cualquier tarea del repo.
- Usar siempre `skills/development-flow/SKILL.md` al inicio de tareas de desarrollo.
- Para desarrollo, usar los roles definidos en `.agents/`.
- `analyst` debe generar un plan antes de cualquier implementacion.
- El plan del `analyst` requiere aprobacion explicita del usuario.
- `team-leader` no puede continuar hasta que el usuario apruebe el plan.
- `developer` solo trabaja sobre tareas tecnicas con ownership claro.
- `reviewer` revisa cambios contra specs, plan, tasks y riesgos.
- `tester` crea, corrige o elimina tests segun comportamiento real.
- Si hay conflicto entre `.agents/` y `skills/`, ganan los `skills/`.

## Flujo de Desarrollo

1. `skills/development-flow` arranca el flujo.
2. `analyst` analiza coherencia, scope y riesgos.
3. `analyst` presenta un plan al usuario.
4. Esperar aprobacion explicita del usuario.
5. `team-leader` traduce el plan aprobado a tareas tecnicas.
6. `developer` implementa una o mas tareas.
7. `reviewer` revisa el resultado.
8. `tester` ajusta tests cuando corresponde.
9. Reportar cambios, validaciones y riesgos restantes.

## Paralelismo

- Usar paralelismo solo cuando los scopes no pisan los mismos archivos.
- Cada developer debe tener objetivo, owner y validacion concretos.
- No ejecutar tareas grandes de forma monolitica si hay divisiones reales y seguras.

## Decisiones Tecnicas Vigentes

- Backend: Go + Gin + GORM + PostgreSQL.
- Frontend: React + Vite + TypeScript.
- No usar Next.js, Tailwind, CQRS formal ni arquitectura hexagonal formal en este challenge.
- PostgreSQL es la fuente de verdad para concurrencia, TTL, release e idempotencia.
- No usar locks en memoria para garantizar correctness.

## Entrega

- Mantener trazabilidad entre historias, specs, plan, tasks, implementacion y README.
- El README debe documentar setup, tests, estrategia de concurrencia, TTL, idempotencia y LLM usado.
- El contrato OpenAPI final debe vivir en `openapi/openapi.yaml`.
