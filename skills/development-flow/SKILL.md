---
name: development-flow
description: "Usar siempre al inicio de toda tarea de desarrollo del repo para activar el flujo de agentes: leer skills, usar analyst para plan aprobado por el usuario, luego team-leader, developers, reviewer y tester segun corresponda."
---

# Flujo de Desarrollo

Este repo debe priorizar entrega clara sobre ceremonia. Usar este skill como punto de entrada obligatorio para tareas de desarrollo.

## Reglas

- Leer todos los archivos bajo `skills/` antes de ejecutar una tarea.
- Usar `spec-kit-architecture` antes de escribir código de aplicación.
- Para desarrollo, usar siempre el flujo de agentes definido en `.agents/`.
- `analyst` siempre debe producir un plan y esperar aprobacion explicita del usuario antes de continuar.
- `team-leader` solo puede continuar despues de que el usuario apruebe el plan del `analyst`.
- Dividir trabajo solo cuando haya scopes independientes reales.
- Preferir pocos artefactos buenos sobre muchos documentos repetidos.

## Flujo Práctico

1. Verificar que el pedido sea coherente con el PDF y las specs.
2. Activar `analyst` para generar un plan de negocio/alcance.
3. Presentar el plan al usuario y esperar aprobacion explicita.
4. Con plan aprobado, activar `team-leader` para traducirlo a tareas tecnicas.
5. Instanciar uno o mas `developer` solo si hay ownership claro y scopes independientes.
6. Pasar cambios por `reviewer`.
7. Usar `tester` para crear, corregir o eliminar tests cuando el cambio sea testeable.
8. Ejecutar validacion posible.
9. Reportar cambios, validacion y riesgos restantes.

## Gate de Implementación

Antes de código deben existir:

- `user_histories/`
- `specs/001-inventory-reservation-system/`
- `plan.md`
- `tasks.md`

Si falta alguno y la tarea es implementación, crear primero el artefacto faltante.
