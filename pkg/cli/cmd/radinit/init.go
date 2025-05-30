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

package radinit

import (
	"context"
	"fmt"
	"os"

	v1 "github.com/radius-project/radius/pkg/armrpc/api/v1"
	"github.com/radius-project/radius/pkg/cli"

	"github.com/radius-project/radius/pkg/cli/aws"
	"github.com/radius-project/radius/pkg/cli/azure"
	"github.com/radius-project/radius/pkg/cli/clierrors"
	"github.com/radius-project/radius/pkg/cli/cmd/commonflags"
	"github.com/radius-project/radius/pkg/cli/connections"
	cli_credential "github.com/radius-project/radius/pkg/cli/credential"
	"github.com/radius-project/radius/pkg/cli/framework"
	"github.com/radius-project/radius/pkg/cli/helm"
	"github.com/radius-project/radius/pkg/cli/kubernetes"
	"github.com/radius-project/radius/pkg/cli/output"
	"github.com/radius-project/radius/pkg/cli/prompt"
	"github.com/radius-project/radius/pkg/cli/setup"
	"github.com/radius-project/radius/pkg/cli/workspaces"
	corerp "github.com/radius-project/radius/pkg/corerp/api/v20231001preview"
	"github.com/radius-project/radius/pkg/to"
	ucp "github.com/radius-project/radius/pkg/ucp/api/v20231001preview"
	"github.com/spf13/cobra"
)

// NOTE: this command is very super big so it's broken up amongst a few files.

// NewCommand creates an instance of the command and runner for the `rad init` command.
//

// This function "NewCommand" creates a new Cobra command with flags and a runner, which can be used to initialize the
// Radius control-plane.
func NewCommand(factory framework.Factory) (*cobra.Command, framework.Runner) {
	runner := NewRunner(factory)

	cmd := &cobra.Command{
		Use:     "initialize",
		Aliases: []string{"init"},
		Short:   "Initialize Radius",
		Long: `
Interactively install the Radius control-plane and setup an environment.

If an environment already exists, 'rad init' will prompt the user to use the existing environment or create a new one.

By default, 'rad init' will optimize for a developer-focused environment with an environment named "default" and Recipes that support prototyping, development and testing using lightweight containers. These environments are great for building and testing your application.

Specifying the '--full' flag will cause 'rad init' to prompt the user for all available configuration options such as Kubernetes context, environment name, and cloud providers. This is useful for fully customizing your environment.
`,
		Example: `
## Create a new development environment named "default"
rad init

## Prompt the user for all available options to create a new environment
rad init --full
`,
		Args: cobra.ExactArgs(0),
		RunE: framework.RunCommand(runner),
	}

	// Define your flags here
	commonflags.AddOutputFlag(cmd)
	cmd.Flags().Bool("full", false, "Prompt user for all available configuration options")
	return cmd, runner
}

// Runner is the runner implementation for the `rad init` command.
type Runner struct {
	azureClient azure.Client
	awsClient   aws.Client

	// ConfigFileInterface is the interface for the config file.
	ConfigFileInterface framework.ConfigFileInterface

	// ConfigHolder is the interface for the config holder.
	ConfigHolder *framework.ConfigHolder

	// ConnectionFactory is the interface for the connection factory.
	ConnectionFactory connections.Factory

	// HelmInterface is the interface for the helm client.
	HelmInterface helm.Interface

	// KubernetesInterface is the interface for the kubernetes client.
	KubernetesInterface kubernetes.Interface

	// Output is the interface for console output.
	Output output.Interface

	// Prompter is the interface for the prompter.
	Prompter prompt.Interface

	// DevRecipeClient is the interface for the dev recipe client.
	DevRecipeClient DevRecipeClient

	// Format is the output format.
	Format string

	// Workspace is the workspace to use. This will be populated by Validate.
	Workspace *workspaces.Workspace

	// Full determines whether or not we ask the user for all options.
	Full bool

	// Options provides the options to used for Radius initialization. This will be populated by Validate.
	Options *initOptions
}

