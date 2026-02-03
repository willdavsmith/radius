# Implementation Plan: Automatic Application Discovery

**Branch**: `001-auto-app-discovery` | **Date**: February 3, 2026 | **Spec**: [spec.md](./spec.md)  
**Input**: Feature specification from `/specs/001-auto-app-discovery/spec.md`

---

## Summary

Radius will implement automatic application discovery to enable developers to go from existing codebase to deployable Radius application with zero manual Resource Type or Recipe authoring. The feature uses a skills-first architecture where all capabilities (dependency detection, service discovery, recipe matching, Bicep generation) are implemented as composable skills exposed via both CLI commands and MCP server for AI agent integration.

**Key Technical Approach**:
- Skills implemented as Go interfaces with JSON input/output schemas
- LLM-powered codebase analysis for language-agnostic dependency detection
- MCP server using `mcp-go` library for AI agent integration
- Template-based Bicep generation with validation
- Pluggable recipe sources (AVM, OCI registries, Terraform registries)

---

## Technical Context

**Language/Version**: Go 1.25 (per go.mod)  
**Primary Dependencies**: Cobra (CLI), mcp-go (MCP server), OpenAI Go SDK (LLM)  
**Storage**: Filesystem (`./radius/` output), `~/.rad/config.yaml` (configuration)  
**Testing**: go test with testify, table-driven tests with fixtures  
**Target Platform**: Linux, macOS, Windows (CLI); Kubernetes (deployment target)  
**Project Type**: Extension to existing radius monorepo  
**Performance Goals**: Discovery of ≤100 source files in <30 seconds (NFR-01)  
**Constraints**: Deterministic output (NFR-04), local-only analysis (no code transmission)  
**Scale/Scope**: Support Python, JavaScript/TypeScript, Go, Java, C# (FR-02)

---

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

| Principle | Status | Notes |
|-----------|--------|-------|
| I. API-First Design | ✅ PASS | Skills have defined JSON schemas; OpenAPI spec in contracts/ |
| II. Idiomatic Code Standards | ✅ PASS | Go code follows Effective Go; CLI uses Cobra patterns |
| III. Multi-Cloud Neutrality | ✅ PASS | Recipe system supports Azure, AWS; abstracted provider interfaces |
| IV. Testing Pyramid Discipline | ✅ PASS | Unit tests for skills, integration tests with fixtures, functional tests |
| V. Collaboration-Centric Design | ✅ PASS | Serves both developers (discovery) and platform engineers (practices) |
| VI. Open Source and Community-First | ✅ PASS | Spec in design-notes; skills are composable and extensible |
| VII. Simplicity Over Cleverness | ✅ PASS | LLM-based analysis avoids language-specific analyzers; template-based generation |
| VIII. Separation of Concerns | ✅ PASS | Skills layer, core engine, interface layer clearly separated |
| IX. Incremental Adoption | ✅ PASS | Multi-step workflow allows review; `--update` preserves manual changes |
| XVI. Repository-Specific Standards | ✅ PASS | Follows existing radius CLI patterns |

### Constitution Violations Requiring Justification

| Violation | Justification |
|-----------|--------------|
| None | No violations detected |

---

## Project Structure

### Documentation (this feature)

```text
specs/001-auto-app-discovery/
├── spec.md              # Feature specification
├── plan.md              # This file
├── research.md          # Phase 0 output - technical research
├── data-model.md        # Phase 1 output - entity definitions
├── quickstart.md        # Phase 1 output - usage guide
├── contracts/           # Phase 1 output - API contracts
│   └── skills-api.json  # OpenAPI spec for skills
├── checklists/          # Existing checklists
│   └── requirements.md
└── tasks.md             # Phase 2 output (created by /speckit.tasks)
```

### Source Code (repository root)

