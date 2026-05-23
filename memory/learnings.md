# Aprendizajes Reutilizables

## Diagramas Markdown

- Aprendizaje: algunos visores Markdown no renderizan Mermaid.
- Regla preventiva: cuando un documento debe ser legible en cualquier visor, incluir diagrama ASCII en bloque `text` y dejar Mermaid solo como opcional.
- Aplicado en: `.agents/ARCHITECTURE.md`.

## Wrappers de Herramientas

- Aprendizaje: no todos los asistentes leen los mismos archivos de instrucciones.
- Regla preventiva: mantener `AGENTS.md` como fuente canonica y crear wrappers minimos por herramienta.
- Aplicado en: `CLAUDE.md` y `opencode.json`.

## OpenCode debe Cargar la Fuente Canonica

- Aprendizaje: cargar solo `.agents/*.md` y `skills/*/SKILL.md` deja afuera reglas globales de `AGENTS.md`.
- Regla preventiva: wrappers/configs deben cargar la fuente canonica y los documentos de soporte necesarios, no solo roles o skills.
- Aplicado en: `opencode.json`.

## Fuentes Externas Fuera de Historias

- Aprendizaje: guardar PDFs o materiales fuente dentro de `sdd/user_histories/` ensucia la carpeta de historias Gherkin.
- Regla preventiva: poner fuentes externas en `references/` y dejar `sdd/user_histories/` solo para `.feature` y README.
- Aplicado en: movimiento del PDF del challenge a `references/`.

## Memoria vs Skills

- Aprendizaje: no todo aprendizaje debe convertirse en skill.
- Regla preventiva: si el aprendizaje guia acciones futuras, actualizar `skills/`; si solo da contexto historico, registrarlo en `memory/`.
- Responsable: `skills-expert`.

## Progreso no es Historial

- Aprendizaje: guardar cada evento de agentes en memoria puede convertir el repo en un log ruidoso.
- Regla preventiva: `memory/progress.md` debe ser snapshot del estado actual, no historial acumulativo.
- Responsable: `team-leader` durante ejecucion y `skills-expert` al cierre.

## Bucle Explícito Evita Ejecución a Ciegas

- Aprendizaje: si el agente no tiene un ciclo claro, puede ejecutar sin contexto, saltarse validaciones o no reportar estado.
- Regla preventiva: todo agente debe seguir el bucle contexto -> plan -> ejecutar -> validar -> status -> aprender -> repetir/cerrar.
- Aplicado en: `AGENTS.md`, `skills/development-flow/SKILL.md`, `.agents/ARCHITECTURE.md` y `HARNESS.md`.

## Evitar Duplicacion del Flujo

- Aprendizaje: duplicar `development-flow` en `.agents/` y `skills/` crea dos fuentes de verdad.
- Regla preventiva: mantener `skills/development-flow/SKILL.md` como unica fuente del flujo inicial; `.agents/` solo define roles.
- Aplicado en: eliminacion de `.agents/development-flow.md`.

## Evitar Dependencia de Marca en el Proceso

- Aprendizaje: nombrar el flujo con una marca especifica ata el harness a una herramienta concreta.
- Regla preventiva: usar terminos generales del proceso, como SDD, specs, plan y tasks.
- Aplicado en: usar `skills/sdd-architecture` y limpiar referencias textuales de marca.

## Init de PostgreSQL no Ejecuta Subdirectorios

- Aprendizaje: la imagen oficial de PostgreSQL ejecuta archivos directos en `/docker-entrypoint-initdb.d`, pero no recorre subdirectorios de migraciones o seeds automaticamente.
- Regla preventiva: si se mantienen `db/migrations/` y `db/seeds/`, agregar un script directo en `db/init/` que ejecute esos directorios en orden.
- Aplicado en: `db/init/001_run_db_scripts.sh` y mounts de `docker-compose.yml`.

## Idempotencia Debe Reproducir el Outcome Almacenado

- Aprendizaje: reconstruir un replay idempotente desde estado vivo rompe el contrato si la reserva luego fue liberada o expirada.
- Regla preventiva: guardar y devolver el `response_body` original para replay, y commitear outcomes determinísticos aun cuando el handler deba responder error.
- Aplicado en: `backend/internal/reservations/postgres_store.go` y `skills/backend-reservation-system/SKILL.md`.

## Tests DB con Apple Container

- Aprendizaje: Apple Container puede publicar PostgreSQL en una IP de contenedor mientras otro servicio local ocupa `127.0.0.1:5432`.
- Regla preventiva: los tests de integracion deben usar `TEST_DATABASE_URL` explicito, resetear schema desde migraciones y no depender de `docker-compose` ni de bootstrap por mounts.
- Aplicado en: `backend/internal/reservations/postgres_store_integration_test.go`.

## Frontend Vite/Vitest y Testing Library

- Aprendizaje: `vite.config.ts` con bloque `test` falla en `tsc -b` si `defineConfig` viene de `vite`.
- Regla preventiva: importar `defineConfig` desde `vitest/config` cuando el proyecto use Vitest integrado a Vite.
- Aplicado en: `frontend/vite.config.ts` y `skills/frontend-reservation-app/SKILL.md`.

## Tests React Aislados

- Aprendizaje: sin `cleanup()` entre tests, React Testing Library puede dejar DOM acumulado y producir queries duplicadas o falsos fallos.
- Regla preventiva: agregar `afterEach(cleanup + vi.clearAllMocks)` en suites de componentes.
- Aplicado en: `frontend/src/App.test.tsx` y `skills/frontend-reservation-app/SKILL.md`.

## CORS en Smoke Frontend Real

- Aprendizaje: los tests mockeados del frontend no detectan que Vite y Gin corren en origins distintos.
- Regla preventiva: cuando `VITE_API_BASE_URL` apunta a otro origin, validar `Access-Control-Allow-Origin` y preflight `OPTIONS` para `Content-Type`, `Idempotency-Key` y `X-Session-ID`.
- Aplicado en: `backend/internal/http/router.go`, `backend/internal/http/router_test.go`, `skills/backend-reservation-system/SKILL.md` y `skills/frontend-reservation-app/SKILL.md`.
