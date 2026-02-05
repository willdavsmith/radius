/*
Copyright 2023 The Radius Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package skills

import (
	"context"
	"encoding/json"
	"strings"

	"github.com/radius-project/radius/pkg/cli/discovery"
)

// SuggestRecipes suggests appropriate Radius recipes for detected dependencies.
type SuggestRecipes struct{}

// NewSuggestRecipes creates a new suggest recipes skill.
func NewSuggestRecipes() *SuggestRecipes {
	return &SuggestRecipes{}
}

// Name returns the skill identifier.
func (s *SuggestRecipes) Name() string {
	return "suggest_recipes"
}

// Description returns a human-readable description.
func (s *SuggestRecipes) Description() string {
	return `Suggests appropriate Radius recipes for detected infrastructure dependencies.
Provides guidance on recipe selection based on environment type (dev, staging, production)
and deployment target (local, Azure, AWS).`
}

// InputSchema returns the JSON Schema for the skill input.
func (s *SuggestRecipes) InputSchema() json.RawMessage {
	return json.RawMessage(`{
		"type": "object",
		"properties": {
			"dependencies": {
				"type": "array",
				"description": "List of detected dependencies",
				"items": {
					"type": "object",
					"properties": {
						"type": {"type": "string"},
						"technology": {"type": "string"}
					}
				}
			},
			"environment": {
				"type": "string",
				"description": "Target environment (dev, staging, production)",
				"default": "dev"
			},
			"platform": {
				"type": "string",
				"description": "Target platform (local, azure, aws)",
				"default": "local"
			}
		},
		"required": ["dependencies"]
	}`)
}

// OutputSchema returns the JSON Schema for the skill output.
func (s *SuggestRecipes) OutputSchema() json.RawMessage {
	return json.RawMessage(`{
		"type": "object",
		"properties": {
			"suggestions": {
				"type": "array",
				"items": {
					"type": "object",
					"properties": {
						"dependencyType": {"type": "string"},
						"recipeName": {"type": "string"},
						"recipeSource": {"type": "string"},
						"description": {"type": "string"},
						"considerations": {"type": "array", "items": {"type": "string"}}
					}
				}
			}
		}
	}`)
}

// SuggestRecipesInput defines the input parameters.
type SuggestRecipesInput struct {
	Dependencies []discovery.DetectedDependency `json:"dependencies"`
	Environment  string                         `json:"environment,omitempty"`
	Platform     string                         `json:"platform,omitempty"`
}

// RecipeSuggestion represents a recipe suggestion.
type RecipeSuggestion struct {
	DependencyType string   `json:"dependencyType"`
	RecipeName     string   `json:"recipeName"`
	RecipeSource   string   `json:"recipeSource"`
	Description    string   `json:"description"`
	Considerations []string `json:"considerations"`
}

// SuggestRecipesOutput defines the output structure.
type SuggestRecipesOutput struct {
	Suggestions []RecipeSuggestion `json:"suggestions"`
}

// Execute runs the recipe suggestion logic.
func (s *SuggestRecipes) Execute(ctx context.Context, input json.RawMessage) (json.RawMessage, error) {
	var in SuggestRecipesInput
	if err := json.Unmarshal(input, &in); err != nil {
		return nil, discovery.WrapError(discovery.ErrCodeInputValidation, "invalid input", err)
	}

	// Set defaults
	if in.Environment == "" {
		in.Environment = "dev"
	}
	if in.Platform == "" {
		in.Platform = "local"
	}

	output := SuggestRecipesOutput{
		Suggestions: make([]RecipeSuggestion, 0, len(in.Dependencies)),
	}

	for _, dep := range in.Dependencies {
		suggestion := s.suggestRecipeForDependency(dep, in.Environment, in.Platform)
		output.Suggestions = append(output.Suggestions, suggestion)
	}

	return json.Marshal(output)
}

// suggestRecipeForDependency generates a recipe suggestion for a single dependency.
func (s *SuggestRecipes) suggestRecipeForDependency(dep discovery.DetectedDependency, env, platform string) RecipeSuggestion {
	suggestion := RecipeSuggestion{
		DependencyType: dep.Type,
		Considerations: make([]string, 0),
	}

	// Recipe catalog based on dependency type, environment, and platform
	switch strings.ToLower(dep.Type) {
	case "postgresql":
		suggestion = s.suggestPostgresRecipe(env, platform)
	case "mysql":
		suggestion = s.suggestMySQLRecipe(env, platform)
	case "mongodb":
		suggestion = s.suggestMongoDBRecipe(env, platform)
	case "redis":
		suggestion = s.suggestRedisRecipe(env, platform)
	case "rabbitmq":
		suggestion = s.suggestRabbitMQRecipe(env, platform)
	case "kafka":
		suggestion = s.suggestKafkaRecipe(env, platform)
	default:
		suggestion.RecipeName = "default"
		suggestion.RecipeSource = "ghcr.io/radius-project/recipes/local-dev/extenders:latest"
		suggestion.Description = "Generic extender recipe for custom resource types"
		suggestion.Considerations = append(suggestion.Considerations,
			"Consider creating a custom recipe for this dependency type",
			"Manual configuration may be required")
	}

	suggestion.DependencyType = dep.Type
	return suggestion
}

func (s *SuggestRecipes) suggestPostgresRecipe(env, platform string) RecipeSuggestion {
	suggestion := RecipeSuggestion{
		Considerations: make([]string, 0),
	}

	switch platform {
	case "azure":
		if env == "production" {
			suggestion.RecipeName = "azure-postgres-flexible"
			suggestion.RecipeSource = "ghcr.io/radius-project/recipes/azure/sqldatabases:latest"
			suggestion.Description = "Azure Database for PostgreSQL Flexible Server - production-ready managed PostgreSQL"
			suggestion.Considerations = append(suggestion.Considerations,
				"Enables high availability and automatic backups",
				"Consider configuring VNet integration for security",
				"Review pricing tier based on expected workload")
		} else {
			suggestion.RecipeName = "azure-postgres-dev"
			suggestion.RecipeSource = "ghcr.io/radius-project/recipes/azure/sqldatabases:latest"
			suggestion.Description = "Azure Database for PostgreSQL - cost-optimized for development"
			suggestion.Considerations = append(suggestion.Considerations,
				"Uses burstable compute tier for cost savings",
				"Not recommended for production workloads")
		}
	case "aws":
		suggestion.RecipeName = "aws-rds-postgres"
		suggestion.RecipeSource = "ghcr.io/radius-project/recipes/aws/sqldatabases:latest"
		suggestion.Description = "Amazon RDS for PostgreSQL"
		suggestion.Considerations = append(suggestion.Considerations,
			"Configure appropriate instance class for workload",
			"Enable Multi-AZ for production high availability")
	default:
		suggestion.RecipeName = "default"
		suggestion.RecipeSource = "ghcr.io/radius-project/recipes/local-dev/sqldatabases:latest"
		suggestion.Description = "Local PostgreSQL container for development"
		suggestion.Considerations = append(suggestion.Considerations,
			"Data is ephemeral - will be lost on container restart",
			"Suitable for local development only")
	}

	return suggestion
}

func (s *SuggestRecipes) suggestMySQLRecipe(env, platform string) RecipeSuggestion {
	suggestion := RecipeSuggestion{
		Considerations: make([]string, 0),
	}

	switch platform {
	case "azure":
		suggestion.RecipeName = "azure-mysql-flexible"
		suggestion.RecipeSource = "ghcr.io/radius-project/recipes/azure/sqldatabases:latest"
		suggestion.Description = "Azure Database for MySQL Flexible Server"
	case "aws":
		suggestion.RecipeName = "aws-rds-mysql"
		suggestion.RecipeSource = "ghcr.io/radius-project/recipes/aws/sqldatabases:latest"
		suggestion.Description = "Amazon RDS for MySQL"
	default:
		suggestion.RecipeName = "default"
		suggestion.RecipeSource = "ghcr.io/radius-project/recipes/local-dev/sqldatabases:latest"
		suggestion.Description = "Local MySQL container for development"
	}

	return suggestion
}

func (s *SuggestRecipes) suggestMongoDBRecipe(env, platform string) RecipeSuggestion {
	suggestion := RecipeSuggestion{
		Considerations: make([]string, 0),
	}

	switch platform {
	case "azure":
		suggestion.RecipeName = "azure-cosmosdb-mongo"
		suggestion.RecipeSource = "ghcr.io/radius-project/recipes/azure/mongodatabases:latest"
		suggestion.Description = "Azure Cosmos DB with MongoDB API"
		suggestion.Considerations = append(suggestion.Considerations,
			"Globally distributed with configurable consistency levels",
			"Consider RU (Request Unit) provisioning for cost control")
	case "aws":
		suggestion.RecipeName = "aws-documentdb"
		suggestion.RecipeSource = "ghcr.io/radius-project/recipes/aws/mongodatabases:latest"
		suggestion.Description = "Amazon DocumentDB (MongoDB compatible)"
	default:
		suggestion.RecipeName = "default"
		suggestion.RecipeSource = "ghcr.io/radius-project/recipes/local-dev/mongodatabases:latest"
		suggestion.Description = "Local MongoDB container for development"
	}

	return suggestion
}

func (s *SuggestRecipes) suggestRedisRecipe(env, platform string) RecipeSuggestion {
	suggestion := RecipeSuggestion{
		Considerations: make([]string, 0),
	}

	switch platform {
	case "azure":
		if env == "production" {
			suggestion.RecipeName = "azure-redis-premium"
			suggestion.RecipeSource = "ghcr.io/radius-project/recipes/azure/rediscaches:latest"
			suggestion.Description = "Azure Cache for Redis Premium tier"
			suggestion.Considerations = append(suggestion.Considerations,
				"Enables clustering and data persistence",
				"Supports VNet integration")
		} else {
			suggestion.RecipeName = "azure-redis-basic"
			suggestion.RecipeSource = "ghcr.io/radius-project/recipes/azure/rediscaches:latest"
			suggestion.Description = "Azure Cache for Redis Basic tier for development"
		}
	case "aws":
		suggestion.RecipeName = "aws-elasticache-redis"
		suggestion.RecipeSource = "ghcr.io/radius-project/recipes/aws/rediscaches:latest"
		suggestion.Description = "Amazon ElastiCache for Redis"
	default:
		suggestion.RecipeName = "default"
		suggestion.RecipeSource = "ghcr.io/radius-project/recipes/local-dev/rediscaches:latest"
		suggestion.Description = "Local Redis container for development"
	}

	return suggestion
}

func (s *SuggestRecipes) suggestRabbitMQRecipe(env, platform string) RecipeSuggestion {
	suggestion := RecipeSuggestion{
		Considerations: make([]string, 0),
	}

	switch platform {
	case "azure":
		suggestion.RecipeName = "azure-servicebus"
		suggestion.RecipeSource = "ghcr.io/radius-project/recipes/azure/rabbitmqqueues:latest"
		suggestion.Description = "Azure Service Bus (RabbitMQ compatible via AMQP)"
		suggestion.Considerations = append(suggestion.Considerations,
			"Native Azure messaging service",
			"May require minor code changes for full compatibility")
	case "aws":
		suggestion.RecipeName = "aws-mq-rabbitmq"
		suggestion.RecipeSource = "ghcr.io/radius-project/recipes/aws/rabbitmqqueues:latest"
		suggestion.Description = "Amazon MQ for RabbitMQ"
	default:
		suggestion.RecipeName = "default"
		suggestion.RecipeSource = "ghcr.io/radius-project/recipes/local-dev/rabbitmqqueues:latest"
		suggestion.Description = "Local RabbitMQ container for development"
	}

	return suggestion
}

func (s *SuggestRecipes) suggestKafkaRecipe(env, platform string) RecipeSuggestion {
	suggestion := RecipeSuggestion{
		Considerations: make([]string, 0),
	}

	switch platform {
	case "azure":
		suggestion.RecipeName = "azure-eventhubs-kafka"
		suggestion.RecipeSource = "ghcr.io/radius-project/recipes/azure/pubsubbrokers:latest"
		suggestion.Description = "Azure Event Hubs with Kafka protocol support"
		suggestion.Considerations = append(suggestion.Considerations,
			"Managed Kafka-compatible streaming platform",
			"May require client configuration changes")
	case "aws":
		suggestion.RecipeName = "aws-msk"
		suggestion.RecipeSource = "ghcr.io/radius-project/recipes/aws/pubsubbrokers:latest"
		suggestion.Description = "Amazon Managed Streaming for Apache Kafka (MSK)"
	default:
		suggestion.RecipeName = "default"
		suggestion.RecipeSource = "ghcr.io/radius-project/recipes/local-dev/pubsubbrokers:latest"
		suggestion.Description = "Local Kafka container for development"
		suggestion.Considerations = append(suggestion.Considerations,
			"Includes Zookeeper for coordination",
			"Resource-intensive for local development")
	}

	return suggestion
}
