#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$ROOT_DIR"

fail() {
  printf 'ERROR: %s\n' "$1" >&2
  exit 1
}

require_file() {
  [[ -f "$1" ]] || fail "falta archivo requerido: $1"
}

require_dir() {
  [[ -d "$1" ]] || fail "falta carpeta requerida: $1"
}

require_contains() {
  local file="$1"
  local pattern="$2"
  rg -q "$pattern" "$file" || fail "$file no contiene patron requerido: $pattern"
}

require_dir skills
require_dir .agents
require_dir memory
require_dir sdd
require_dir harness

require_file AGENTS.md
require_file HARNESS.md
require_file skills/SELECTING_SKILLS.md
require_file skills/development-flow/SKILL.md
require_file .agents/README.md
require_file .agents/ARCHITECTURE.md
require_file memory/decisions.md
require_file memory/learnings.md
require_file memory/progress.md
require_file memory/current-task.md
require_file harness/checklist.md
require_file harness/brief-template.md
require_file harness/status-template.md
require_file harness/final-summary-template.md
require_file harness/risk-levels.md
require_file harness/escalation.md
require_file sdd/TRACEABILITY.md
require_file sdd/plans.md
require_file sdd/plans/001-inventory-reservation-system.md
require_file sdd/tasks.md
require_file sdd/tasks/001-inventory-reservation-system.md

for skill_dir in skills/*; do
  [[ -d "$skill_dir" ]] || continue
  require_file "$skill_dir/SKILL.md"
  require_contains "$skill_dir/SKILL.md" '^name: '
  require_contains "$skill_dir/SKILL.md" '^description: '
  require_file "$skill_dir/agents/openai.yaml"
done

for agent_file in .agents/*.md; do
  base="$(basename "$agent_file")"
  [[ "$base" == "README.md" || "$base" == "ARCHITECTURE.md" ]] && continue
  require_contains "$agent_file" '^name: '
  require_contains "$agent_file" '^description: '
  require_contains "$agent_file" '^model_profile: '
  require_contains "$agent_file" '^reasoning: '
  require_contains "$agent_file" '^## Objetivo'
  require_contains "$agent_file" '^## Modelo'
  require_contains "$agent_file" '^## Entradas'
  require_contains "$agent_file" '^## Proceso'
  require_contains "$agent_file" '^## Salida Esperada'
  require_contains "$agent_file" '^## Handoff'
  require_contains "$agent_file" '^## Reglas'
done

require_contains AGENTS.md 'delivery-manager'
require_contains AGENTS.md 'Contexto de Sub-agentes'
require_contains skills/development-flow/SKILL.md 'delivery-manager'
require_contains skills/development-flow/SKILL.md 'Contexto de Sub-agentes'
require_contains skills/SELECTING_SKILLS.md 'delivery-manager'
require_contains .agents/ARCHITECTURE.md 'delivery-manager'
require_contains .agents/ARCHITECTURE.md 'Contexto de Sub-agentes'
require_contains .agents/ARCHITECTURE.md 'brief minimo, masticado y verificable'
require_contains .agents/team-leader.md 'Brief de Sub-agente'
require_contains .agents/team-leader.md 'Resumen masticado'
require_contains .agents/team-leader.md 'Modelo/perfil recomendado'
require_contains AGENTS.md 'Optimización de Contexto y Modelo'
require_contains skills/development-flow/SKILL.md 'Optimización de Contexto y Modelo'
for sub_agent in .agents/developer.md .agents/reviewer.md .agents/tester.md .agents/delivery-manager.md .agents/skills-expert.md; do
  require_contains "$sub_agent" '^## Contexto Aislado'
  require_contains "$sub_agent" 'No debe cargar todo desde cero'
done
require_contains sdd/TRACEABILITY.md 'T-003'
require_contains memory/progress.md 'memory/current-task.md'
require_contains memory/current-task.md '^## Identificacion'
require_contains memory/current-task.md '^## Owner y Alcance'
require_contains memory/current-task.md '^## Riesgo y Ejecucion'
require_contains memory/current-task.md '^## Validacion Esperada'
require_contains memory/current-task.md 'Risk level:'
require_contains memory/current-task.md 'Modelo sugerido:'
require_contains memory/current-task.md 'Agentes requeridos:'
require_contains memory/current-task.md 'Puede paralelizarse:'
require_contains memory/current-task.md 'Criterio para escalar:'
require_contains harness/brief-template.md '## Brief de Sub-agente'
require_contains harness/status-template.md '## Status'
require_contains harness/final-summary-template.md '## Resumen Final'
require_contains harness/risk-levels.md '## Alto'
require_contains harness/escalation.md '## Escalamiento'

node -e "const fs=require('fs'); JSON.parse(fs.readFileSync('opencode.json','utf8'));" >/dev/null

if rg -q 'Spec Kit|spec-kit|SpecKit|SPEC KIT|spec kit|08_spec_kit' . \
  --glob '!/.git/**' \
  --glob '!scripts/validate-harness.sh' \
  --glob '!harness/checklist.md'; then
  fail "hay referencias obsoletas a procesos eliminados"
fi

if rg -q '(^|[^s]/)plan\.md' . \
  --glob '!/.git/**' \
  --glob '!scripts/validate-harness.sh' \
  --glob '!harness/checklist.md'; then
  fail "hay referencias obsoletas a plan.md raiz"
fi

printf 'Harness OK\n'