```text
# New packages for discovery feature
pkg/cli/discovery/
├── skills/              # Skill interface and registry
│   ├── skill.go         # Skill interface definition
│   ├── registry.go      # Skill registry for CLI and MCP
│   └── registry_test.go
├── dependencies/        # Dependency detection skill
│   ├── detector.go      # Main detection logic
│   ├── detector_test.go
│   └── testdata/        # Test fixtures
├── services/            # Service detection skill
│   ├── detector.go
│   ├── detector_test.go
│   └── testdata/
├── practices/           # Team practices detection skill
│   ├── detector.go
│   ├── iac_parser.go    # Parse Terraform/Bicep for practices
│   ├── doc_parser.go    # Parse documentation via LLM
│   └── testdata/
├── recipes/             # Recipe discovery skill
│   ├── discoverer.go
│   ├── source.go        # Recipe source interface
│   ├── avm_source.go    # Azure Verified Modules source
│   ├── oci_source.go    # OCI registry source
│   ├── terraform_source.go
│   └── testdata/
├── resourcetypes/       # Resource type generation skill
│   ├── generator.go
│   ├── catalog.go       # Pre-defined resource type catalog
│   └── testdata/
├── generator/           # Application definition generation skill
│   ├── bicep_generator.go
│   ├── templates/       # Bicep templates
│   │   ├── app.bicep.tmpl
│   │   ├── container.bicep.tmpl
│   │   └── resource.bicep.tmpl
│   ├── validator.go
│   └── testdata/
└── llm/                 # LLM provider abstraction
    ├── provider.go      # Provider interface
    ├── openai.go        # OpenAI implementation
    ├── ollama.go        # Ollama implementation
    ├── mock.go          # Mock for testing
    └── prompts/         # Prompt templates
        ├── dependencies.txt
        ├── services.txt
        └── practices.txt

# New CLI commands
pkg/cli/cmd/app/
├── discover/            # rad app discover command
│   ├── discover.go
│   └── discover_test.go
├── generate/            # rad app generate command
│   ├── generate.go
│   └── generate_test.go
└── scaffold/            # rad app scaffold command (P3)
    ├── scaffold.go
    └── scaffold_test.go

pkg/cli/cmd/mcp/
├── serve/               # rad mcp serve command
│   ├── serve.go
│   ├── handler.go       # MCP request handler
│   └── serve_test.go
└── mcp.go               # rad mcp parent command

pkg/cli/cmd/recipe/
└── source/              # rad recipe source commands
    ├── add/
    │   └── add.go
    ├── list/
    │   └── list.go
    └── remove/
        └── remove.go

# Configuration extensions
pkg/cli/config/
└── discovery_config.go  # Discovery-specific config handling

# Tests
test/functional-portable/
└── discovery/           # Functional tests
    ├── discover_test.go
    ├── generate_test.go
    ├── mcp_test.go
    └── testdata/        # Sample codebases
        ├── nodejs-express/
        ├── python-flask/
        └── go-gin/
```

**Structure Decision**: Extends existing radius monorepo structure following established patterns in `pkg/cli/cmd/` for CLI commands and `pkg/cli/` for shared packages. New `discovery/` package contains all feature-specific logic organized by skill.

---

## Architecture Overview

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                           USER INTERFACES                                   │
├─────────────────────────────────────────────────────────────────────────────┤
│  ┌──────────────────┐  ┌──────────────────┐  ┌──────────────────┐          │
│  │    rad CLI       │  │   MCP Server     │  │   Go API         │          │
│  │  (Cobra cmds)    │  │  (mcp-go)        │  │  (pkg import)    │          │
│  └────────┬─────────┘  └────────┬─────────┘  └────────┬─────────┘          │
│           │                     │                     │                     │
│           └─────────────────────┼─────────────────────┘                     │
│                                 │                                           │
│                                 ▼                                           │
├─────────────────────────────────────────────────────────────────────────────┤
│                           SKILLS LAYER                                      │
├─────────────────────────────────────────────────────────────────────────────┤
│  ┌────────────────┐  ┌────────────────┐  ┌────────────────┐                │
│  │ discover_      │  │ discover_      │  │ discover_      │                │
│  │ dependencies   │  │ services       │  │ team_practices │                │
│  └───────┬────────┘  └───────┬────────┘  └───────┬────────┘                │
│          │                   │                   │                          │
│  ┌────────────────┐  ┌────────────────┐  ┌────────────────┐                │
│  │ generate_      │  │ discover_      │  │ generate_      │                │
│  │ resource_types │  │ recipes        │  │ app_definition │                │
│  └───────┬────────┘  └───────┬────────┘  └───────┬────────┘                │
│          │                   │                   │                          │
│          └───────────────────┼───────────────────┘                          │
│                              │                                              │
│                              ▼                                              │
├─────────────────────────────────────────────────────────────────────────────┤
│                           CORE ENGINE                                       │
├─────────────────────────────────────────────────────────────────────────────┤
│  ┌────────────────┐  ┌────────────────┐  ┌────────────────┐                │
│  │  LLM Provider  │  │ IaC Parsers    │  │ Bicep Generator│                │
│  │  (OpenAI,      │  │ (HCL, Bicep,   │  │ (Templates)    │                │
│  │   Ollama)      │  │  ARM)          │  │                │                │
│  └───────┬────────┘  └───────┬────────┘  └───────┬────────┘                │
│          │                   │                   │                          │
│  ┌────────────────┐  ┌────────────────┐  ┌────────────────┐                │
│  │ Recipe Sources │  │ Resource Type  │  │ Bicep          │                │
│  │ (AVM, OCI,     │  │ Catalog        │  │ Validator      │                │
│  │  Terraform)    │  │                │  │                │                │
│  └────────────────┘  └────────────────┘  └────────────────┘                │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## Implementation Phases

