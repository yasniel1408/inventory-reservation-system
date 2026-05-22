---
name: development-flow
description: "Usar al inicio de toda tarea no trivial del repo para mantener el flujo simple del challenge: leer skills, verificar coherencia, dividir solo si aporta valor, aplicar Spec Kit antes de código, implementar con ownership claro, revisar cambios y validar con tests/comandos concretos."
---

# Flujo de Desarrollo

Este repo debe priorizar entrega clara sobre ceremonia. Usar este skill como punto de entrada para tareas no triviales.

## Reglas

- Leer todos los archivos bajo `skills/` antes de ejecutar una tarea.
- Usar `spec-kit-architecture` antes de escribir código de aplicación.
- Dividir trabajo solo cuando haya scopes independientes reales.
- Evitar roles/agentes simulados si una tarea puede resolverse con una ejecución directa y validada.
- Preferir pocos artefactos buenos sobre muchos documentos repetidos.

## Flujo Práctico

1. Verificar que el pedido sea coherente con el PDF y las specs.
2. Identificar subtareas solo si reducen riesgo o permiten paralelismo real.
3. Elegir los skills mínimos desde `skills/SELECTING_SKILLS.md`.
4. Ejecutar cambios con ownership claro.
5. Revisar contra acceptance criteria y test spec.
6. Ejecutar validación posible.
7. Reportar cambios, validación y riesgos restantes.

## Gate de Implementación

Antes de código deben existir:

- `user_histories/`
- `specs/001-inventory-reservation-system/`
- `plan.md`
- `tasks.md`
- `spec-kit-notes.md`

Si falta alguno y la tarea es implementación, crear primero el artefacto faltante.
