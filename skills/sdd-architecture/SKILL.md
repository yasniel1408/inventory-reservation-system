---
name: sdd-architecture
description: Usar al crear, actualizar o validar artefactos SDD para este challenge de reservas de inventario, incluyendo user histories, spec.md, plans.md, plans/, tasks.md, tasks/, decisiones de arquitectura, supuestos y trazabilidad antes de implementación.
---

# SDD Architecture

Usar este skill para preservar el flujo SDD del challenge. No escribir código de aplicación hasta que existan `spec.md`, `plans.md`, `plans/`, `tasks.md` y `tasks/`, y hasta que mapeen con las historias de usuario.

## Entradas Requeridas

- Leer todos los archivos en `skills/` antes de iniciar cualquier tarea del repo.
- Usar `user_histories/` como fuente de comportamiento.
- Usar `specs/001-inventory-reservation-system/` como paquete actual de especificación.
- Usar el PDF original del challenge solo cuando requisitos no estén claros o se deba confirmar un artefacto faltante.

## Flujo

1. Identificar qué fase de artefacto pide el usuario: historias, spec, plan, tasks, implementación o delivery.
2. Mantener artefactos en orden cronológico:
   - `user_histories/`
   - `specs/001-inventory-reservation-system/spec.md`
   - `specs/001-inventory-reservation-system/domain-model.md`
   - `specs/001-inventory-reservation-system/api-spec.md`
   - `specs/001-inventory-reservation-system/acceptance-criteria.md`
   - `specs/001-inventory-reservation-system/test-spec.md`
   - `plans.md` como indice
   - `plans/001-inventory-reservation-system.md` como plan activo
   - `tasks.md` como indice
   - `tasks/001-inventory-reservation-system.md` como tablero activo
3. Documentar supuestos, ambigüedades y decisiones en el artefacto más cercano: `spec.md`, `plans/<id>.md`, `tasks/<id>.md` o `README.md`.
4. Mantener trazabilidad explícita: cada tarea de implementación debe apuntar a al menos una spec, criterio de aceptación o requisito de test.
5. Preferir artefactos pequeños y revisables antes que documentos narrativos amplios.
6. Si el usuario pide implementar antes de que existan artefactos requeridos, crear o actualizar primero los artefactos faltantes.
7. Al trabajar con agentes, cada agente debe declarar las skills seleccionadas para su rol y scope.

## Reglas de Artefactos

- `spec.md` define comportamiento de producto, scope, requisitos no funcionales, edge cases y supuestos.
- `domain-model.md` define entidades, estados, transiciones, invariantes, concurrencia y límites de tiempo.
- `api-spec.md` define la superficie REST prevista antes de generar el OpenAPI formal.
- `acceptance-criteria.md` define criterios observables de pass/fail.
- `test-spec.md` define escenarios de verificación requeridos.
- `plans.md` debe ser un indice corto de planes.
- `plans/<id>.md` debe elegir tecnologías y estrategias concretas de implementación.
- `tasks.md` debe ser un indice corto de tableros.
- `tasks/<id>.md` debe contener work items ejecutables con ownership claro de archivo/módulo.

## Checks Específicos del Challenge

- Confirmar comportamiento de reserva atómica antes de escribir código backend.
- Confirmar expiración TTL, release manual y carreras de idempotencia antes de schema work.
- Confirmar sincronización UI antes de state management frontend.
- Confirmar OpenAPI, seed data, README y chat history completo antes de delivery.
- Confirmar estructura repo y comandos ejecutables antes de scaffold de backend/frontend.

## Guía de Trabajo Paralelo

Dividir trabajo solo cuando los scopes de escritura sean independientes. Buenos splits:

- Plan backend y plan frontend.
- Contrato OpenAPI y planificación de schema DB.
- Plan de tests y README/delivery.

Evitar ediciones paralelas al mismo artefacto salvo que un agente tenga ownership claro de una sección.