// NewRunner creates a new instance of the `rad init` runner.
//

// NewRunner creates a new Runner struct with the given factory's ConfigHolder, Output, ConnectionFactory, Prompter,
// ConfigFileInterface, KubernetesInterface, HelmInterface, DevRecipeClient, AWSClient, and AzureClient.
func NewRunner(factory framework.Factory) *Runner {
	return &Runner{
		ConfigHolder:        factory.GetConfigHolder(),
		Output:              factory.GetOutput(),
		ConnectionFactory:   factory.GetConnectionFactory(),
		Prompter:            factory.GetPrompter(),
		ConfigFileInterface: factory.GetConfigFileInterface(),
		KubernetesInterface: factory.GetKubernetesInterface(),
		HelmInterface:       factory.GetHelmInterface(),
		DevRecipeClient:     NewDevRecipeClient(),
		awsClient:           factory.GetAWSClient(),
		azureClient:         factory.GetAzureClient(),
	}
}

// Validate runs validation for the `rad init` command.
// Validates the user prompts, values provided and builds the picture for the backend to execute
//

// Validate gathers input from the user, creates a workspace and options, and confirms the options with the user before
// returning the options and workspace. If the user does not confirm the options, the function will loop and gather input again.
// If an error occurs, the function will return an error.
func (r *Runner) Validate(cmd *cobra.Command, args []string) error {
	format, err := cli.RequireOutput(cmd)
	if err != nil {
		return err
	}
	r.Format = format

	r.Full, err = cmd.Flags().GetBool("full")
	if err != nil {
		return err
	}

	for {
		options, workspace, err := r.enterInitOptions(cmd.Context())
		if err != nil {
			return err
		}

		// Show a confirmation screen if we're in full mode.
		confirmed := true
		if r.Full {
			confirmed, err = r.confirmOptions(cmd.Context(), options)
			if err != nil {
				return err
			}
		}

		if confirmed {
			r.Options = options
			r.Workspace = workspace
			return nil
		}

		// User did not confirm the summary, so gather input again.
	}
}

// Run runs the `rad init` command.
// Creates Radius resources, azure resources if required based on the user input, command flags
//

// Run creates a progress channel, installs the radius control plane, creates an environment, configures cloud
// providers, scaffolds an application, and updates the config file, all while displaying progress updates to the UI.
func (r *Runner) Run(ctx context.Context) error {
	config := r.ConfigFileInterface.ConfigFromContext(ctx)

	// Use this channel to send progress updates to the UI.
	progressChan := make(chan progressMsg)
	progressCompleteChan := make(chan error)
	progress := progressMsg{}

	go func() {
		// Show dynamic UI.
		err := r.showProgress(ctx, r.Options, progressChan)
		if err != nil {
			progressCompleteChan <- err
		}
		close(progressCompleteChan)
	}()

	if r.Options.Cluster.Install {
		cliOptions := helm.CLIClusterOptions{
			Radius: helm.ChartOptions{
				SetArgs: r.Options.SetValues,
			},
		}

		clusterOptions := helm.PopulateDefaultClusterOptions(cliOptions)

		err := r.HelmInterface.InstallRadius(ctx, clusterOptions, r.Options.Cluster.Context)
		if err != nil {
			return clierrors.MessageWithCause(err, "Failed to install Radius.")
		}
	}
	progress.InstallComplete = true
	progressChan <- progress

	if r.Options.Environment.Create {
		err := r.CreateEnvironment(ctx)
		if err != nil {
			return err
		}
	}
	progress.EnvironmentComplete = true
	progressChan <- progress

	if r.Options.Application.Scaffold {
		client, err := r.ConnectionFactory.CreateApplicationsManagementClient(ctx, *r.Workspace)
		if err != nil {
			return err
		}

		// Initialize the application resource if it's not found. This supports the scenario where the application
		// resource is not defined in bicep.
		err = client.CreateApplicationIfNotFound(ctx, r.Options.Application.Name, &corerp.ApplicationResource{
			Location: to.Ptr(v1.LocationGlobal),
			Properties: &corerp.ApplicationProperties{
				Environment: &r.Workspace.Environment,
			},
		})
		if err != nil {
			return err
		}

		// Scaffold application files in the current directory
		wd, err := os.Getwd()
		if err != nil {
			return err
		}

		err = setup.ScaffoldApplication(wd, r.Options.Application.Name)
		if err != nil {
			return err
		}
	}
	progress.ApplicationComplete = true
	progressChan <- progress

	err := r.ConfigFileInterface.EditWorkspaces(ctx, config, r.Workspace)
	if err != nil {
		return err
	}
	progress.ConfigComplete = true
	progressChan <- progress

	// Wait for UI to complete.
	err = <-progressCompleteChan
	if err != nil {
		return err
	}

	return nil
}

