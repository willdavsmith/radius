# Tasks: Automatic Application Discovery

**Input**: Design documents from `/specs/001-auto-app-discovery/`  
**Prerequisites**: plan.md ✓, spec.md ✓, research.md ✓, data-model.md ✓, contracts/ ✓

---

## Format: `[ID] [P?] [Story?] Description`

- **[P]**: Can run in parallel (different files, no dependencies)
- **[Story]**: Which user story this task belongs to (e.g., US1, US2)
- Include exact file paths in descriptions

## User Stories Summary

| ID | Priority | Title | Skills |
|----|----------|-------|--------|
| US1 | P1 | Analyze Existing Application Codebase | `discover_dependencies`, `discover_services` |
| US2 | P1 | Map to Resource Types | `generate_resource_types` |
| US4 | P1 | Generate Application Definition | `generate_app_definition`, `validate_app_definition` |
| US6 | P1 | AI Coding Agent Integration | MCP server, all skills |
| US3 | P2 | Match Recipes from Configured Sources | `discover_recipes` |
| US7 | P2 | Apply Team Infrastructure Practices | `discover_team_practices` |
| US5 | P3 | New Application Scaffolding | `rad app scaffold` |

---

## Phase 1: Setup (Shared Infrastructure)

**Purpose**: Project initialization and core abstractions

- [x] T001 Create discovery package structure in pkg/cli/discovery/
- [x] T002 [P] Add new dependencies to go.mod: mcp-go, go-openai, hcl/v2
- [x] T003 [P] Create Skill interface and registry in pkg/cli/discovery/skills/skill.go
- [x] T004 [P] Create skill registry for CLI and MCP in pkg/cli/discovery/skills/registry.go
- [x] T005 Create LLM provider interface in pkg/cli/discovery/llm/provider.go
- [x] T006 [P] Implement mock LLM provider for testing in pkg/cli/discovery/llm/mock.go
- [x] T007 Create data model types in pkg/cli/discovery/types.go (DiscoveryResult, DetectedDependency, DetectedService, etc.)
- [x] T008 [P] Create error types for discovery in pkg/cli/discovery/errors.go

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: Core infrastructure that MUST be complete before ANY user story can be implemented

**⚠️ CRITICAL**: No user story work can begin until this phase is complete

- [x] T009 Implement OpenAI LLM provider in pkg/cli/discovery/llm/openai.go
- [x] T010 [P] Implement Ollama LLM provider in pkg/cli/discovery/llm/ollama.go
- [x] T011 Create LLM prompt templates in pkg/cli/discovery/llm/prompts/
- [x] T012 Extend config.yaml schema for discovery settings in pkg/cli/config/discovery_config.go
- [x] T013 [P] Create discovery output writer for discovery.md in pkg/cli/discovery/output.go
- [x] T014 Add unit tests for skill registry in pkg/cli/discovery/skills/registry_test.go
- [x] T015 [P] Add unit tests for LLM providers in pkg/cli/discovery/llm/provider_test.go
- [x] T016 Create test fixtures directory structure in test/functional-portable/discovery/testdata/

**Checkpoint**: Foundation ready - user story implementation can now begin ✓

---

## Phase 3: User Story 1 - Analyze Existing Application Codebase (Priority: P1) 🎯 MVP

**Goal**: Detect infrastructure dependencies and deployable services from codebase

**Independent Test**: Point Radius at sample codebase with known dependencies (Node.js with PostgreSQL, Redis) and verify correct identification

### Implementation for User Story 1

