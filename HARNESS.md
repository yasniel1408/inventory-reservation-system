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
        user_histories/
        specs/
        plan.md
        tasks.md
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
user_histories -> specs -> plan.md -> tasks.md -> implementacion -> validacion
```

Harness Engineering organiza como trabajan las herramientas y agentes:

```text
AGENTS.md -> skills -> agents -> memory -> wrappers
```

Ambos conviven: SDD define que construir; el harness define como trabajar de forma consistente sin importar la herramienta.

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
