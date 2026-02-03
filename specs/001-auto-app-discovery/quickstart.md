# Quickstart: Automatic Application Discovery

This guide walks through the end-to-end workflow for using Radius Automatic Application Discovery to go from an existing codebase to a deployed application.

## Prerequisites

- Radius CLI installed and configured (`rad init` completed)
- A codebase with infrastructure dependencies (databases, caches, etc.)
- LLM provider configured (see [LLM Configuration](#llm-configuration))

## The Three-Phase Workflow

```
┌─────────────┐     ┌─────────────┐     ┌─────────────┐
│  DISCOVER   │────▶│   GENERATE  │────▶│   DEPLOY    │
│             │     │             │     │             │
│ rad app     │     │ rad app     │     │ rad deploy  │
│ discover .  │     │ generate    │     │ ./radius/   │
│             │     │             │     │ app.bicep   │
└─────────────┘     └─────────────┘     └─────────────┘
```

---

## Step 1: Discover (`rad app discover`)

Analyze your codebase to detect infrastructure dependencies, services, and team practices.

```bash
# Navigate to your project root
cd /path/to/my-app

# Run discovery
rad app discover .
```

### Example Output

```
🔍 Analyzing codebase...

✅ Discovery Complete

Detected Dependencies:
┌────────────────┬────────────┬────────────┐
│ Technology     │ Confidence │ Evidence   │
├────────────────┼────────────┼────────────┤
│ PostgreSQL     │ 95%        │ pg@8.11.0  │
│ Redis          │ 92%        │ ioredis@5.3│
│ Azure Blob     │ 88%        │ @azure/blob│
└────────────────┴────────────┴────────────┘

Services Detected:
• api-server (Express.js, port 3000)
• worker (background job processor)

Team Practices Detected:
• Naming: {env}-{service}-{resource}
• Required Tags: cost-center, owner, env
  (Source: /infra/main.tf)

📄 Output: ./radius/discovery.md

Run 'rad app generate' to create your application definition.
```

### Discovery Output

The discovery phase creates `./radius/discovery.md` containing:
- All detected dependencies with confidence scores
- Detected services and their configurations
- Team practices extracted from existing IaC files

You can review and edit this file before proceeding to generation.

---

## Step 2: Generate (`rad app generate`)

Create a Radius application definition from the discovery results.

```bash
# Generate application definition
rad app generate
```

### Interactive Recipe Selection

When multiple recipe options exist, you'll be prompted to choose:

```
📋 Generating Resource Types...
   ✓ Applications.Datastores/postgreSql
   ✓ Applications.Datastores/redisCaches
   ✓ Applications.Datastores/blobStorage

🔍 Discovering IaC implementations...

PostgreSQL:
  [1] Azure Database for PostgreSQL (AVM)
      ✓ Production  ✓ HA  ✓ Auto-backup
  [2] PostgreSQL Container (Dev)
      ⚡ Fast  ⚠ Not for production

Select Recipe for PostgreSQL [1]: 1

Redis:
  [1] Azure Cache for Redis (AVM)
  [2] Redis Container (Dev)

Select Recipe for Redis [1]: 1

📝 Generating application definition...
   ✓ Wiring container connections
   ✓ Mapping environment variables
   ✓ Validating Bicep syntax

✅ Created: ./radius/app.bicep
```

### Non-Interactive Mode

For CI/CD pipelines, use `--accept-defaults`:

```bash
rad app generate --accept-defaults
```

### Specifying Environment Profile

Use `--recipe-profile` to select environment-specific recipes:

```bash
# Use production-ready recipes
rad app generate --recipe-profile production

# Use development recipes (containers, fast startup)
rad app generate --recipe-profile dev
```

---

## Step 3: Deploy (`rad deploy`)

Deploy your generated application to Radius.

```bash
# Deploy to default environment
rad deploy ./radius/app.bicep

# Deploy to a specific environment
rad deploy ./radius/app.bicep -e production
```

### Example Output

```
🚀 Deploying my-ecommerce-app...

✓ Creating Resource Types (if needed)
✓ Provisioning product-db... [2m 34s]
✓ Provisioning session-cache... [1m 12s]
✓ Provisioning image-storage... [45s]
✓ Deploying api-server... [38s]
✓ Deploying worker... [32s]

✅ Deployment Complete!

Application Endpoints:
• api-server: https://my-ecommerce-app-api.azurecontainerapps.io

Run 'rad app status' for health info.
```

---

## AI Agent Integration (MCP)

Use the same discovery capabilities through AI coding agents like GitHub Copilot.

### Start the MCP Server

```bash
# Start MCP server for VS Code integration (stdio)
rad mcp serve

# Start MCP server for remote agents (HTTP)
rad mcp serve --transport http --port 8080
```

### VS Code Configuration

Add to `.vscode/mcp.json`:

```json
{
  "servers": {
    "radius": {
      "command": "rad",
      "args": ["mcp", "serve"]
    }
  }
}
```

### Conversational Workflow

With the MCP server running, you can use natural language:

> "Help me deploy my Node.js e-commerce app using Radius"

The AI agent will:
1. Invoke `discover_dependencies` and `discover_services` skills
2. Present findings and ask for confirmation
3. Invoke `discover_recipes` to find suitable IaC implementations
4. Invoke `generate_app_definition` to create `app.bicep`
5. Optionally deploy with your confirmation

---

## Configuration

### Recipe Sources

Configure recipe sources in `~/.rad/config.yaml`:

```yaml
recipeSources:
  - type: avm
    location: mcr.microsoft.com/bicep/avm
    priority: 1
  - type: oci
    location: myregistry.azurecr.io/recipes
    priority: 2
  - type: terraform-registry
    location: app.terraform.io/my-org
    priority: 3
```

### Add a New Recipe Source

```bash
rad recipe source add myregistry \
  --type oci \
  --location myregistry.azurecr.io/recipes \
  --priority 2
```

### LLM Configuration

Configure LLM provider for codebase analysis in `~/.rad/config.yaml`:

```yaml
llmProvider:
  provider: openai  # or ollama, anthropic
  model: gpt-4o     # optional, uses default
  apiKeyEnvVar: OPENAI_API_KEY
```

For air-gapped environments using Ollama:

```yaml
llmProvider:
  provider: ollama
  model: llama3.1
  endpoint: http://localhost:11434
```

### Team Practices

Create `.radius/team-practices.yaml` in your project:

```yaml
naming:
  pattern: "{env}-{service}-{resource}"
  environments:
    dev: d
    staging: s
    production: p

tags:
  required:
    - cost-center
    - owner
    - environment
  defaults:
    managed-by: radius

sizing:
  dev:
    defaultTier: Basic
    haEnabled: false
  production:
    defaultTier: Premium
    haEnabled: true
    geoRedundant: true
```

---

## Common Workflows

### Update an Existing Application

When your codebase changes (new dependencies added):

```bash
# Re-run discovery
rad app discover .

# Generate with update mode (preserves manual changes)
rad app generate --update

# Review diff and deploy
rad deploy ./radius/app.bicep
```

### Add a Manual Dependency

Add a dependency without detection:

```bash
rad app generate --add-dependency redis
```

### Scaffold a New Application

Start a new application from scratch:

```bash
rad app scaffold --name my-new-app \
  --dependencies postgresql,redis \
  --container api-server
```

---

## Troubleshooting

### Low Confidence Detections

If important dependencies are detected with low confidence:

1. Check `./radius/discovery.md` for the evidence
2. Add explicit configuration in your codebase (e.g., environment variables)
3. Manually add to discovery output if needed

### No Recipes Found

If no recipes match a dependency:

1. Check configured recipe sources: `rad recipe source list`
2. Add additional sources or create a custom recipe
3. Use `--add-dependency` to specify the type manually

### Bicep Validation Errors

If generated Bicep fails validation:

1. Check `./radius/app.bicep` for TODO placeholders
2. Ensure Dockerfiles exist or specify container images
3. Run `rad bicep build ./radius/app.bicep` for detailed errors

---

## Next Steps

- **Customize Generated Definition**: Edit `./radius/app.bicep` to add custom configuration
- **Add Custom Recipes**: Create and register your own infrastructure recipes
- **Configure Team Practices**: Set up organizational standards in `.radius/team-practices.yaml`
- **Explore MCP Skills**: Use individual skills for advanced automation