### Phase 0: Research (Completed)

See [research.md](./research.md) for:
- Skills architecture decisions
- LLM integration approach
- MCP server implementation
- Recipe source integration
- Team practices detection
- Discovery output format
- Bicep generation strategy

### Phase 1: Design (Completed)

See:
- [data-model.md](./data-model.md) - Entity definitions and relationships
- [contracts/skills-api.json](./contracts/skills-api.json) - OpenAPI specification
- [quickstart.md](./quickstart.md) - Usage guide and examples

### Phase 2: Tasks (Next Step)

Tasks will be generated by `/speckit.tasks` command to create actionable implementation work items organized by user story priority.

---

## Dependencies

### External Dependencies (New)

| Package | Version | Purpose |
|---------|---------|---------|
| `github.com/mark3labs/mcp-go` | latest | MCP server implementation |
| `github.com/sashabaranov/go-openai` | latest | OpenAI API client |
| `github.com/hashicorp/hcl/v2` | latest | Terraform HCL parsing |

### Internal Dependencies (Existing)

| Package | Purpose |
|---------|---------|
| `pkg/cli/framework` | CLI command framework |
| `pkg/cli/bicep` | Bicep compilation and validation |
| `pkg/cli/config` | Configuration management |
| `pkg/cli/connections` | Radius API connectivity |
| `pkg/cli/workspaces` | Workspace context |
| `pkg/cli/output` | CLI output formatting |

---

## Risk Assessment

| Risk | Impact | Mitigation |
|------|--------|------------|
| LLM response variability | Medium | Use structured output mode; implement deterministic fallback |
| AVM API changes | Low | Abstract behind interface; version lock discovery |
| Performance for large codebases | Medium | Implement file filtering; batch LLM calls; add caching |
| MCP protocol evolution | Low | Use stable mcp-go library; abstract transport layer |

---

## Success Metrics (from Spec)

| ID | Metric | Target |
|----|--------|--------|
| SC-002 | Dependency detection accuracy | ≥85% for sample frameworks |
| SC-003 | First-deploy success rate | 90% when recipes available |
| SC-004 | Zero manual authoring | Complete workflow without Resource Type/Recipe authoring |
| SC-005 | Recipe source onboarding | <10 minutes |

---

## Next Steps

1. **Run `/speckit.tasks`** to generate implementation tasks from this plan
2. **Create GitHub issues** for each task
3. **Begin implementation** starting with P1 user stories:
   - US-1: Analyze Existing Application Codebase
   - US-2: Map to Resource Types
   - US-4: Generate Application Definition
   - US-6: AI Coding Agent Integration

---

## References

- [Feature Specification](./spec.md)
- [Research Findings](./research.md)
- [Data Model](./data-model.md)
- [API Contracts](./contracts/skills-api.json)
- [Quickstart Guide](./quickstart.md)
- [Constitution](.specify/memory/constitution.md)
