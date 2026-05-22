# Agentes del Proyecto

Estos agentes describen roles de trabajo para Codex dentro de este repo. No reemplazan los `skills/`; los `skills/` siguen siendo las reglas tecnicas y de arquitectura. El flujo de arranque vive en `skills/development-flow/SKILL.md`.

## Orden de Uso

1. `skills/development-flow`: se consulta primero y siempre para tareas de desarrollo.
2. `analyst`: valida coherencia del pedido, lo divide en piezas de negocio y genera un plan.
3. Aprobacion del usuario: el flujo se detiene hasta que el usuario apruebe el plan del `analyst`.
4. `team-leader`: con el plan aprobado, traduce negocio a tareas tecnicas y asigna ownership.
5. `developer`: implementa una tarea tecnica concreta usando `tdd-development` si escribe codigo productivo. Puede instanciarse N veces si hay scopes independientes.
6. `reviewer`: revisa cambios contra specs, plan, tasks y riesgos.
7. `tester`: crea, corrige o elimina tests segun comportamiento real.
8. `delivery-manager`: verifica README, OpenAPI, seeds, setup, comandos y checklist de entrega.
9. `skills-expert`: corre al final para mantener `skills/` actualizadas, descubribles y sin redundancia.

## Reglas

- Leer siempre `skills/` antes de ejecutar tareas del repo.
- Usar `skills/development-flow/SKILL.md` como unica fuente de verdad del flujo inicial.
- Para desarrollo, usar siempre los agentes.
- El `team-leader` no debe actuar hasta que el usuario apruebe el plan del `analyst`.
- El `team-leader` debe reportar status continuo del flujo: etapa, agente activo, tarea, bloqueos y siguiente paso.
- `skills-expert` debe revisar al cierre si los cambios exigen actualizar, fusionar o eliminar skills.
- Cada agente debe seleccionar y declarar las skills necesarias para su tarea actual antes de ejecutar.
- Todo `developer` debe incluir `tdd-development` cuando implementa feature, bugfix, refactor o cambio de comportamiento.
- Cada agente debe usar el perfil de modelo recomendado para su rol cuando la herramienta lo permita.
- El resumen final debe incluir que se hizo, skills aplicadas, agentes usados, camino tomado, validaciones y riesgos.
- Cada agente debe entregar `## Handoff` cuando pasa trabajo al siguiente agente.
- Paralelizar solo cuando los archivos o responsabilidades no se pisen.
- Toda salida debe ser accionable: archivos, decisiones, riesgos o validaciones.
- La fuente de verdad del challenge sigue siendo: `sdd/user_histories/`, `sdd/specs/001-inventory-reservation-system/`, `sdd/plans.md`/`sdd/plans/` y `sdd/tasks.md`/`sdd/tasks/`.

## Tareas Menores

No hace falta activar todo el flujo para preguntas conceptuales, lecturas cortas, comandos simples, correcciones menores o inspecciones rapidas. Si aparece riesgo, bug, cambio de arquitectura o implementacion real, se escala al flujo completo desde `analyst`.

## Perfiles de Modelo

| Agente | Perfil recomendado | Razonamiento |
| --- | --- | --- |
| `analyst` | Codex alto | high |
| `team-leader` | Codex maximo disponible, preferentemente `gpt-5.5-codex` si existe | xhigh |
| `developer` | Codex estandar | medium |
| `reviewer` | Codex alto | high |
| `tester` | Codex estandar | medium |
| `delivery-manager` | Codex estandar | medium |
| `skills-expert` | Codex alto | high |

Si la herramienta no soporta el nombre exacto, usar el modelo disponible mas cercano.
