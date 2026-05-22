# Memoria del Harness

Esta carpeta guarda memoria versionada dentro del repo. Es auditable y portable entre Codex, Claude Code, OpenCode u otra herramienta.

## Archivos

- `decisions.md`: decisiones estables que no queremos rediscutir sin motivo.
- `learnings.md`: aprendizajes reutilizables, especialmente bugs, fallos o ajustes de proceso.
- `progress.md`: snapshot del estado actual del flujo y agentes.

## Cuándo Leer

Leer esta carpeta cuando:

- una tarea toca el harness, agentes, skills o forma de trabajo;
- se retoma trabajo pausado y hace falta conocer el estado actual;
- una decision previa parece relevante;
- aparece un bug o fallo que podria tener una causa ya aprendida;
- hay duda entre una decision nueva y una decision historica.

No es obligatorio leer `memory/` para cada microtarea si `skills/`, specs y plan ya bastan.
Si existe duda sobre el estado actual del flujo, leer `memory/progress.md`.

## Cuándo Actualizar

Actualizar memoria cuando:

- se toma una decision estable del harness;
- se corrige un bug o fallo con causa reusable;
- se descarta una decision previa;
- se descubre un patron que conviene recordar, pero no necesariamente convertir en skill.
- cambia el agente activo, etapa, bloqueo o siguiente paso relevante del flujo.

`progress.md` debe actualizarse como snapshot: reemplazar estado actual, no acumular historial largo.

## Relacion con Skills

- Si el aprendizaje debe guiar ejecucion futura, debe ir a `skills/`.
- Si el aprendizaje solo explica contexto o trazabilidad, puede ir a `memory/`.
- Si algo esta en `skills/` y `memory/` contradice, gana `skills/`.
- `skills-expert` decide si un aprendizaje queda en `skills/`, en `memory/` o en ambos.
