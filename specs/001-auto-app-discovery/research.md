# Research: Automatic Application Discovery

**Feature Branch**: `001-auto-app-discovery`  
**Date**: February 3, 2026  
**Spec**: [spec.md](./spec.md)

## Research Summary

This document captures research findings for implementing the Automatic Application Discovery feature. All "NEEDS CLARIFICATION" items from the technical context have been resolved through codebase analysis and best practices research.

---

## 1. Skills Architecture for CLI and MCP

### Decision
Implement skills as Go interfaces with JSON input/output, usable by both CLI commands and MCP server.

### Rationale
- The Radius CLI uses Cobra with a consistent `framework.Runner` pattern (see `pkg/cli/framework/`)
- MCP servers in Go can be implemented using the `github.com/mark3labs/mcp-go` library (widely adopted for Go MCP servers)
- Skills should be pure functions with JSON schemas for input/output, enabling deterministic behavior regardless of invocation path

### Alternatives Considered
1. **Separate implementations for CLI and MCP** - Rejected because it violates DRY and the spec's "skills-first" architecture (DC-006)
2. **Shell script-based skills** - Rejected because Go provides better type safety and testing

### Implementation Pattern
```go
// pkg/cli/skills/skill.go
type Skill interface {
    Name() string
    Description() string
    InputSchema() json.RawMessage
    OutputSchema() json.RawMessage
    Execute(ctx context.Context, input json.RawMessage) (json.RawMessage, error)
}

// CLI commands invoke skills via the Skills Registry
// MCP server exposes skills as tools via the same registry
```

---

## 2. LLM Integration for Codebase Analysis

### Decision
Use a pluggable LLM provider interface with OpenAI as the default, supporting local models via Ollama for air-gapped environments.

### Rationale
- The spec explicitly calls for LLM-based analysis (Q-9: "Use LLMs to read and analyze the codebase")
- OpenAI's structured output mode ensures deterministic JSON responses
- Ollama support enables air-gapped/enterprise deployments

### Alternatives Considered
1. **Language-specific static analyzers** - Rejected per spec Q-9; LLM approach is simpler and supports all languages
2. **Single LLM provider** - Rejected because enterprises need flexibility for approved providers

### Implementation Pattern
```go
// pkg/cli/discovery/llm/provider.go
type LLMProvider interface {
    Analyze(ctx context.Context, prompt string, schema json.RawMessage) (json.RawMessage, error)
}

// Implementations: OpenAIProvider, OllamaProvider, AnthropicProvider
// Configuration via ~/.rad/config.yaml or environment variables
```

### Security Considerations
- Local-only analysis: file contents sent to LLM for analysis
- User must configure LLM provider; no default cloud provider without explicit opt-in
- Support for self-hosted models addresses data sovereignty concerns

---

## 3. MCP Server Implementation

### Decision
Implement MCP server using `github.com/mark3labs/mcp-go` with stdio transport (primary) and HTTP transport (secondary).

### Rationale
- FR-29 requires `rad mcp serve` command
- FR-30 requires stdio (VS Code) and HTTP (remote agents) transports
- The `mcp-go` library is the most mature Go MCP implementation

### Implementation Pattern
```go
// pkg/cli/cmd/mcp/serve.go
// Registers all skills as MCP tools
// Supports --transport=stdio|http, --port, --auth-mode flags
```

---

## 4. Recipe Source Integration

### Decision
Support three recipe source types: Azure Verified Modules (AVM), internal Bicep registries (OCI), and internal Terraform registries (Terraform Cloud/Enterprise or filesystem).

### Rationale
- FR-14 requires searching configured recipe sources
- AVM provides authoritative, well-tested Azure modules
- Enterprise customers have existing internal module registries

### Alternatives Considered
1. **AVM only** - Rejected because FR-14 explicitly requires internal sources
2. **Git repositories as sources** - Deferred to future iteration; focus on registry-based sources first

### Implementation Pattern
```go
// pkg/cli/discovery/recipes/source.go
type RecipeSource interface {
    Search(ctx context.Context, resourceType string) ([]RecipeMatch, error)
    Metadata() SourceMetadata
}

// Implementations: AVMSource, OCIRegistrySource, TerraformRegistrySource
```

---

## 5. Team Practices Detection

### Decision
Detect team practices from existing IaC files using pattern matching, with optional LLM enhancement for unstructured documentation.

### Rationale
- FR-06 requires detecting practices from Terraform/Bicep/ARM files
- FR-07 requires extracting practices from wikis/Confluence/Notion
- Pattern matching handles structured IaC; LLM handles unstructured docs

