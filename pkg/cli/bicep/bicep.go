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

package bicep

import (
	"context"
	"fmt"
	"os"

	"github.com/radius-project/radius/pkg/cli/bicep/tools"
	biceptools "github.com/radius-project/radius/pkg/cli/bicep/tools"
	"github.com/radius-project/radius/pkg/retry"
)

const (
	radBicepEnvVar                  = "RAD_BICEP"
	binaryName                      = "rad-bicep"
	manifestToBicepExtensionCLIName = "manifest-to-bicep-extension"
)

func GetBicepFilePath() (string, error) {
	return tools.GetLocalFilepath(radBicepEnvVar, binaryName)
}

// GetBicepManifestToBicepExtensionCLIPath returns the path to the manifest-to-bicep-extension CLI tool.
func GetBicepManifestToBicepExtensionCLIPath() (string, error) {
	// The manifest-to-bicep-extension CLI is expected to be in the same directory as the Bicep binary.
	filepath, err := GetBicepFilePath()
	if err != nil {
		return "", fmt.Errorf("failed to get Bicep file path: %v", err)
	}

	// Remove the binary name and append the manifest-to-bicep-extension CLI name.
	filepath = filepath[:len(filepath)-len(binaryName)]

	// Append the manifest-to-bicep-extension CLI name.
	filepath = filepath + manifestToBicepExtensionCLIName

	return filepath, nil
}

// IsBicepInstalled checks if the Bicep binary is installed on the local machine.
func IsBicepInstalled() (bool, error) {
	filepath, err := GetBicepFilePath()
	if err != nil {
		return false, err
	}

	_, err = os.Stat(filepath)
	if err != nil && os.IsNotExist(err) {
		return false, nil
	} else if err != nil {
		return false, fmt.Errorf("error checking for %s: %v", filepath, err)
	}

	return true, nil
}

// DeleteBicep removes the Bicep binary from the local machine.
func DeleteBicep() error {
	filepath, err := GetBicepFilePath()
	if err != nil {
		return err
	}

	err = os.Remove(filepath)
	if err != nil {
		return fmt.Errorf("failed to delete %s: %v", filepath, err)
	}

	return nil
}

// DownloadBicep downloads the Bicep binary to the local machine.
func DownloadBicep() error {
	ctx := context.Background()
	retryer := retry.NewDefaultRetryer()

	// Download the Bicep binary
	filepath, err := GetBicepFilePath()
	if err != nil {
		return err
	}

	err = retryer.RetryFunc(ctx, func(ctx context.Context) error {
		return biceptools.DownloadToFolder(filepath)
	})
	if err != nil {
		return fmt.Errorf("failed to download bicep: %v", err)
	}

	// Download the manifest-to-bicep-extension CLI
	manifestToBicepExtensionCLIPath, err := GetBicepManifestToBicepExtensionCLIPath()
	if err != nil {
		return err
	}

	err = retryer.RetryFunc(ctx, func(ctx context.Context) error {
		return biceptools.DownloadManifestToBicepExtensionCLI(manifestToBicepExtensionCLIPath)
	})
	if err != nil {
		return fmt.Errorf("failed to download manifest-to-bicep-extension CLI: %v", err)
	}

	return nil
}
