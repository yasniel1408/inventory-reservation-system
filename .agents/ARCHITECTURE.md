# Arquitectura del Flujo de Agentes

Este documento muestra como interactuan el skill `development-flow` y los agentes definidos en `.agents/`. Los agentes coordinan el trabajo; los `skills/` siguen definiendo reglas tecnicas y de arquitectura.

## Flujo Principal

```mermaid
flowchart TD
    U["Usuario pide cambio o implementacion"] --> DF["skill development-flow<br/>arranque obligatorio"]
    DF --> S["Leer skills/"]
    S --> A["analyst<br/>analiza coherencia y arma plan"]
    A --> AP{"Usuario aprueba<br/>el plan?"}
    AP -- "No" --> AR["analyst ajusta plan"]
    AR --> AP
    AP -- "Si" --> TL["team-leader<br/>traduce a tareas tecnicas"]
    TL --> DQ{"Hay tareas<br/>independientes?"}
    DQ -- "Si" --> D1["developer #1"]
    DQ -- "Si" --> D2["developer #2"]
    DQ -- "Si" --> DN["developer #N"]
    DQ -- "No" --> D["developer unico"]
    D1 --> R["reviewer"]
    D2 --> R
    DN --> R
    D --> R
    R --> T{"Requiere tests?"}
    T -- "Si" --> TS["tester"]
    T -- "No" --> V["validacion final"]
    TS --> V
    V --> OUT["Reporte al usuario"]
```

## Gate de Aprobacion

```mermaid
sequenceDiagram
    participant U as Usuario
    participant DF as skill development-flow
    participant A as analyst
    participant TL as team-leader
    participant DEV as developer
    participant REV as reviewer
    participant TES as tester

    U->>DF: Solicita desarrollo
    DF->>DF: Lee skills y detecta flujo
    DF->>A: Pide analisis y plan
    A->>U: Presenta plan para aprobacion
    U-->>A: Aprueba o pide cambios
    alt Usuario aprueba
        A->>TL: Entrega plan aprobado
        TL->>DEV: Asigna tareas tecnicas
        DEV->>REV: Entrega cambios
        REV->>TES: Solicita tests si aplica
        TES->>REV: Reporta resultados
        REV->>U: Estado final y riesgos
    else Usuario pide cambios
        A->>A: Ajusta plan
        A->>U: Reenvia plan
    end
```

## Responsabilidades

```mermaid
flowchart LR
    A["analyst"] -->|negocio, coherencia, scope| P["Plan aprobado por usuario"]
    P --> TL["team-leader"]
    TL -->|tareas tecnicas, ownership, paralelismo| DEV["developer(s)"]
    DEV -->|codigo y validacion local| REV["reviewer"]
    REV -->|hallazgos y gaps| TES["tester"]
    TES -->|tests y resultados| REV
    REV -->|cierre tecnico| U["usuario"]
```

## Reglas de Control

- `skills/development-flow/SKILL.md` se usa siempre al inicio de tareas de desarrollo.
- `analyst` siempre produce un plan antes de ejecutar desarrollo.
- El plan del `analyst` requiere aprobacion explicita del usuario.
- `team-leader` no continua hasta que el usuario apruebe el plan.
- `developer` puede instanciarse N veces solo si los scopes no pisan archivos.
- `reviewer` revisa contra specs, plan, tasks y riesgos.
- `tester` crea, corrige o elimina tests segun comportamiento real.
- Si aparece conflicto entre agentes y `skills/`, ganan los `skills/`.
