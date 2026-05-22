# Harness Engineering

Este archivo describe el mapa operativo del proyecto. La fuente canonica sigue siendo `AGENTS.md`; este documento explica como se conectan las piezas.

## Mapa

```text
AGENTS.md
  fuente canonica tool-agnostic
  |
  +-- skills/
  |     reglas operativas, arquitectura y aprendizajes accionables
  |
  +-- .agents/
  |     roles de coordinacion: analyst, team-leader, developer, reviewer, tester, skills-expert
  |
  +-- memory/
  |     decisiones estables y aprendizajes historicos versionados
  |
  +-- wrappers
  |     CLAUDE.md
  |     opencode.json
  |
  +-- SDD
        sdd/user_histories/
        sdd/specs/
        sdd/plans.md -> sdd/plans/
        sdd/tasks.md -> sdd/tasks/
```

## Prioridad de Instrucciones

1. `skills/`
2. `AGENTS.md`
3. `.agents/`
4. `memory/`
5. wrappers de herramienta

Si hay conflicto entre `memory/` y `skills/`, ganan `skills/`.
Si hay conflicto entre wrappers y `AGENTS.md`, gana `AGENTS.md`.

## Relacion con SDD

SDD organiza el trabajo de producto:

```text
sdd/user_histories -> sdd/specs -> sdd/plans.md + sdd/plans -> sdd/tasks.md + sdd/tasks -> implementacion -> validacion
```

Harness Engineering organiza como trabajan las herramientas y agentes:

```text
AGENTS.md -> skills -> agents -> modelos -> memory -> wrappers
```

Ambos conviven: SDD define que construir; el harness define como trabajar de forma consistente sin importar la herramienta.

## Skills y Modelos por Agente

- Cada agente selecciona las skills necesarias para su tarea actual antes de ejecutar.
- Cada agente declara en status las `skills seleccionadas`.
- Los developers cargan `tdd-development` cuando escriben codigo productivo.
- Cada agente usa el perfil de modelo definido en `AGENTS.md` y `.agents/` cuando la herramienta lo permite.
- Si el modelo exacto no existe en la herramienta, se usa el modelo disponible mas cercano y se reporta fallback.
- El `team-leader` concentra el razonamiento mas alto porque coordina tradeoffs, ownership, paralelismo y cierre.

## Bucle del Agente

Todo agente sigue el mismo ciclo operativo:

```text
leer contexto -> entender tarea -> planificar -> ejecutar -> validar -> reportar -> aprender -> repetir o cerrar
```

Este bucle evita ejecucion a ciegas y deja trazabilidad en status, `memory/progress.md`, skills o memoria cuando corresponde.

## Wrappers

- `CLAUDE.md`: wrapper para Claude Code. Apunta a `AGENTS.md` y `skills/SELECTING_SKILLS.md`.
- `opencode.json`: wrapper para OpenCode. Carga `AGENTS.md`, `.agents/*.md`, `skills/*/SKILL.md` y `memory/*.md`.

Los wrappers no deben duplicar reglas. Si una regla aplica a todas las herramientas, va en `AGENTS.md` o `skills/`.

## Memoria

- `memory/decisions.md`: decisiones estables.
- `memory/learnings.md`: aprendizajes reutilizables y contexto historico.

Si el aprendizaje debe guiar ejecucion futura, `skills-expert` debe evaluar moverlo o copiarlo a `skills/` como regla preventiva.

## Fuente del Challenge

El PDF original del challenge vive en `references/`.
