# Tarea Actual

Este archivo es un snapshot de la tarea tecnica activa. No es historial; se reemplaza cuando cambia la tarea en curso.

## Identificacion

- Fecha: 2026-05-22.
- Task activa: T-003 Crear estructura de carpetas del proyecto.
- Tablero: `sdd/tasks/001-inventory-reservation-system.md`.
- Plan asociado: `sdd/plans/001-inventory-reservation-system.md`.
- Estado: pendiente.

## Owner y Alcance

- Owner: raiz del repo.
- Archivos/carpetas esperadas:
  - `backend/`
  - `frontend/`
  - `db/migrations/`
  - `db/seeds/`
  - `openapi/`
  - `docs/`

## Riesgo y Ejecucion

- Risk level: bajo.
- Modelo sugerido: Codex rapido/estandar.
- Agentes requeridos: team-leader, developer, reviewer, skills-expert.
- Puede paralelizarse: no necesario; scaffold inicial chico.
- Criterio para escalar: si la estructura base implica cambiar stack, contratos, scripts de build o decisiones del plan.

## Contexto Necesario

- Skills base: `development-flow`, `sdd-architecture`.
- Si se escribe codigo productivo: agregar `tdd-development`.
- Referencia directa: `sdd/plans/001-inventory-reservation-system.md#entrega`.
- Antes de implementar, el `team-leader` debe confirmar ownership y si la tarea puede dividirse.

## Validacion Esperada

- La estructura de carpetas existe.
- `sdd/tasks/001-inventory-reservation-system.md` se actualiza si T-003 queda completa.
- `sdd/TRACEABILITY.md` se revisa si cambia cobertura.
- `scripts/validate-harness.sh` sigue pasando.

## Bloqueos

- Ninguno.

## Siguiente Paso

- Activar `team-leader` para traducir T-003 a una ejecucion tecnica acotada.
