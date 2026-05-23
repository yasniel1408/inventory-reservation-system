# Chat History

El historial completo esta en la conversacion de Codex donde se construyo este repo. Este archivo deja una version resumida y revisable para la entrega.

## Resumen del flujo

1. Se leyo el PDF del challenge y se transformo en historias de usuario en `sdd/user_histories/`.
2. Se crearon specs SDD en `sdd/specs/001-inventory-reservation-system/`.
3. Se definio un plan simple en `sdd/plans/001-inventory-reservation-system.md`.
4. Se creo el tablero ejecutable en `sdd/tasks/001-inventory-reservation-system.md`.
5. Se implemento estructura base, DB, backend, tests backend, frontend, tests frontend y entrega.
6. Se mantuvo trazabilidad en `sdd/TRACEABILITY.md`.
7. Se actualizaron `memory/` y `skills/` cuando aparecieron aprendizajes reutilizables.
8. Se ejecuto una fase final de smoke end-to-end con backend, frontend y PostgreSQL real; se corrigio CORS y un ajuste responsive.

## Decisiones principales

- No usar Next.js ni Tailwind porque el PDF pide React + Vite + TypeScript.
- No usar arquitectura hexagonal/CQRS formal para evitar complejidad innecesaria en un challenge corto.
- Usar PostgreSQL como fuente de verdad para stock, TTL e idempotencia.
- Evitar locks en memoria para correctness de concurrencia.
- Usar expiracion lazy en vez de worker/cron para mantener el scope simple.
- Documentar contrato final en OpenAPI.
- Habilitar CORS en backend para que Vite pueda consumir la API local desde navegador real.

## Evidencia principal

- Historias: `sdd/user_histories/`
- Specs: `sdd/specs/001-inventory-reservation-system/`
- Plan: `sdd/plans/001-inventory-reservation-system.md`
- Tasks: `sdd/tasks/001-inventory-reservation-system.md`
- Trazabilidad: `sdd/TRACEABILITY.md`
- OpenAPI: `openapi/openapi.yaml`
- README: `README.md`

## Validaciones esperadas

```bash
cd backend && go test ./...
cd frontend && npm test
cd frontend && npm run lint
cd frontend && npm run build
ruby -e "require 'yaml'; YAML.load_file('openapi/openapi.yaml'); puts 'OpenAPI YAML OK'"
./scripts/validate-harness.sh
```

Smoke adicional:

```bash
curl -i -H 'Origin: http://localhost:5173' http://localhost:8080/items
curl -i -X OPTIONS http://localhost:8080/reservations \
  -H 'Origin: http://localhost:5173' \
  -H 'Access-Control-Request-Method: POST' \
  -H 'Access-Control-Request-Headers: content-type,idempotency-key'
```
