---
name: delivery-artifacts
description: "Usar para artefactos finales del challenge: OpenAPI contract, README.md, seed data PostgreSQL, docker-compose o instrucciones de PostgreSQL, comandos de test, explicación de estrategia de concurrencia, LLM usado, chat history completo y checklist de entrega."
---

# Delivery Artifacts

Usar este skill para preparar lo que el evaluador revisa fuera del código.

## Obligatorio

- OpenAPI para endpoints implementados.
- README con estrategia de concurrencia y comandos de test.
- Seed data PostgreSQL.
- Chat history completo o ruta documentada.

## PostgreSQL Local

- La imagen oficial de PostgreSQL solo ejecuta archivos directos en `/docker-entrypoint-initdb.d` durante la inicializacion de un volumen nuevo.
- Si el repo separa `db/migrations/` y `db/seeds/`, agregar un script directo en `db/init/` que ejecute ambos directorios en orden.
- Documentar que los scripts de init no vuelven a correr sobre un volumen ya inicializado; para re-bootstrap local se debe recrear el volumen.

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