- [x] T017 [US1] Create dependency detector skill in pkg/cli/discovery/dependencies/detector.go
- [x] T018 [P] [US1] Create dependency detection prompts in pkg/cli/discovery/llm/prompts/dependencies.txt
- [x] T019 [US1] Implement file walker with exclusion patterns in pkg/cli/discovery/dependencies/walker.go
- [x] T020 [US1] Create service detector skill in pkg/cli/discovery/services/detector.go
- [x] T021 [P] [US1] Create service detection prompts in pkg/cli/discovery/llm/prompts/services.txt
- [x] T022 [US1] Implement entrypoint detection (Dockerfile, package.json, main files) in pkg/cli/discovery/services/entrypoints.go
- [x] T023 [US1] Create `rad app discover` command in pkg/cli/cmd/app/discover/discover.go
- [x] T024 [US1] Wire discover command into root.go in cmd/rad/cmd/root.go
- [x] T025 [US1] Implement discovery.md output generation in pkg/cli/discovery/output.go
- [x] T026 [US1] Add progress indicators for long-running operations (NFR-02)
- [x] T027 [US1] Add partial failure handling with warnings (FR-10)
- [x] T028 [US1] Add unit tests for dependency detector in pkg/cli/discovery/dependencies/detector_test.go
- [x] T029 [P] [US1] Add unit tests for service detector in pkg/cli/discovery/services/detector_test.go
- [x] T030 [P] [US1] Create Node.js Express test fixture in test/functional-portable/discovery/testdata/nodejs-express/
- [x] T031 [P] [US1] Create Python Flask test fixture in test/functional-portable/discovery/testdata/python-flask/
- [x] T032 [P] [US1] Create Go Gin test fixture in test/functional-portable/discovery/testdata/go-gin/
- [x] T033 [US1] Add integration test for discover command in test/functional-portable/discovery/discover_test.go

**Checkpoint**: `rad app discover .` works and produces discovery.md with dependencies and services ✓

---

## Phase 4: User Story 2 - Map to Resource Types (Priority: P1)

**Goal**: Map detected dependencies to Resource Type definitions from catalog or generate minimal types

**Independent Test**: Verify that after dependency detection, valid Resource Type definitions are available

### Implementation for User Story 2

- [x] T034 [US2] Create Resource Type catalog in pkg/cli/discovery/resourcetypes/catalog.go
- [x] T035 [US2] Define catalog entries for common dependencies (PostgreSQL, Redis, MongoDB, etc.) in pkg/cli/discovery/resourcetypes/catalog_entries.go
- [x] T036 [US2] Create resource type generator skill in pkg/cli/discovery/resourcetypes/generator.go
- [x] T037 [US2] Implement catalog lookup and matching logic in pkg/cli/discovery/resourcetypes/matcher.go
- [x] T038 [US2] Implement fallback generation for unknown types in pkg/cli/discovery/resourcetypes/fallback.go
- [x] T039 [US2] Register generate_resource_types skill in pkg/cli/discovery/skills/builtin.go
- [x] T040 [US2] Add unit tests for catalog in pkg/cli/discovery/resourcetypes/catalog_test.go
- [x] T041 [P] [US2] Add unit tests for generator in pkg/cli/discovery/resourcetypes/generator_test.go

**Checkpoint**: Detected dependencies are mapped to valid Resource Type definitions ✓

---

## Phase 5: User Story 4 - Generate Application Definition (Priority: P1)

**Goal**: Generate complete app.bicep from discovery results

**Independent Test**: Generate app definition from sample codebase and verify it can be deployed with Radius

### Implementation for User Story 4

- [x] T042 [US4] Create Bicep template for application in pkg/cli/discovery/generator/templates/app.bicep.tmpl
- [x] T043 [P] [US4] Create Bicep template for containers in pkg/cli/discovery/generator/templates/container.bicep.tmpl
- [x] T044 [P] [US4] Create Bicep template for resources in pkg/cli/discovery/generator/templates/resource.bicep.tmpl
- [x] T045 [US4] Create Bicep generator skill in pkg/cli/discovery/generator/bicep_generator.go
- [x] T046 [US4] Implement container-to-resource connection wiring in pkg/cli/discovery/generator/connections.go
- [x] T047 [US4] Create Bicep validator skill in pkg/cli/discovery/generator/validator.go
- [x] T048 [US4] Create `rad app generate` command in pkg/cli/cmd/app/generate/generate.go
- [x] T049 [US4] Wire generate command into root.go in cmd/rad/cmd/root.go
- [x] T050 [US4] Implement --accept-defaults flag for non-interactive mode (FR-32)
- [x] T051 [US4] Implement --output flag for custom output paths (FR-34)
- [x] T052 [US4] Implement existing app.bicep handling: prompt for overwrite/merge/diff/cancel (FR-21)
- [x] T053 [US4] Register generate_app_definition and validate_app_definition skills
- [x] T054 [US4] Add unit tests for Bicep generator in pkg/cli/discovery/generator/bicep_generator_test.go
- [x] T055 [P] [US4] Add unit tests for validator in pkg/cli/discovery/generator/validator_test.go
- [x] T056 [US4] Add integration test for generate command in test/functional-portable/discovery/generate_test.go

