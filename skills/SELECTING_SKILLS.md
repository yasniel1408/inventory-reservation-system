# Selección de Skills

Usar `development-flow` primero para toda tarea no trivial. Luego cargar solo el skill específico necesario.

## Regla por Agente

Cada agente debe recolectar sus propias skills para la tarea actual:

1. Leer este archivo.
2. Leer `development-flow` si la tarea es no trivial o de desarrollo.
3. Elegir solo las skills necesarias para el scope actual.
4. Declarar `skills seleccionadas` en el status inicial.
5. Agregar otra skill solo si el scope cambia.
6. Reportar `skills aplicadas` en el cierre.

## Skills Activos

- `sdd-architecture`: historias, specs, plan, tasks, supuestos, decisiones y trazabilidad SDD.
- `backend-reservation-system`: backend Go, Gin, GORM, PostgreSQL, migraciones, concurrencia, TTL, release e idempotencia.
- `frontend-reservation-app`: frontend React + Vite + TypeScript, UI de inventario/reservas, timers, estado API y feedback.
- `tdd-development`: ciclo test-first obligatorio para developers antes de escribir codigo productivo.
- `delivery-artifacts`: OpenAPI, README, seed data, chat history, comandos de ejecución/test y checklist final.

## Regla de Austeridad

No crear más skills salvo que aparezca una responsabilidad estable y repetida que no encaje en estos cuatro.

## Mapa Rapido por Rol

- `analyst`: `development-flow`, `sdd-architecture`; agregar `delivery-artifacts` si analiza entrega.
- `team-leader`: `development-flow`, `sdd-architecture`; agregar backend, frontend o delivery segun tareas.
- `developer`: `development-flow`, `tdd-development` y la skill tecnica asignada.
- `reviewer`: `development-flow`, `sdd-architecture`, `tdd-development` y skills tecnicas segun diff.
- `tester`: `development-flow`, `tdd-development` y backend/frontend segun suite.
- `delivery-manager`: `development-flow`, `sdd-architecture` y `delivery-artifacts`.
- `skills-expert`: `development-flow`, `sdd-architecture` y toda skill afectada por el aprendizaje.
