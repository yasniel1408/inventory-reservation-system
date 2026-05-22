# Trazabilidad SDD

Este archivo conecta historias, especificaciones, plan, tareas y validaciones. No reemplaza los documentos fuente; sirve como indice para detectar huecos.

## Artefactos Base

| Tipo | Ruta |
| --- | --- |
| Historias | `sdd/user_histories/` |
| Specs | `sdd/specs/001-inventory-reservation-system/` |
| Plan activo | `sdd/plans/001-inventory-reservation-system.md` |
| Tasks activas | `sdd/tasks/001-inventory-reservation-system.md` |
| Skill SDD | `skills/sdd-architecture/SKILL.md` |

## Matriz Principal

| Area | Historia | Spec | Plan | Tasks | Validacion esperada |
| --- | --- | --- | --- | --- | --- |
| Inventario | `sdd/user_histories/01_inventory_dashboard.feature` | `spec.md`, `api-spec.md#get-items`, `acceptance-criteria.md#inventario` | `#backend`, `#frontend` | T-009, T-023 | Backend endpoint + UI dashboard + estado vacio/carga |
| Reserva atomica | `sdd/user_histories/02_atomic_reservations.feature` | `spec.md`, `domain-model.md`, `test-spec.md#concurrencia-ultima-unidad` | `#estrategia-de-concurrencia` | T-010, T-016, T-017 | Tests concurrentes sin oversell |
| TTL | `sdd/user_histories/03_reservation_ttl.feature` | `domain-model.md`, `test-spec.md`, `api-spec.md` | `#estrategia-de-ttl-y-release`, `#estado-frontend` | T-014, T-026, T-027 | Expiracion devuelve stock una vez + timer reconciliado |
| Release manual | `sdd/user_histories/04_manual_release.feature` | `api-spec.md#delete-reservationsid`, `test-spec.md#idempotencia-de-release` | `#estrategia-de-ttl-y-release` | T-013, T-019, T-025 | Release doble no duplica stock |
| Idempotencia | `sdd/user_histories/05_idempotency.feature` | `domain-model.md`, `api-spec.md#post-reservations`, `test-spec.md#idempotencia-de-reserva` | `#estrategia-de-idempotencia` | T-011, T-018 | Misma key no duplica reserva ni decremento |
| UI feedback | `sdd/user_histories/06_ui_feedback_and_state.feature` | `acceptance-criteria.md`, `test-spec.md` | `#frontend`, `#estado-frontend` | T-012, T-024, T-028, T-029 | Feedback visible, refetch y tests de componentes |
| OpenAPI | `sdd/user_histories/07_openapi_contract.feature` | `api-spec.md` | `#entrega` | T-030 | `openapi/openapi.yaml` documenta endpoints, headers y errores |
| Entrega y trazabilidad | `sdd/user_histories/08_sdd_traceability.feature` | `spec.md`, `acceptance-criteria.md`, `test-spec.md` | `#entrega` | T-031, T-032, T-033, T-034 | README, chat history, comandos y validacion final |
| Harness | `sdd/user_histories/08_sdd_traceability.feature` | `sdd/TRACEABILITY.md` | `sdd/plans/001-inventory-reservation-system.md` | T-001, T-002, T-003 | `scripts/validate-harness.sh` |

## Reglas

- Cada nueva tarea debe apuntar a historia, spec, plan o criterio de aceptacion.
- Cada nueva validacion debe quedar asociada a una tarea.
- Si una tarea no tiene fuente de requerimiento, debe revisarla `analyst`.
- Si una validacion revela una regla reusable, debe pasar por `skills-expert`.
- Mantener esta matriz corta; el detalle vive en specs, plan y tasks.
