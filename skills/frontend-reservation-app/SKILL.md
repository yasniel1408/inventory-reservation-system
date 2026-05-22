---
name: frontend-reservation-app
description: Usar para implementar, planificar o revisar el frontend React + Vite + TypeScript del challenge, incluyendo dashboard de inventario, reserve action, reservas activas, release button, timer de expiración, loading/error states, feedback de validación/conflicto, sincronización con backend por polling/refetch y tests de componentes.
---

# Frontend Reservation App

Usar este skill para todo el frontend. Seguir el PDF: React + Vite + TypeScript.

## Stack

- React.
- Vite.
- TypeScript.
- CSS simple y componentes propios; no introducir librerías visuales pesadas.

## UI Requerida

- Lista de items con nombre, stock total, stock disponible y acción de reserva.
- Feedback visible para success, stock insuficiente, cantidad inválida y release.
- Vista de reservas activas del usuario con botón release.
- Loading y error states.
- Timer de expiración que reconcilia con backend.

## Reglas de Estado

- Backend es canónico.
- Refetch después de reserve, release y expiración de timer.
- No crear reservas falsas ante errores.
- Generar una `Idempotency-Key` por intento de reserva y reutilizarla solo para retries del mismo intento.

## Tests Frontend Obligatorios

- Unit tests para lógica de timer.
- Component test de happy path de reserva.
- Component test de estado de error, por ejemplo stock insuficiente.
