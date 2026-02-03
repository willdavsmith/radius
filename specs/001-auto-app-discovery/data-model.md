# Data Model: Automatic Application Discovery

**Feature Branch**: `001-auto-app-discovery`  
**Date**: February 3, 2026  
**Spec**: [spec.md](./spec.md)

---

## Entity Relationship Diagram

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                          DISCOVERY DOMAIN                                   │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│  ┌──────────────┐         ┌──────────────────┐       ┌──────────────────┐  │
│  │    Skill     │         │ DiscoveryResult  │       │  DetectedDep     │  │
│  ├──────────────┤         ├──────────────────┤       ├──────────────────┤  │
│  │ name         │         │ analyzedAt       │  1:N  │ type             │  │
│  │ description  │         │ codebasePath     │◀─────▶│ technology       │  │
│  │ inputSchema  │         │ dependencies[]   │       │ confidence       │  │
│  │ outputSchema │         │ services[]       │       │ evidence[]       │  │
│  └──────────────┘         │ practices        │       │ connectionInfo   │  │
│                           └──────────────────┘       └──────────────────┘  │
│                                    │                                        │
│                                    │ 1:N                                    │
│                                    ▼                                        │
│                           ┌──────────────────┐       ┌──────────────────┐  │
│                           │ DetectedService  │       │  TeamPractices   │  │
│                           ├──────────────────┤       ├──────────────────┤  │
│                           │ name             │       │ naming           │  │
│                           │ type             │       │ tags[]           │  │
│                           │ entrypoint       │       │ sizing{}         │  │
│                           │ port             │       │ sources[]        │  │
│                           │ dockerfile       │       │ environments{}   │  │
│                           └──────────────────┘       └──────────────────┘  │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────────────────┐
│                          GENERATION DOMAIN                                  │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│  ┌──────────────────┐     ┌──────────────────┐     ┌──────────────────┐    │
│  │  ResourceType    │     │     Recipe       │     │  RecipeSource    │    │
│  ├──────────────────┤     ├──────────────────┤     ├──────────────────┤    │
│  │ name             │1:N  │ name             │ N:1 │ type             │    │
│  │ type             │◀───▶│ iacType          │◀───▶│ location         │    │
│  │ schema           │     │ templatePath     │     │ priority         │    │
│  │ outputs[]        │     │ parameters[]     │     │ authConfig       │    │
│  │ maturity         │     │ source           │     │ refreshFreq      │    │
│  └──────────────────┘     │ environments[]   │     └──────────────────┘    │
│           │               └──────────────────┘                              │
│           │                        │                                        │
│           │ N:1                    │ N:M                                    │
│           ▼                        ▼                                        │
│  ┌──────────────────┐     ┌──────────────────┐                             │
│  │ ResourceTypeCat  │     │ RecipeSelection  │                             │
│  ├──────────────────┤     ├──────────────────┤                             │
│  │ catalog{}        │     │ dependency       │                             │
│  │ version          │     │ recipe           │                             │
│  │ lastUpdated      │     │ parameters{}     │                             │
│  └──────────────────┘     │ environment      │                             │
│                           └──────────────────┘                             │
│                                    │                                        │
│                                    │ N:1                                    │
│                                    ▼                                        │
│                           ┌──────────────────┐                             │
│                           │ AppDefinition    │                             │
│                           ├──────────────────┤                             │
│                           │ name             │                             │
│                           │ containers[]     │                             │
│                           │ resources[]      │                             │
│                           │ connections[]    │                             │
│                           │ bicepContent     │                             │
│                           └──────────────────┘                             │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## Entity Definitions

### Skill

Composable capability exposed via MCP and CLI.

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `name` | string | ✓ | Unique skill identifier (e.g., `discover_dependencies`) |
| `description` | string | ✓ | Human-readable description for MCP tool listing |
| `inputSchema` | JSONSchema | ✓ | JSON Schema defining valid input structure |
| `outputSchema` | JSONSchema | ✓ | JSON Schema defining output structure |

**Validation Rules**:
- `name` must match pattern `^[a-z][a-z0-9_]*$`
- Schemas must be valid JSON Schema Draft-07

---

### DiscoveryResult

Output of the discovery phase containing all detected information.

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `version` | string | ✓ | Schema version (e.g., "1.0") |
| `analyzedAt` | timestamp | ✓ | When analysis was performed |
| `codebasePath` | string | ✓ | Absolute path to analyzed codebase |
| `dependencies` | DetectedDependency[] | ✓ | List of detected infrastructure dependencies |
| `services` | DetectedService[] | ✓ | List of detected deployable services |
| `practices` | TeamPractices | ✓ | Detected or configured team practices |
| `warnings` | Warning[] | | Partial failures or low-confidence detections |

**State Transitions**:
- `pending` → `analyzing` → `complete` | `failed`

---

### DetectedDependency

An infrastructure dependency detected from codebase analysis.

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `type` | string | ✓ | Normalized dependency type (e.g., "postgresql", "redis") |
| `technology` | string | ✓ | Specific technology/service (e.g., "PostgreSQL", "Azure Redis") |
| `confidence` | float | ✓ | Detection confidence score (0.0 - 1.0) |
| `evidence` | Evidence[] | ✓ | Sources supporting this detection |
| `connectionInfo` | ConnectionInfo | | Detected connection configuration |