**Checkpoint**: `rad app generate` produces valid app.bicep that can be deployed ✓

---

## Phase 6: User Story 6 - AI Coding Agent Integration (Priority: P1)

**Goal**: Expose all skills via MCP server for AI agent integration

**Independent Test**: Invoke each skill via MCP and verify structured JSON output

### Implementation for User Story 6

- [x] T057 [US6] Create MCP parent command in pkg/cli/cmd/mcp/mcp.go
- [x] T058 [US6] Create `rad mcp serve` command in pkg/cli/cmd/mcp/serve/serve.go
- [x] T059 [US6] Implement MCP request handler in pkg/cli/cmd/mcp/serve/handler.go
- [x] T060 [US6] Implement stdio transport for VS Code integration (FR-30)
- [x] T061 [US6] Implement HTTP transport for remote agents (FR-30)
- [x] T062 [US6] Wire MCP command into root.go in cmd/rad/cmd/root.go
- [x] T063 [US6] Add MCP server configuration flags (port, allowed origins, auth mode) (FR-31)
- [x] T064 [US6] Register all skills as MCP tools in pkg/cli/cmd/mcp/serve/tools.go
- [x] T065 [US6] Ensure CLI commands use same skill implementations (FR-27)
- [x] T066 [US6] Add unit tests for MCP handler in pkg/cli/cmd/mcp/serve/handler_test.go
- [x] T067 [US6] Add integration test for MCP server in test/functional-portable/discovery/mcp_test.go

**Checkpoint**: MCP server starts and all skills are invocable via MCP protocol ✓

---

## Phase 7: User Story 3 - Match Recipes from Configured Sources (Priority: P2)

**Goal**: Find and suggest Recipes from AVM and internal repositories

**Independent Test**: Configure AVM and internal source, verify dependencies are matched with appropriate ranking

### Implementation for User Story 3

- [ ] T068 [US3] Create recipe source interface in pkg/cli/discovery/recipes/source.go
- [ ] T069 [US3] Implement AVM recipe source in pkg/cli/discovery/recipes/avm_source.go
- [ ] T070 [P] [US3] Implement OCI registry recipe source in pkg/cli/discovery/recipes/oci_source.go
- [ ] T071 [P] [US3] Implement Terraform registry recipe source in pkg/cli/discovery/recipes/terraform_source.go
- [ ] T072 [US3] Create recipe discoverer skill in pkg/cli/discovery/recipes/discoverer.go
- [ ] T073 [US3] Implement recipe ranking and prioritization in pkg/cli/discovery/recipes/ranker.go
- [ ] T074 [US3] Create `rad recipe source add` command in pkg/cli/cmd/recipe/source/add/add.go
- [ ] T075 [P] [US3] Create `rad recipe source list` command in pkg/cli/cmd/recipe/source/list/list.go
- [ ] T076 [P] [US3] Create `rad recipe source remove` command in pkg/cli/cmd/recipe/source/remove/remove.go
- [ ] T077 [US3] Wire recipe source commands into root.go in cmd/rad/cmd/root.go
- [ ] T078 [US3] Implement interactive recipe selection in rad app generate (FR-15)
- [ ] T079 [US3] Implement --recipe-profile flag for environment-specific recipes (FR-33)
- [ ] T080 [US3] Implement graceful degradation when sources unavailable (FR-22)
- [ ] T081 [US3] Implement auth to private sources (FR-37)
- [ ] T082 [US3] Add unit tests for recipe sources in pkg/cli/discovery/recipes/source_test.go
- [ ] T083 [P] [US3] Add unit tests for discoverer in pkg/cli/discovery/recipes/discoverer_test.go

