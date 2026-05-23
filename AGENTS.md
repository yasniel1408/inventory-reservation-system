# AGENTS.md

Instrucciones obligatorias para Codex y agentes en este repositorio.

## Contexto

- Este repo implementa el challenge de reservas de inventario.
- La fuente principal de requerimientos esta en:
  - `sdd/user_histories/`
  - `sdd/specs/001-inventory-reservation-system/`
  - `sdd/plans.md` como indice de planes
  - `sdd/plans/001-inventory-reservation-system.md` como plan activo
  - `sdd/tasks.md` como indice de tareas
  - `sdd/tasks/001-inventory-reservation-system.md` como tablero activo
- Los `skills/` definen reglas tecnicas y de arquitectura.
- `.agents/` define roles de coordinacion.
- `memory/` guarda decisiones y aprendizajes versionados del harness.
- `memory/progress.md` guarda el snapshot actual del flujo y agentes.
- `memory/current-task.md` guarda el snapshot de la tarea tecnica activa.
- `HARNESS.md` explica el mapa operativo entre harness, SDD, wrappers, skills y memoria.
- `harness/` contiene templates, niveles de riesgo, escalamiento y checklist operativo.
- `AGENTS.md` es la fuente canonica tool-agnostic del harness.
- `CLAUDE.md` y `opencode.json` son wrappers para herramientas especificas y no deben duplicar reglas.

## Reglas Obligatorias

- Leer siempre todos los archivos dentro de `skills/` antes de ejecutar cualquier tarea del repo.
- Usar siempre `skills/development-flow/SKILL.md` al inicio de tareas de desarrollo.
- Para desarrollo, usar los roles definidos en `.agents/`.
- `analyst` debe generar un plan antes de cualquier implementacion.
- El plan del `analyst` requiere aprobacion explicita del usuario.
- `team-leader` no puede continuar hasta que el usuario apruebe el plan.
- `team-leader` debe mantener status continuo del flujo: etapa actual, agente activo, tarea en curso, bloqueos y siguiente paso.
- `team-leader` puede actualizar `memory/progress.md` cuando cambie etapa, agente activo, bloqueo o siguiente paso.
- `developer` solo trabaja sobre tareas tecnicas con ownership claro.
- `reviewer` revisa cambios contra specs, plan, tasks y riesgos.
- `tester` crea, corrige o elimina tests segun comportamiento real.
- `delivery-manager` verifica README, OpenAPI, seeds, setup, comandos y checklist de entrega antes del cierre.
- `skills-expert` corre al final para mantener `skills/` actualizadas, descubribles y sin redundancia.
- Cada agente debe recolectar y declarar las skills necesarias para su tarea actual antes de ejecutar.
- Cada agente debe usar el perfil de modelo recomendado para su rol cuando la herramienta lo permita.
- Cada agente debe entregar un bloque `## Handoff` al pasar trabajo al siguiente agente.
- Los sub-agentes no deben heredar todo el contexto interno del `team-leader`; deben recibir un brief minimo y autosuficiente.
- Si aparece un bug, regresion o validacion fallida con causa reusable, `skills-expert` debe documentar la regla preventiva en `skills/`.
- Usar `memory/` para decisiones estables y aprendizajes historicos que no necesariamente son reglas operativas.
- Usar `memory/progress.md` para retomar el estado actual sin convertirlo en log historico.
- Usar `memory/current-task.md` para retomar la tarea tecnica activa sin mezclarla con el estado del flujo.
- Si hay conflicto entre `.agents/` y `skills/`, ganan los `skills/`.
- Si hay conflicto entre wrappers de herramienta y `AGENTS.md`, gana `AGENTS.md`.
- Si hay conflicto entre `memory/` y `skills/`, ganan los `skills/`.

## Flujo de Desarrollo

