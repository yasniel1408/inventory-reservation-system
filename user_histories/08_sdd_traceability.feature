# language: es
Característica: Trazabilidad SDD y entregables
  Como evaluador del challenge
  Quiero que las decisiones y tareas estén documentadas antes de codificar
  Para verificar que el proyecto siguió un enfoque Architecture First

  Escenario: Documentar requisitos y casos límite antes del código
    Cuando reviso los artefactos SDD del repositorio
    Entonces encuentro un `spec.md` que describe requisitos funcionales y no funcionales
    Y el `spec.md` identifica casos límite de concurrencia, TTL, idempotencia y desincronización de UI
    Y las decisiones ambiguas quedan registradas como supuestos o preguntas resueltas

  Escenario: Trazar plan técnico hacia tareas implementables
    Cuando reviso `plan.md` y `tasks.md`
    Entonces cada decisión arquitectónica relevante tiene tareas asociadas
    Y cada tarea puede mapearse a código, pruebas o documentación verificable
    Y no existen tareas vagas sin resultado observable

  Escenario: Incluir pruebas requeridas por el challenge
    Cuando reviso la suite de pruebas planificada o implementada
    Entonces existe una prueba Go para 50 o más reservas simultáneas sobre la última unidad
    Y existe una prueba Go donde 100 reservas concurrentes para 10 unidades producen 10 éxitos y 90 rechazos
    Y existe una prueba Go de idempotencia para reservas paralelas con la misma clave
    Y existe una prueba Go de idempotencia para liberar dos veces devolviendo stock una sola vez
    Y existe una prueba React para lógica de temporizador
    Y existen pruebas React de flujo exitoso de reserva y estado de error

  Escenario: Entregar documentación final del proyecto
    Cuando reviso los entregables del repositorio
    Entonces encuentro código fuente Go del backend
    Y encuentro código fuente React con Vite y TypeScript del frontend
    Y encuentro seed data para PostgreSQL
    Y encuentro un `README.md` con estrategia de concurrencia, ejecución de tests y LLM usado
    Y encuentro trazabilidad suficiente en `spec.md`, `plan.md`, `tasks.md` y `README.md`
