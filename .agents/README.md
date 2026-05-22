# Agentes del Proyecto

Estos agentes describen roles de trabajo para Codex dentro de este repo. No reemplazan los `skills/`; los `skills/` siguen siendo las reglas tecnicas y de arquitectura. El flujo de arranque vive en `skills/development-flow/SKILL.md`.

## Orden de Uso

1. `skills/development-flow`: se consulta primero y siempre para tareas de desarrollo.
2. `analyst`: valida coherencia del pedido, lo divide en piezas de negocio y genera un plan.
3. Aprobacion del usuario: el flujo se detiene hasta que el usuario apruebe el plan del `analyst`.
4. `team-leader`: con el plan aprobado, traduce negocio a tareas tecnicas y asigna ownership.
5. `developer`: implementa una tarea tecnica concreta. Puede instanciarse N veces si hay scopes independientes.
6. `reviewer`: revisa cambios contra specs, plan, tasks y riesgos.
7. `tester`: crea, corrige o elimina tests segun comportamiento real.

## Reglas

- Leer siempre `skills/` antes de ejecutar tareas del repo.
- Usar `skills/development-flow/SKILL.md` como unica fuente de verdad del flujo inicial.
- Para desarrollo, usar siempre los agentes.
- El `team-leader` no debe actuar hasta que el usuario apruebe el plan del `analyst`.
- Paralelizar solo cuando los archivos o responsabilidades no se pisen.
- Toda salida debe ser accionable: archivos, decisiones, riesgos o validaciones.
- La fuente de verdad del challenge sigue siendo: `user_histories/`, `specs/001-inventory-reservation-system/`, `plan.md` y `tasks.md`.