**Confidence Tiers** (per FR-05):
- High: ≥ 0.80
- Medium: 0.50 - 0.79
- Low: < 0.50 (excluded by default)

**Validation Rules**:
- `confidence` must be between 0.0 and 1.0
- At least one evidence item required

---

### Evidence

Supporting evidence for a detected dependency.

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `source` | string | ✓ | File path relative to codebase root |
| `type` | string | ✓ | Evidence type: "package", "import", "config", "connection_string" |
| `value` | string | ✓ | The detected value (e.g., "pg@8.11.0") |
| `line` | int | | Line number in source file |

---

### ConnectionInfo

Detected connection configuration for a dependency.

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `envVar` | string | | Environment variable name (e.g., "DATABASE_URL") |
| `defaultPort` | int | | Default port for the service |
| `hostPattern` | string | | Detected hostname pattern |

---

### DetectedService

A deployable service/entrypoint detected in the codebase.

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `name` | string | ✓ | Service name (derived from directory or config) |
| `type` | string | ✓ | Framework type (e.g., "express", "flask", "spring") |
| `entrypoint` | string | | Main file or command |
| `port` | int | | Detected or configured port |
| `dockerfile` | string | | Path to Dockerfile if detected |
| `buildContext` | string | | Build context directory |

**Validation Rules**:
- `name` must be valid Kubernetes name (lowercase, alphanumeric, hyphens)
- `port` must be 1-65535 if specified

---

### TeamPractices

Infrastructure conventions detected or configured for the team.

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `naming` | NamingConvention | | Naming pattern for resources |
| `tags` | TagPolicy | | Required and optional tags |
| `sizing` | map[string]SizingPolicy | | Environment-specific sizing (key: environment name) |
| `sources` | PracticeSource[] | | Where practices were detected |

---

### NamingConvention

Naming pattern for infrastructure resources.

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `pattern` | string | ✓ | Pattern with placeholders (e.g., "{env}-{service}-{resource}") |
| `maxLength` | int | | Maximum resource name length |
| `environmentAbbreviations` | map[string]string | | Short forms for environments |

---

### TagPolicy

Required and recommended tags for resources.

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `required` | string[] | ✓ | Tags that must be present |
| `recommended` | string[] | | Tags that should be present |
| `defaults` | map[string]string | | Default values for tags |

---

### SizingPolicy

Environment-specific sizing and configuration.

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `defaultTier` | string | ✓ | Default pricing tier |
| `haEnabled` | bool | ✓ | High availability enabled |
| `backupEnabled` | bool | | Backup enabled |
| `geoRedundant` | bool | | Geo-redundancy enabled |

---

### PracticeSource

Location where a practice was detected.

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `type` | string | ✓ | Source type: "iac", "config", "wiki", "manual" |
| `location` | string | ✓ | File path or URL |
| `detectedAt` | timestamp | ✓ | When this source was analyzed |

---

### ResourceType

Interface defining what an application consumes from infrastructure.

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `name` | string | ✓ | Human-friendly name (e.g., "PostgreSQL Database") |
| `type` | string | ✓ | Radius type (e.g., "Applications.Datastores/postgreSql") |
| `apiVersion` | string | ✓ | API version |
| `schema` | ResourceSchema | ✓ | Input/output schema |
| `outputs` | OutputDefinition[] | ✓ | Outputs provided to containers |
| `maturity` | string | ✓ | "alpha", "beta", "stable" |

---

### ResourceSchema

Schema defining properties for a Resource Type.

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `properties` | map[string]PropertyDef | ✓ | Property definitions |
| `required` | string[] | | Required property names |

---

### PropertyDef

Definition of a single property.

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `type` | string | ✓ | JSON type |
| `description` | string | ✓ | Property description |
| `default` | any | | Default value |
| `enum` | any[] | | Allowed values |
| `sensitive` | bool | | Whether value is sensitive |

---

### OutputDefinition

Output provided by a Resource Type to containers.

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `name` | string | ✓ | Output name (e.g., "connectionString") |
| `type` | string | ✓ | Output type |
| `description` | string | ✓ | Description |
| `sensitive` | bool | | Whether output is sensitive |

---

### Recipe

IaC implementation that provisions infrastructure.

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `name` | string | ✓ | Recipe name |
| `description` | string | | Human-readable description |
| `iacType` | string | ✓ | "bicep" or "terraform" |
| `templatePath` | string | ✓ | Registry path or URL |
| `parameters` | map[string]ParameterDef | | Recipe parameters |
| `source` | RecipeSource | ✓ | Where recipe was found |
| `environments` | string[] | | Suitable environments (e.g., ["dev", "production"]) |
| `resourceType` | string | ✓ | Target Resource Type |
| `verified` | bool | | Whether from verified source (e.g., AVM) |

---

### RecipeSource