1. `skills/development-flow` arranca el flujo.
2. `analyst` analiza coherencia, scope y riesgos.
3. `analyst` presenta un plan al usuario.
4. Esperar aprobacion explicita del usuario.
5. `team-leader` traduce el plan aprobado a tareas tecnicas.
6. `team-leader` reporta status al iniciar cada etapa y cuando cambie el agente activo.
7. `developer` implementa una o mas tareas.
8. `reviewer` revisa el resultado.
9. `tester` ajusta tests cuando corresponde.
10. `delivery-manager` valida artefactos de entrega.
11. `skills-expert` revisa si hay que actualizar, fusionar o eliminar skills, incluyendo aprendizajes de bugs corregidos.
12. Reportar resumen final con trabajo realizado, skills aplicadas, agentes usados, camino tomado, validaciones y riesgos restantes.

## Cuándo No Usar Todo el Flujo

No hace falta activar todos los agentes para tareas menores:

- preguntas conceptuales;
- lectura o resumen corto;
- comandos simples;
- correcciones menores de texto o formato;
- inspecciones rapidas sin cambio tecnico.

Si una tarea menor descubre riesgo, bug, cambio de arquitectura o implementacion real, escalar al flujo completo desde `analyst`.

## Seleccion de Skills por Agente

Cada agente hace una seleccion local de skills para reducir contexto sin perder reglas importantes.

Proceso obligatorio:

1. Leer `skills/SELECTING_SKILLS.md`.
2. Leer `skills/development-flow/SKILL.md` si la tarea es no trivial o de desarrollo.
3. Elegir solo las skills especificas necesarias para el trabajo actual.
4. Declarar en el status inicial: `skills seleccionadas`.
5. Si durante el trabajo aparece otro scope, agregar la skill correspondiente y reportarlo.
6. En el cierre, reportar `skills aplicadas`.

| Agente | Skills minimas | Skills condicionales |
| --- | --- | --- |
| `analyst` | `development-flow`, `sdd-architecture` | `delivery-artifacts` si analiza entrega |
| `team-leader` | `development-flow`, `sdd-architecture` | backend/frontend/delivery segun tareas |
| `developer` | `development-flow`, `tdd-development` + skill tecnica asignada | `sdd-architecture` si toca artefactos SDD |
| `reviewer` | `development-flow`, `sdd-architecture` | `tdd-development` y backend/frontend/delivery segun diff |
| `tester` | `development-flow`, `tdd-development` | backend/frontend segun suite |
| `delivery-manager` | `development-flow`, `sdd-architecture`, `delivery-artifacts` | backend/frontend si verifica comandos especificos |
| `skills-expert` | `development-flow`, `sdd-architecture` | toda skill afectada por el aprendizaje |

## Perfiles de Modelo

Estos perfiles son una politica del harness para ahorrar tokens. Si la herramienta no soporta el nombre exacto, usar el modelo disponible mas cercano al perfil.

| Agente | Perfil recomendado | Razonamiento | Uso esperado |
| --- | --- | --- | --- |
| `analyst` | Codex alto | high | Coherencia, scope, riesgos y plan aprobable. |
| `team-leader` | Codex maximo disponible, preferentemente `gpt-5.5-codex` si existe | xhigh | Coordinacion global, paralelismo, ownership y decisiones de tradeoff. |
| `developer` | Codex rapido/estandar | medium | Implementacion acotada con brief masticado y archivos owner claros. |
| `reviewer` | Codex alto | high | Revision de bugs, riesgos, contratos y gaps. |
| `tester` | Codex rapido/estandar | medium | Tests concretos y validacion de comportamiento. |
| `delivery-manager` | Codex rapido/estandar | low/medium | Checklist de entrega, README, OpenAPI, setup y comandos. |
| `skills-expert` | Codex alto | high | Mantener skills sin duplicacion ni drift. |

Reglas:

- Subir temporalmente el perfil si hay concurrencia, seguridad, datos, migraciones o decisiones irreversibles.
- Bajar temporalmente el perfil para cambios mecanicos, indices o formato.
- Usar modelos rapidos/chicos para sub-agentes cuando el brief este masticado, el scope sea claro y la validacion sea concreta.
- Usar `harness/risk-levels.md` para decidir modelo, paralelismo y revision.
- Usar `harness/escalation.md` cuando un sub-agente deba parar y devolver al `team-leader`.
- No reducir contexto critico: specs, plan activo, task activa, diff y skill tecnica aplicable siempre se leen cuando corresponden.
- El resumen final debe indicar si se uso un perfil distinto al recomendado.