**Checkpoint**: `rad app generate` suggests recipes from configured sources with ranking

---

## Phase 8: User Story 7 - Apply Team Infrastructure Practices (Priority: P2)

**Goal**: Detect and apply team practices (naming, tags, sizing) to generated definitions

**Independent Test**: Configure team practices and verify generated Resource Types incorporate them

### Implementation for User Story 7

- [ ] T084 [US7] Create team practices detector skill in pkg/cli/discovery/practices/detector.go
- [ ] T085 [US7] Implement IaC parser for Terraform in pkg/cli/discovery/practices/iac_parser.go
- [ ] T086 [P] [US7] Implement IaC parser for Bicep/ARM in pkg/cli/discovery/practices/bicep_parser.go
- [ ] T087 [US7] Create practices detection prompts in pkg/cli/discovery/llm/prompts/practices.txt
- [ ] T088 [US7] Implement documentation parser using LLM in pkg/cli/discovery/practices/doc_parser.go
- [ ] T089 [US7] Create .radius/team-practices.yaml schema and loader in pkg/cli/discovery/practices/config.go
- [ ] T090 [US7] Apply practices to Resource Type generation in pkg/cli/discovery/resourcetypes/generator.go
- [ ] T091 [US7] Apply naming conventions to generated Bicep in pkg/cli/discovery/generator/naming.go
- [ ] T092 [US7] Implement environment-specific practices (FR-08) in pkg/cli/discovery/practices/environments.go
- [ ] T093 [US7] Register discover_team_practices skill
- [ ] T094 [US7] Add unit tests for practices detector in pkg/cli/discovery/practices/detector_test.go
- [ ] T095 [P] [US7] Add unit tests for IaC parsers in pkg/cli/discovery/practices/iac_parser_test.go

**Checkpoint**: Generated definitions include team naming conventions, tags, and sizing

---

## Phase 9: User Story 5 - New Application Scaffolding (Priority: P3)

**Goal**: Scaffold new applications with infrastructure patterns from scratch

**Independent Test**: Use scaffolding command to create new app with PostgreSQL and Redis

### Implementation for User Story 5

- [ ] T096 [US5] Create `rad app scaffold` command in pkg/cli/cmd/app/scaffold/scaffold.go
- [ ] T097 [US5] Implement --name flag for application name
- [ ] T098 [US5] Implement --dependencies flag for specifying infrastructure types
- [ ] T099 [US5] Implement --container flag for specifying container services
- [ ] T100 [US5] Wire scaffold command into root.go in cmd/rad/cmd/root.go
- [ ] T101 [US5] Reuse Resource Type and Recipe discovery for scaffolding
- [ ] T102 [US5] Add unit tests for scaffold command in pkg/cli/cmd/app/scaffold/scaffold_test.go

**Checkpoint**: `rad app scaffold` creates starter app.bicep without existing code

---

## Phase 10: Polish & Cross-Cutting Concerns

**Purpose**: Improvements that affect multiple user stories

- [ ] T103 [P] Add --add-dependency flag to rad app generate (FR-24)
- [ ] T104 [P] Add --update flag for diff/patch mode (FR-20)
- [ ] T105 Implement structured JSON logging (NFR-05)
- [ ] T106 [P] Ensure deterministic output (NFR-04) - sort all collections, stable timestamps
- [ ] T107 Add performance optimization for large codebases (file batching, LLM caching)
- [ ] T108 [P] Add comprehensive error messages with actionable guidance (NFR-03)
- [ ] T109 Add documentation in docs/ for new commands
- [ ] T110 Run quickstart.md validation end-to-end
- [ ] T111 Security review: validate no secrets in output, no code execution

---

## Dependencies & Execution Order

### Phase Dependencies

