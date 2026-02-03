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

package resourcetypes

import "github.com/radius-project/radius/pkg/cli/discovery"

// builtinCatalogEntries contains the built-in resource type catalog entries.
var builtinCatalogEntries = []*CatalogEntry{
	{
		DependencyType: "postgresql",
		Aliases:        []string{"postgres", "pg", "pgsql"},
		DefaultRecipe:  "default",
		ResourceType: discovery.ResourceType{
			Name:       "PostgreSQL Database",
			Type:       "Applications.Datastores/sqlDatabases",
			APIVersion: "2023-10-01-preview",
			Maturity:   "stable",
			Schema: &discovery.ResourceSchema{
				Properties: map[string]discovery.PropertyDef{
					"database": {
						Type:        "string",
						Description: "The name of the SQL database",
					},
					"server": {
						Type:        "string",
						Description: "The fully qualified domain name of the SQL database server",
					},
					"port": {
						Type:        "integer",
						Description: "The port of the SQL database server",
						Default:     5432,
					},
					"username": {
						Type:        "string",
						Description: "The username for the SQL database",
						Sensitive:   true,
					},
					"password": {
						Type:        "string",
						Description: "The password for the SQL database",
						Sensitive:   true,
					},
				},
				Required: []string{"database"},
			},
			Outputs: []discovery.OutputDefinition{
				{Name: "connectionString", Type: "string", Description: "The connection string for the database", Sensitive: true},
				{Name: "database", Type: "string", Description: "The database name"},
				{Name: "host", Type: "string", Description: "The database host"},
				{Name: "port", Type: "integer", Description: "The database port"},
				{Name: "username", Type: "string", Description: "The database username", Sensitive: true},
				{Name: "password", Type: "string", Description: "The database password", Sensitive: true},
			},
		},
	},
	{
		DependencyType: "redis",
		Aliases:        []string{"redis-cache", "cache"},
		DefaultRecipe:  "default",
		ResourceType: discovery.ResourceType{
			Name:       "Redis Cache",
			Type:       "Applications.Datastores/redisCaches",
			APIVersion: "2023-10-01-preview",
			Maturity:   "stable",
			Schema: &discovery.ResourceSchema{
				Properties: map[string]discovery.PropertyDef{
					"host": {
						Type:        "string",
						Description: "The Redis cache host",
					},
					"port": {
						Type:        "integer",
						Description: "The Redis cache port",
						Default:     6379,
					},
					"tls": {
						Type:        "boolean",
						Description: "Whether TLS is enabled",
						Default:     false,
					},
				},
			},
			Outputs: []discovery.OutputDefinition{
				{Name: "connectionString", Type: "string", Description: "The connection string for the cache", Sensitive: true},
				{Name: "host", Type: "string", Description: "The cache host"},
				{Name: "port", Type: "integer", Description: "The cache port"},
				{Name: "password", Type: "string", Description: "The cache password", Sensitive: true},
			},
		},
	},
	{
		DependencyType: "mongodb",
		Aliases:        []string{"mongo", "documentdb"},
		DefaultRecipe:  "default",
		ResourceType: discovery.ResourceType{
			Name:       "MongoDB Database",
			Type:       "Applications.Datastores/mongoDatabases",
			APIVersion: "2023-10-01-preview",
			Maturity:   "stable",
			Schema: &discovery.ResourceSchema{
				Properties: map[string]discovery.PropertyDef{
					"database": {
						Type:        "string",
						Description: "The name of the MongoDB database",
					},
					"host": {
						Type:        "string",
						Description: "The MongoDB host",
					},
					"port": {
						Type:        "integer",
						Description: "The MongoDB port",
						Default:     27017,
					},
				},
				Required: []string{"database"},
			},
			Outputs: []discovery.OutputDefinition{
				{Name: "connectionString", Type: "string", Description: "The connection string for the database", Sensitive: true},
				{Name: "database", Type: "string", Description: "The database name"},
				{Name: "host", Type: "string", Description: "The database host"},
				{Name: "port", Type: "integer", Description: "The database port"},
			},
		},
	},
	{
		DependencyType: "rabbitmq",
		Aliases:        []string{"amqp", "rabbit"},
		DefaultRecipe:  "default",
		ResourceType: discovery.ResourceType{
			Name:       "RabbitMQ Message Queue",
			Type:       "Applications.Messaging/rabbitMQQueues",
			APIVersion: "2023-10-01-preview",
			Maturity:   "stable",
			Schema: &discovery.ResourceSchema{
				Properties: map[string]discovery.PropertyDef{
					"queue": {
						Type:        "string",
						Description: "The name of the queue",
					},
					"host": {
						Type:        "string",
						Description: "The RabbitMQ host",
					},
					"port": {
						Type:        "integer",
						Description: "The RabbitMQ port",
						Default:     5672,
					},
					"vhost": {
						Type:        "string",
						Description: "The RabbitMQ virtual host",
						Default:     "/",
					},
				},
				Required: []string{"queue"},
			},
			Outputs: []discovery.OutputDefinition{
				{Name: "connectionString", Type: "string", Description: "The connection string for the queue", Sensitive: true},
				{Name: "host", Type: "string", Description: "The queue host"},
				{Name: "port", Type: "integer", Description: "The queue port"},
			},
		},
	},
	{
		DependencyType: "kafka",
		Aliases:        []string{"event-stream"},
		DefaultRecipe:  "default",
		ResourceType: discovery.ResourceType{
			Name:       "Kafka Event Stream",
			Type:       "Applications.Messaging/kafkaTopics",
			APIVersion: "2023-10-01-preview",
			Maturity:   "alpha",
			Schema: &discovery.ResourceSchema{
				Properties: map[string]discovery.PropertyDef{
					"topic": {
						Type:        "string",
						Description: "The name of the Kafka topic",
					},
					"brokers": {
						Type:        "string",
						Description: "Comma-separated list of Kafka brokers",
					},
				},
				Required: []string{"topic"},
			},
			Outputs: []discovery.OutputDefinition{
				{Name: "brokers", Type: "string", Description: "The Kafka brokers"},
				{Name: "topic", Type: "string", Description: "The Kafka topic"},
			},
		},
	},
	{
		DependencyType: "mysql",
		Aliases:        []string{"mariadb"},
		DefaultRecipe:  "default",
		ResourceType: discovery.ResourceType{
			Name:       "MySQL Database",
			Type:       "Applications.Datastores/sqlDatabases",
			APIVersion: "2023-10-01-preview",
			Maturity:   "stable",
			Schema: &discovery.ResourceSchema{
				Properties: map[string]discovery.PropertyDef{
					"database": {
						Type:        "string",
						Description: "The name of the SQL database",
					},
					"server": {
						Type:        "string",
						Description: "The fully qualified domain name of the SQL database server",
					},
					"port": {
						Type:        "integer",
						Description: "The port of the SQL database server",
						Default:     3306,
					},
				},
				Required: []string{"database"},
			},
			Outputs: []discovery.OutputDefinition{
				{Name: "connectionString", Type: "string", Description: "The connection string for the database", Sensitive: true},
				{Name: "database", Type: "string", Description: "The database name"},
				{Name: "host", Type: "string", Description: "The database host"},
				{Name: "port", Type: "integer", Description: "The database port"},
				{Name: "username", Type: "string", Description: "The database username", Sensitive: true},
				{Name: "password", Type: "string", Description: "The database password", Sensitive: true},
			},
		},
	},
	{
		DependencyType: "blob-storage",
		Aliases:        []string{"s3", "azure-blob", "object-storage", "minio"},
		DefaultRecipe:  "default",
		ResourceType: discovery.ResourceType{
			Name:       "Blob Storage",
			Type:       "Applications.Datastores/blobContainers",
			APIVersion: "2023-10-01-preview",
			Maturity:   "alpha",
			Schema: &discovery.ResourceSchema{
				Properties: map[string]discovery.PropertyDef{
					"container": {
						Type:        "string",
						Description: "The name of the blob container",
					},
				},
				Required: []string{"container"},
			},
			Outputs: []discovery.OutputDefinition{
				{Name: "connectionString", Type: "string", Description: "The connection string for blob storage", Sensitive: true},
				{Name: "container", Type: "string", Description: "The container name"},
			},
		},
	},
	{
		DependencyType: "dapr-statestore",
		Aliases:        []string{"dapr-state"},
		DefaultRecipe:  "default",
		ResourceType: discovery.ResourceType{
			Name:       "Dapr State Store",
			Type:       "Applications.Dapr/stateStores",
			APIVersion: "2023-10-01-preview",
			Maturity:   "stable",
			Schema: &discovery.ResourceSchema{
				Properties: map[string]discovery.PropertyDef{
					"type": {
						Type:        "string",
						Description: "The type of state store component",
					},
				},
			},
			Outputs: []discovery.OutputDefinition{
				{Name: "componentName", Type: "string", Description: "The Dapr component name"},
			},
		},
	},
	{
		DependencyType: "dapr-pubsub",
		Aliases:        []string{"dapr-pubsub-broker"},
		DefaultRecipe:  "default",
		ResourceType: discovery.ResourceType{
			Name:       "Dapr Pub/Sub Broker",
			Type:       "Applications.Dapr/pubSubBrokers",
			APIVersion: "2023-10-01-preview",
			Maturity:   "stable",
			Schema: &discovery.ResourceSchema{
				Properties: map[string]discovery.PropertyDef{
					"type": {
						Type:        "string",
						Description: "The type of pub/sub component",
					},
				},
			},
			Outputs: []discovery.OutputDefinition{
				{Name: "componentName", Type: "string", Description: "The Dapr component name"},
			},
		},
	},
	{
		DependencyType: "dapr-secretstore",
		Aliases:        []string{"dapr-secrets"},
		DefaultRecipe:  "default",
		ResourceType: discovery.ResourceType{
			Name:       "Dapr Secret Store",
			Type:       "Applications.Dapr/secretStores",
			APIVersion: "2023-10-01-preview",
			Maturity:   "stable",
			Schema: &discovery.ResourceSchema{
				Properties: map[string]discovery.PropertyDef{
					"type": {
						Type:        "string",
						Description: "The type of secret store component",
					},
				},
			},
			Outputs: []discovery.OutputDefinition{
				{Name: "componentName", Type: "string", Description: "The Dapr component name"},
			},
		},
	},
}
