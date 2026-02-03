# Discovery Test Fixtures

This directory contains test fixtures for the automatic application discovery feature.

## Structure

- `nodejs-express/` - Sample Node.js Express application with PostgreSQL and Redis
- `python-flask/` - Sample Python Flask application with MongoDB
- `go-gin/` - Sample Go Gin application with PostgreSQL

## Usage

These fixtures are used by integration tests to verify that the discovery feature
correctly identifies dependencies and services in various codebases.

Each fixture contains:
- Application source code
- Dependency manifests (package.json, requirements.txt, go.mod)
- Expected discovery results for validation
