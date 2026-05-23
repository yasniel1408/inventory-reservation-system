# Tarea Actual

Este archivo es un snapshot de la tarea tecnica activa. No es historial; se reemplaza cuando cambia la tarea en curso.

## Identificacion

- Fecha: 2026-05-22.
- Task activa: T-016 Test de concurrencia para última unidad.
- Tablero: `sdd/tasks/001-inventory-reservation-system.md`.
- Plan asociado: `sdd/plans/001-inventory-reservation-system.md`.
- Estado: pendiente.

## Owner y Alcance

- Owner: `backend/internal/reservations/`.
- Archivos/carpetas esperadas:
  - `backend/internal/reservations/*_test.go`

## Riesgo y Ejecucion

- Risk level: alto.
- Modelo sugerido: Codex alto para tests concurrentes y estado final DB.
- Agentes requeridos: team-leader, developer, reviewer, tester, skills-expert.
- Puede paralelizarse: no sobre `internal/reservations/`; los tests de concurrencia pisan el mismo contrato.
- Criterio para escalar: si Docker/PostgreSQL no esta disponible, si aparecen flakiness o si falla atomicidad/idempotencia.

## Contexto Necesario

- Skills base: `development-flow`, `sdd-architecture`, `backend-reservation-system`.
- Si se escribe o corrige código productivo por fallos de test: agregar `tdd-development`.
- Referencias directas:
  - `sdd/specs/001-inventory-reservation-system/test-spec.md#concurrencia-última-unidad`
  - `skills/backend-reservation-system/SKILL.md`
  - `sdd/tasks/001-inventory-reservation-system.md#fase-4---tests-backend`
- Antes de implementar, el `team-leader` debe confirmar si hay PostgreSQL disponible o definir estrategia local de test.

## Validacion Esperada

- Test concurrente por ultima unidad demuestra exactamente 1 exito.
- El test verifica estado final de DB, no solo conteo de responses.
- `go test ./...` desde `backend/` pasa.
- `sdd/tasks/001-inventory-reservation-system.md` se actualiza si T-016 queda completa.
- `scripts/validate-harness.sh` sigue pasando.

## Bloqueos

- Ninguno.

## Siguiente Paso

- Activar `team-leader` para preparar brief de T-016 y ejecutar tests backend concurrentes.