func (r *Runner) getAzureCredential() (ucp.AzureCredentialResource, error) {
	switch r.Options.CloudProviders.Azure.CredentialKind {
	case azure.AzureCredentialKindServicePrincipal:
		return ucp.AzureCredentialResource{
			Location: to.Ptr(v1.LocationGlobal),
			Type:     to.Ptr(cli_credential.AzureCredential),
			Properties: &ucp.AzureServicePrincipalProperties{
				Storage: &ucp.CredentialStorageProperties{
					Kind: to.Ptr(ucp.CredentialStorageKindInternal),
				},
				TenantID:     &r.Options.CloudProviders.Azure.ServicePrincipal.TenantID,
				ClientID:     &r.Options.CloudProviders.Azure.ServicePrincipal.ClientID,
				ClientSecret: &r.Options.CloudProviders.Azure.ServicePrincipal.ClientSecret,
			},
		}, nil
	case azure.AzureCredentialKindWorkloadIdentity:
		return ucp.AzureCredentialResource{
			Location: to.Ptr(v1.LocationGlobal),
			Type:     to.Ptr(cli_credential.AzureCredential),
			Properties: &ucp.AzureWorkloadIdentityProperties{
				Storage: &ucp.CredentialStorageProperties{
					Kind: to.Ptr(ucp.CredentialStorageKindInternal),
				},
				TenantID: &r.Options.CloudProviders.Azure.WorkloadIdentity.TenantID,
				ClientID: &r.Options.CloudProviders.Azure.WorkloadIdentity.ClientID,
			},
		}, nil
	default:
		return ucp.AzureCredentialResource{}, fmt.Errorf("unsupported Azure credential kind: %s", r.Options.CloudProviders.Azure.CredentialKind)
	}
}

func (r *Runner) getAWSCredential() (ucp.AwsCredentialResource, error) {
	switch r.Options.CloudProviders.AWS.CredentialKind {
	case aws.AWSCredentialKindAccessKey:
		return ucp.AwsCredentialResource{
			Location: to.Ptr(v1.LocationGlobal),
			Type:     to.Ptr(cli_credential.AWSCredential),
			Properties: &ucp.AwsAccessKeyCredentialProperties{
				Storage: &ucp.CredentialStorageProperties{
					Kind: to.Ptr(ucp.CredentialStorageKindInternal),
				},
				AccessKeyID:     &r.Options.CloudProviders.AWS.AccessKey.AccessKeyID,
				SecretAccessKey: &r.Options.CloudProviders.AWS.AccessKey.SecretAccessKey,
			},
		}, nil
	case aws.AWSCredentialKindIRSA:
		return ucp.AwsCredentialResource{
			Location: to.Ptr(v1.LocationGlobal),
			Type:     to.Ptr(cli_credential.AWSCredential),
			Properties: &ucp.AwsIRSACredentialProperties{
				Storage: &ucp.CredentialStorageProperties{
					Kind: to.Ptr(ucp.CredentialStorageKindInternal),
				},
				RoleARN: &r.Options.CloudProviders.AWS.IRSA.RoleARN,
			},
		}, nil
	default:
		return ucp.AwsCredentialResource{}, fmt.Errorf("unsupported AWS credential kind: %s", r.Options.CloudProviders.AWS.CredentialKind)
	}
}
