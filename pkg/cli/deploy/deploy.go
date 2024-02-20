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

package deploy

import (
	"context"
	"fmt"
	"os/signal"
	"sync"

	"github.com/radius-project/radius/pkg/cli/bicep"
	"github.com/radius-project/radius/pkg/cli/clients"
	"github.com/radius-project/radius/pkg/cli/output"
)

// DeployWithProgress runs a deployment and displays progress to the user. This is intended to be used
// from the CLI and thus logs to the console.
//

// DeployWithProgress injects environment and application parameters into the template, displays progress updates while
// deploying, and logs the deployment results and public endpoints. If an error occurs, an error is returned.
func DeployWithProgress(ctx context.Context, options Options) (clients.DeploymentResult, error) {
	deploymentClient, err := options.ConnectionFactory.CreateDeploymentClient(ctx, options.Workspace)
	if err != nil {
		return clients.DeploymentResult{}, err
	}

	err = bicep.InjectEnvironmentParam(options.Template, options.Parameters, options.Providers.Radius.EnvironmentID)
	if err != nil {
		return clients.DeploymentResult{}, err
	}

	err = bicep.InjectApplicationParam(options.Template, options.Parameters, options.Providers.Radius.ApplicationID)
	if err != nil {
		return clients.DeploymentResult{}, err
	}

	step := output.BeginStep(options.ProgressText)
	output.LogInfo("")

	// Watch for progress while we're deploying.
	progressChan := make(chan clients.ResourceProgress, 1)
	listener := NewProgressListener(progressChan)
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		listener.Run()
		wg.Done()
	}()

	result, err := deploymentClient.Deploy(ctx, clients.DeploymentOptions{
		Template:     options.Template,
		Parameters:   options.Parameters,
		Providers:    options.Providers,
		ProgressChan: progressChan,
	})

	// Drain any UI progress updates before we process the results of the deployment.
	wg.Wait()
	if err != nil {
		return clients.DeploymentResult{}, err
	}

	output.LogInfo("")
	output.CompleteStep(step)

	output.LogInfo(options.CompletionText)
	output.LogInfo("")

	if len(result.Resources) > 0 {
		output.LogInfo("Resources:")

		for _, resource := range result.Resources {
			if output.ShowResource(resource) {
				output.LogInfo("    " + output.FormatResourceForDisplay(resource))
			}
		}
		var diagnosticsClient clients.DiagnosticsClient
		diagnosticsClient, err = options.ConnectionFactory.CreateDiagnosticsClient(ctx, options.Workspace)
		if err != nil {
			return clients.DeploymentResult{}, err
		}

		endpoints, err := FindPublicEndpoints(ctx, diagnosticsClient, result)
		if err != nil {
			return clients.DeploymentResult{}, err
		}

		if len(endpoints) > 0 {
			output.LogInfo("")
			output.LogInfo("Public Endpoints:")

			for _, entry := range endpoints {
				output.LogInfo("    %s %s", output.FormatResourceForDisplay(entry.Resource), entry.Endpoint)
			}
		}

		exposeDashboardOptions := clients.ExposeDashboardOptions{
			Port:       7007,
			RemotePort: 7007,
		}
		failed, stop, signals, err := diagnosticsClient.ExposeDashboard(ctx, exposeDashboardOptions)
		if err != nil {
			// We don't want to fail deployment if the dashboard fails to launch.
			output.LogInfo("")
			output.LogInfo("Dashboard failed to launch. You may be able to access it by running:")
			output.LogInfo("    kubectl port-forward --namespace=radius-system svc/dashboard 7007:80")
		}

		go func() error {
			// We own stopping the signal created by Expose
			defer signal.Stop(signals)

			for {
				select {
				case <-signals:
					// shutting down... wait for socket to close
					close(stop)
					continue
				case err := <-failed:
					if err != nil {
						return fmt.Errorf("failed to port-forward: %w", err)
					}

					return nil
				}
			}
		}()

		if err != nil {
			// We don't want to fail deployment if the dashboard fails to launch.
			output.LogInfo("")
			output.LogInfo("Dashboard failed to launch. You may be able to access it by running:")
			output.LogInfo("    kubectl port-forward --namespace=radius-system svc/dashboard 7007:80")
		} else {
			output.LogInfo("")
			output.LogInfo("Dashboard URL: %s", "http://localhost:7007")
		}
	}

	return result, nil
}
