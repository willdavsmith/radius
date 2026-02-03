---
analyzed_at: "2026-02-03T13:46:18-08:00"
codebase_path: /Users/willsmith/dev/radius/radius/test/functional-portable/discovery/testdata/nodejs-express
dependencies:
  - confidence: 0.7
    evidence:
      - package.json:postgres
      - package.json:pg
    technology: PostgreSQL
    type: postgresql
  - confidence: 0.7
    evidence:
      - package.json:redis
    technology: Redis
    type: redis
services:
  - entrypoint: src/index.js
    name: src
    port: 3000
    type: http
version: "1.0"
---

# Discovery Report

Generated: 2026-02-03 13:46:18
Codebase: /Users/willsmith/dev/radius/radius/test/functional-portable/discovery/testdata/nodejs-express

## Dependencies Detected

| Type | Technology | Confidence | Evidence |
|------|------------|------------|----------|
| postgresql | PostgreSQL | 70% | package.json, package.json |
| redis | Redis | 70% | package.json |



## Services Detected

| Name | Type | Port | Entrypoint |
|------|------|------|------------|
| src | http | 3000 | src/index.js |





## Next Steps

1. Review the detected dependencies and services above
2. Run `rad app generate` to create app.bicep
3. Customize the generated application definition as needed
4. Run `rad deploy` to deploy your application