```
Phase 1 (Setup)
    ↓
Phase 2 (Foundational) ─── BLOCKS ALL USER STORIES
    ↓
┌───────────────────────────────────────────────────────┐
│  P1 User Stories (can be parallel):                   │
│  ├── Phase 3: US1 (Dependency Detection) ─────────┐   │
│  ├── Phase 4: US2 (Resource Types) ──────────────├───┤
│  ├── Phase 5: US4 (App Generation) depends on ───┘   │
│  └── Phase 6: US6 (MCP Server) depends on all skills │
└───────────────────────────────────────────────────────┘
    ↓
┌───────────────────────────────────────────────────────┐
│  P2 User Stories (can be parallel):                   │
│  ├── Phase 7: US3 (Recipe Matching)                   │
│  └── Phase 8: US7 (Team Practices)                    │
└───────────────────────────────────────────────────────┘
    ↓
Phase 9: US5 (Scaffolding) - P3
    ↓
Phase 10: Polish
```

### User Story Dependencies

| Story | Depends On | Can Parallelize With |
|-------|------------|---------------------|
| US1 | Foundational | US2 (partial) |
| US2 | Foundational | US1 |
| US4 | US1, US2 | - |
| US6 | US1, US2, US4 | - |
| US3 | US4 | US7 |
| US7 | US1, US2 | US3 |
| US5 | US2, US3 | - |

### Parallel Opportunities per Phase

**Phase 1 (Setup)**:
- T002, T003, T004 can run in parallel
- T006, T008 can run in parallel

**Phase 2 (Foundational)**:
- T010 can run parallel with T009
- T013, T015 can run in parallel

**Phase 3 (US1)**:
- T018, T021 (prompts) can run parallel
- T030, T031, T032 (test fixtures) can run parallel
- T028, T029 (tests) can run parallel

**Phase 5 (US4)**:
- T042, T043, T044 (templates) can run parallel
- T054, T055 (tests) can run parallel

---

## Parallel Execution Examples

### Example: Phase 1 Setup
```bash
# Parallel group 1:
Task T002: Add dependencies to go.mod
Task T003: Create Skill interface
Task T004: Create skill registry

# Parallel group 2:
Task T006: Mock LLM provider
Task T008: Error types
```

### Example: Phase 3 (User Story 1)
```bash
# Parallel group - Prompts:
Task T018: Dependency detection prompts
Task T021: Service detection prompts

# Parallel group - Test fixtures:
Task T030: Node.js Express test fixture
Task T031: Python Flask test fixture
Task T032: Go Gin test fixture
```

---

## Implementation Strategy

### MVP First (P1 User Stories Only)

1. **Complete Setup + Foundational** → Core infrastructure ready
2. **Complete US1** (Analyze Codebase) → `rad app discover` works
3. **Complete US2** (Resource Types) → Dependencies map to types
4. **Complete US4** (Generate App) → `rad app generate` produces app.bicep
5. **Complete US6** (MCP Server) → AI agents can invoke skills
6. **STOP and VALIDATE**: Test full discover → generate → deploy workflow

### Incremental Delivery

| Milestone | Stories | Capabilities |
|-----------|---------|--------------|
| MVP | US1, US2, US4 | Discover → Generate workflow via CLI |
| + MCP | + US6 | AI agent integration |
| + Recipes | + US3 | External recipe matching |
| + Practices | + US7 | Team convention detection |
| + Scaffolding | + US5 | New app creation |

---

## Summary

| Metric | Value |
|--------|-------|
| **Total Tasks** | 111 |
| **Phase 1 (Setup)** | 8 tasks |
| **Phase 2 (Foundational)** | 8 tasks |
| **US1 (P1 - Discover)** | 17 tasks |
| **US2 (P1 - Resource Types)** | 8 tasks |
| **US4 (P1 - Generate)** | 15 tasks |
| **US6 (P1 - MCP)** | 11 tasks |
| **US3 (P2 - Recipes)** | 16 tasks |
| **US7 (P2 - Practices)** | 12 tasks |
| **US5 (P3 - Scaffold)** | 7 tasks |
| **Phase 10 (Polish)** | 9 tasks |
| **Parallel Opportunities** | 35+ tasks marked [P] |

---

## Notes

- All tasks marked [P] can run in parallel with other [P] tasks in same phase
- [Story] labels map tasks to user stories for traceability
- Each user story checkpoint enables independent testing
- MVP = Phases 1-6 (Setup through US6)
- Commit after each task or logical group
- Run `go test ./pkg/cli/discovery/...` after each implementation task