### Implementation Pattern
```go
// pkg/cli/discovery/practices/detector.go
type PracticesDetector interface {
    Detect(ctx context.Context, path string) (*TeamPractices, error)
}

// IaCDetector: parses HCL/Bicep/ARM for naming patterns, tags, sizing
// DocumentDetector: uses LLM to extract practices from markdown/wiki content
```

### Practices Schema
```yaml
# .radius/team-practices.yaml
naming:
  pattern: "{env}-{service}-{resource}"
  environments:
    dev: "d"
    staging: "s"
    production: "p"
tags:
  required:
    - cost-center
    - owner
    - environment
sizing:
  dev:
    default_tier: "Basic"
    ha_enabled: false
  production:
    default_tier: "Premium"
    ha_enabled: true
```

---

## 6. Discovery Output Format

### Decision
Output discovery results to `./radius/discovery.md` as human-readable Markdown with embedded YAML frontmatter for machine parsing.

### Rationale
- FR-09 requires output to `./radius/discovery.md`
- Markdown is reviewable by developers
- YAML frontmatter enables re-parsing for `rad app generate`

### Output Structure
```markdown
---
version: "1.0"
analyzed_at: "2026-02-03T10:00:00Z"
dependencies:
  - type: "postgresql"
    confidence: 0.95
    evidence: ["package.json:pg@8.11.0"]
services:
  - name: "api-server"
    type: "express"
    port: 3000
practices:
  naming:
    pattern: "{env}-{service}-{resource}"
---

# Discovery Report

## Dependencies Detected
...
```

---

## 7. Bicep Generation Strategy

### Decision
Generate Bicep application definitions using the existing `pkg/cli/bicep/` infrastructure, with template-based generation for containers and resources.

### Rationale
- DC-001 requires Bicep as the target format
- Radius already has Bicep generation/validation in `pkg/cli/bicep/`
- Template-based generation with Go text/template provides flexibility

### Implementation Pattern
```go
// pkg/cli/discovery/generator/bicep.go
type BicepGenerator struct {
    templates *template.Template
}

func (g *BicepGenerator) Generate(ctx context.Context, discovery *DiscoveryResult, recipes []RecipeSelection) (*GeneratedApp, error)
```

---

## 8. Resource Type Catalog vs. Dynamic Generation

### Decision
Maintain a pre-defined catalog of Resource Types in-memory, with fallback to minimal schema generation for unknown types.

### Rationale
- Per Q-1: "predefined catalog of Resource Types maintained in the Resource Types repository"
- Catalog provides quality/consistency for common dependencies
- Fallback generation ensures extensibility

### Catalog Structure
```go
// pkg/cli/discovery/resourcetypes/catalog.go
var Catalog = map[string]ResourceType{
    "postgresql": {
        Type: "Applications.Datastores/postgreSql",
        Schema: ResourceSchema{...},
    },
    "redis": {
        Type: "Applications.Datastores/redisCaches",
        Schema: ResourceSchema{...},
    },
    // ...
}
```

---

## 9. Existing Radius Integration Points

### Decision
Integrate with existing Radius infrastructure:
- Use `pkg/cli/connections/` for Radius API connectivity
- Use `pkg/cli/config/` for configuration management
- Use `pkg/cli/workspaces/` for workspace context
- Use `pkg/cli/deploy/` for deployment operations

### Rationale
- Radius has mature infrastructure for these concerns
- Reuse avoids duplication and ensures consistency

---

## 10. Testing Strategy

### Decision
- Unit tests for individual skills using mock LLM providers
- Integration tests using pre-analyzed codebases with expected outputs
- Functional tests deploying generated apps to test environments

### Rationale
- NFR-04 requires deterministic output (same input → identical output)
- Testing against known codebases ensures accuracy metrics (SC-002: 85% detection rate)

---

## Technology Decisions Summary

| Area | Decision | Primary Library/Tool |
|------|----------|---------------------|
| CLI Framework | Cobra with `framework.Runner` | `github.com/spf13/cobra` |
| MCP Server | stdio + HTTP transports | `github.com/mark3labs/mcp-go` |
| LLM Integration | Pluggable provider interface | OpenAI Go SDK, Ollama API |
| Recipe Sources | Registry-based discovery | OCI client, Terraform registry API |
| Bicep Generation | Template-based | Go `text/template` |
| Team Practices | Pattern matching + LLM | HCL parser, Bicep AST |
| Testing | Table-driven with fixtures | `testing` + `testify` |

---

## Open Items for Implementation

1. **LLM Prompt Engineering**: Develop and refine prompts for codebase analysis and documentation extraction
2. **AVM Integration**: Research Azure Verified Modules API/registry structure for recipe discovery
3. **Performance Optimization**: Implement caching for LLM responses during multi-file analysis
4. **Error Taxonomy**: Define structured error types for partial failures (FR-10)
