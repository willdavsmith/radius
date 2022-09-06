// ------------------------------------------------------------
// Copyright (c) Microsoft Corporation.
// Licensed under the MIT License.
// ------------------------------------------------------------

package radcli

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	armrpcv1 "github.com/project-radius/radius/pkg/armrpc/api/v1"
	"github.com/project-radius/radius/pkg/cli/clients_new/generated"
	"github.com/project-radius/radius/pkg/cli/framework"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/require"
)

type ValidateInput struct {
	Name          string
	Input         []string
	ExpectedValid bool
	ConfigHolder  framework.ConfigHolder
}

func SharedCommandValidation(t *testing.T, factory func(framework framework.Factory) (*cobra.Command, framework.Runner)) {
	cmd, _ := factory(&framework.Impl{})
	require.NotNil(t, cmd.Args, "Args is required")
	require.NotEmpty(t, cmd.Example, "Example is required")
	require.NotEmpty(t, cmd.Long, "Long is required")
	require.NotEmpty(t, cmd.Short, "Short is required")
	require.NotEmpty(t, cmd.Use, "Use is required")
	require.NotNil(t, cmd.RunE, "RunE is required")
}

func SharedValidateValidation(t *testing.T, factory func(framework framework.Factory) (*cobra.Command, framework.Runner), testcases []ValidateInput) {
	for _, testcase := range testcases {
		t.Run(testcase.Name, func(t *testing.T) {
			framework := &framework.Impl{
				ConnectionFactory: nil,
				ConfigHolder:      &testcase.ConfigHolder,
				Output:            nil}
			cmd, runner := factory(framework)
			cmd.SetArgs(testcase.Input)

			err := cmd.ParseFlags(testcase.Input)
			require.NoError(t, err, "flag parsing failed")

			err = cmd.ValidateArgs(cmd.Flags().Args())
			if testcase.ExpectedValid && err != nil {
				require.NoError(t, err, "validation should have passed but it failed")
			} else if !testcase.ExpectedValid && err != nil {
				return
			}

			err = runner.Validate(cmd, cmd.Flags().Args())
			if testcase.ExpectedValid {
				require.NoError(t, err, "validation should have passed but it failed")
			} else {
				require.Error(t, err, "validation should have failed but it passed")
			}
		})
	}
}

const (
	TestWorkspaceName = "test-workspace"
)

func LoadConfigWithWorkspace(t *testing.T) *viper.Viper {

	var yamlData = []byte(`
workspaces: 
  default: test-workspace
  items: 
    test-workspace: 
      connection: 
        context: test-context
        kind: kubernetes
      environment: /planes/radius/local/resourceGroups/test-resource-group/providers/Applications.Core/environments/test-environment
      scope: /planes/radius/local/resourceGroups/test-resource-group
`)

	v := viper.New()
	v.SetConfigType("yaml")
	err := v.ReadConfig(bytes.NewBuffer(yamlData))
	require.NoError(t, err)
	return v
}

func LoadConfigWithoutWorkspace(t *testing.T) *viper.Viper {

	var yamlData = []byte(`
workspaces: 
`)

	v := viper.New()
	v.SetConfigType("yaml")
	err := v.ReadConfig(bytes.NewBuffer(yamlData))
	require.NoError(t, err)
	return v
}

func Create404Error() error {
	code := armrpcv1.CodeNotFound
	return &azcore.ResponseError{
		ErrorCode:  code,
		StatusCode: 404,
	}
}

func CreateResource(resourceType string, resourceName string) generated.GenericResource {
	id := fmt.Sprintf("/planes/radius/local/resourcegroups/test-environment/providers/%s/%s", resourceType, resourceName)
	location := "global"

	return generated.GenericResource{
		ID:       &id,
		Name:     &resourceName,
		Type:     &resourceType,
		Location: &location,
	}
}