## Optimización de Contexto y Modelo

El `team-leader` concentra el razonamiento caro: lee contexto amplio, decide scope, define ownership y prepara briefs masticados y verificables. Los sub-agentes deben ejecutar con contexto reducido.

Un sub-agente puede usar modelo rapido/chico si:

- el objetivo esta claro;
- los archivos owner estan definidos;
- las fuentes relevantes ya estan resumidas en el brief;
- las restricciones estan explicitas;
- la validacion esperada es concreta;
- no toca concurrencia critica, migraciones destructivas, seguridad, datos sensibles ni decisiones arquitectonicas.

Debe escalar a modelo alto si:

- el brief es ambiguo o insuficiente;
- cambia un contrato;
- toca datos, concurrencia, migraciones o seguridad;
- aparecen fallos raros o no reproducibles;
- debe tomar una decision arquitectonica.

## Bucle del Agente

Cada agente debe trabajar en ciclos cortos y observables:

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

Reglas del bucle:

- No ejecutar sin contexto suficiente.
- No avanzar al siguiente agente si falta aprobacion, validacion o ownership.
- Si aparece un bloqueo, reportarlo y volver a planificar.
- Si aparece un aprendizaje reusable, marcarlo para `skills-expert`.
- Si el flujo puede pausarse, actualizar `memory/progress.md`.
- Si la tarea tecnica puede pausarse, actualizar `memory/current-task.md`.

## Handoff

Cada agente debe pasar un bloque reutilizable al siguiente:

```md
## Handoff

- Para:
- Contexto:
- Archivos tocados:
- Decisiones:
- Validacion:
- Riesgos:
- Proximo paso:
```

## Contexto de Sub-agentes

Cuando `team-leader` instancia o deriva trabajo a un sub-agente, no debe pasarle todo su contexto interno. Debe entregar un brief minimo:

- objetivo concreto;
- archivos owner;
- specs, tasks o criterios relevantes;
- skills requeridas;
- restricciones;
- validacion esperada;
- formato de handoff esperado.

El brief debe ser masticado: incluir resumen de lo relevante, decisiones ya tomadas, fuentes ya revisadas y fuentes que el sub-agente debe abrir solo si hay duda o si va a editar/validar esa superficie.

El sub-agente no debe cargar todo desde cero. Debe verificar solo lo necesario y no depender de memoria implicita, razonamiento privado o dudas internas del `team-leader`.

Usar `harness/brief-template.md` para el brief. Usar `harness/escalation.md` si el brief no alcanza o sube el riesgo.

## Status y Resumen

- Durante ejecucion, el usuario debe ver status recurrente y claro de que esta pasando.
- Usar `harness/status-template.md` para status recurrente.
- Usar `harness/final-summary-template.md` para cierres.
- Si el flujo puede pausarse o cambiar de herramienta, reflejar el estado actual en `memory/progress.md`.
- Si la tarea tecnica puede pausarse o cambiar de herramienta, reflejar tarea, owner, scope, bloqueo y siguiente paso en `memory/current-task.md`.
- El status debe indicar:
  - etapa actual del flujo
  - agente activo
  - tarea en curso
  - skills seleccionadas
  - modelo o perfil usado
  - decision o bloqueo relevante
  - siguiente paso
- El resumen final debe incluir:
  - que se hizo
  - que archivos cambiaron
  - que skills se aplicaron
  - que agentes participaron
  - que perfiles de modelo se usaron o si hubo fallback
  - que camino tomo el flujo
  - que validaciones se ejecutaron
  - que riesgos o pendientes quedan

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

- Mantener trazabilidad entre historias, specs, `sdd/plans.md`, `sdd/plans/`, `sdd/tasks.md`, `sdd/tasks/`, implementacion y README.
- El README debe documentar setup, tests, estrategia de concurrencia, TTL, idempotencia y LLM usado.
- El contrato OpenAPI final debe vivir en `openapi/openapi.yaml`.
