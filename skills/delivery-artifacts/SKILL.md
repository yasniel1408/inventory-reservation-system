---
name: delivery-artifacts
description: "Usar para artefactos finales del challenge: OpenAPI contract, README.md, spec-kit-notes.md, seed data PostgreSQL, docker-compose o instrucciones de PostgreSQL, comandos de test, explicación de estrategia de concurrencia, LLM usado, chat history completo y checklist de entrega."
---

# Delivery Artifacts

Usar este skill para preparar lo que el evaluador revisa fuera del código.

## Obligatorio

- OpenAPI para endpoints implementados.
- README con estrategia de concurrencia y comandos de test.
- `spec-kit-notes.md` con comandos, supuestos, refinamientos y pivots.
- Seed data PostgreSQL.
- Chat history completo o ruta documentada.

## OpenAPI Mínimo

- `GET /items`
- `POST /reservations`
- `GET /reservations`
- `DELETE /reservations/{id}`
- Header `Idempotency-Key`.
- Error responses estables.

## README Debe Explicar

- Cómo correr backend.
- Cómo correr frontend.
- Cómo iniciar PostgreSQL y seed.
- Cómo correr tests.
- Cómo se evita oversell.
- Cómo funcionan TTL e idempotencia.
- Qué LLM se usó y por qué.
