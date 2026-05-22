# Aprendizajes Reutilizables

## Diagramas Markdown

- Aprendizaje: algunos visores Markdown no renderizan Mermaid.
- Regla preventiva: cuando un documento debe ser legible en cualquier visor, incluir diagrama ASCII en bloque `text` y dejar Mermaid solo como opcional.
- Aplicado en: `.agents/ARCHITECTURE.md`.

## Wrappers de Herramientas

- Aprendizaje: no todos los asistentes leen los mismos archivos de instrucciones.
- Regla preventiva: mantener `AGENTS.md` como fuente canonica y crear wrappers minimos por herramienta.
- Aplicado en: `CLAUDE.md` y `opencode.json`.

## OpenCode debe Cargar la Fuente Canonica

- Aprendizaje: cargar solo `.agents/*.md` y `skills/*/SKILL.md` deja afuera reglas globales de `AGENTS.md`.
- Regla preventiva: wrappers/configs deben cargar la fuente canonica y los documentos de soporte necesarios, no solo roles o skills.
- Aplicado en: `opencode.json`.

## Fuentes Externas Fuera de Historias

- Aprendizaje: guardar PDFs o materiales fuente dentro de `user_histories/` ensucia la carpeta de historias Gherkin.
- Regla preventiva: poner fuentes externas en `references/` y dejar `user_histories/` solo para `.feature` y README.
- Aplicado en: movimiento del PDF del challenge a `references/`.

## Memoria vs Skills

- Aprendizaje: no todo aprendizaje debe convertirse en skill.
- Regla preventiva: si el aprendizaje guia acciones futuras, actualizar `skills/`; si solo da contexto historico, registrarlo en `memory/`.
- Responsable: `skills-expert`.

## Evitar Duplicacion del Flujo

- Aprendizaje: duplicar `development-flow` en `.agents/` y `skills/` crea dos fuentes de verdad.
- Regla preventiva: mantener `skills/development-flow/SKILL.md` como unica fuente del flujo inicial; `.agents/` solo define roles.
- Aplicado en: eliminacion de `.agents/development-flow.md`.