Location for discovering Recipes.

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `type` | string | ✓ | "avm", "oci", "terraform-registry", "filesystem" |
| `location` | string | ✓ | Registry URL or path |
| `priority` | int | ✓ | Search priority (lower = higher priority) |
| `authConfig` | AuthConfig | | Authentication configuration |
| `refreshFrequency` | duration | | How often to refresh cache |

---

### AuthConfig

Authentication configuration for recipe sources.

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `type` | string | ✓ | "none", "token", "oidc", "env" |
| `tokenEnvVar` | string | | Environment variable containing token |
| `credentialHelper` | string | | Credential helper command |

---

### RecipeSelection

User's selection of a recipe for a dependency.

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `dependency` | string | ✓ | Dependency type |
| `recipe` | Recipe | ✓ | Selected recipe |
| `parameters` | map[string]any | | Parameter overrides |
| `environment` | string | | Target environment |

---

### AppDefinition

Generated Radius application definition.

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `name` | string | ✓ | Application name |
| `containers` | ContainerDef[] | ✓ | Container definitions |
| `resources` | ResourceDef[] | ✓ | Resource definitions |
| `connections` | ConnectionDef[] | ✓ | Container-to-resource connections |
| `bicepContent` | string | ✓ | Generated Bicep content |
| `valid` | bool | ✓ | Whether Bicep validated successfully |
| `errors` | ValidationError[] | | Validation errors if any |

---

### ContainerDef

Container definition in the application.

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `name` | string | ✓ | Container name |
| `image` | string | | Container image (or TODO placeholder) |
| `ports` | PortDef[] | | Exposed ports |
| `env` | map[string]EnvValue | | Environment variables |

---

### ResourceDef

Resource definition in the application.

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `name` | string | ✓ | Resource name |
| `type` | string | ✓ | Resource Type |
| `recipe` | string | | Recipe reference |
| `properties` | map[string]any | | Resource properties |

---

### ConnectionDef

Connection between container and resource.

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `container` | string | ✓ | Container name |
| `resource` | string | ✓ | Resource name |
| `envMapping` | map[string]string | | Environment variable mappings |

---

## Configuration Entities

### RadiusConfig

Global Radius configuration (extends existing `~/.rad/config.yaml`).

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `discovery` | DiscoveryConfig | | Discovery-specific configuration |
| `recipeSources` | RecipeSource[] | | Configured recipe sources |
| `llmProvider` | LLMConfig | | LLM provider configuration |

---

### DiscoveryConfig

Discovery-specific configuration.

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `outputPath` | string | | Custom output path (default: "./radius") |
| `confidenceThreshold` | float | | Minimum confidence (default: 0.5) |
| `excludePaths` | string[] | | Paths to exclude from analysis |

---

### LLMConfig

LLM provider configuration.

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `provider` | string | ✓ | "openai", "ollama", "anthropic" |
| `model` | string | | Model name (default: provider-specific) |
| `endpoint` | string | | Custom endpoint URL |
| `apiKeyEnvVar` | string | | Environment variable for API key |

---

## Validation Rules Summary

| Entity | Rule | Error Code |
|--------|------|------------|
| Skill.name | Must match `^[a-z][a-z0-9_]*$` | INVALID_SKILL_NAME |
| DetectedDependency.confidence | Must be 0.0 - 1.0 | INVALID_CONFIDENCE |
| DetectedService.name | Must be valid K8s name | INVALID_SERVICE_NAME |
| DetectedService.port | Must be 1-65535 | INVALID_PORT |
| ResourceType.type | Must match Radius type pattern | INVALID_RESOURCE_TYPE |
| Recipe.iacType | Must be "bicep" or "terraform" | INVALID_IAC_TYPE |
| RecipeSource.priority | Must be positive integer | INVALID_PRIORITY |

---

## State Machines

### Discovery State Machine

```
┌─────────┐     start()      ┌───────────┐
│ pending │──────────────────▶│ analyzing │
└─────────┘                   └─────┬─────┘
                                    │
                    ┌───────────────┴───────────────┐
                    │                               │
                success()                       error()
                    │                               │
                    ▼                               ▼
            ┌──────────┐                      ┌────────┐
            │ complete │                      │ failed │
            └──────────┘                      └────────┘
```

### Generation State Machine

```
┌─────────┐   selectRecipes()   ┌────────────┐   generate()   ┌────────────┐
│ pending │─────────────────────▶│ selecting  │───────────────▶│ generating │
└─────────┘                      └────────────┘                └─────┬──────┘
                                                                     │
                                               ┌─────────────────────┴──────────────────────┐
                                               │                                            │
                                           validate()                                   error()
                                               │                                            │
                                               ▼                                            ▼
                                        ┌───────────┐                                 ┌────────┐
                                        │ validated │                                 │ failed │
                                        └─────┬─────┘                                 └────────┘
                                              │
                              ┌───────────────┴───────────────┐
                              │                               │
                          valid()                         invalid()
                              │                               │
                              ▼                               ▼
                       ┌──────────┐                    ┌────────────┐
                       │ complete │                    │ validation │
                       └──────────┘                    │   errors   │
                                                       └────────────┘
```
