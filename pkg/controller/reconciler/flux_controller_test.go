/*
Copyright 2024 The Radius Authors.
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

package reconciler

import (
	"context"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"
	"testing"
	"time"

	sourcev1 "github.com/fluxcd/source-controller/api/v1"
	"github.com/radius-project/radius/pkg/cli/bicep"
	"github.com/radius-project/radius/pkg/cli/filesystem"
	radappiov1alpha3 "github.com/radius-project/radius/pkg/controller/api/radapp.io/v1alpha3"
	"github.com/radius-project/radius/pkg/to"
	"github.com/radius-project/radius/test/testcontext"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/wait"
	ctrl "sigs.k8s.io/controller-runtime"
	k8sclient "sigs.k8s.io/controller-runtime/pkg/client"
	crconfig "sigs.k8s.io/controller-runtime/pkg/config"
	"sigs.k8s.io/controller-runtime/pkg/metrics/server"
)

type setupFluxControllerTestOptions struct {
	archiveFetcher ArchiveFetcher
	filesystem     filesystem.FileSystem
	bicep          bicep.Interface
}

func SetupFluxControllerTest(t *testing.T, options setupFluxControllerTestOptions) k8sclient.Client {
	SkipWithoutEnvironment(t)

	// For debugging, you can set uncomment this to see logs from the controller. This will cause tests to fail
	// because the logging will continue after the test completes.
	//
	// Add runtimelog "sigs.k8s.io/controller-runtime/pkg/log" to imports.
	//
	// runtimelog.SetLogger(ucplog.FromContextOrDiscard(testcontext.New(t)))

	// Shut down the manager when the test exits.
	ctx, cancel := testcontext.NewWithCancel(t)
	t.Cleanup(cancel)

	mgr, err := ctrl.NewManager(config, ctrl.Options{
		Scheme: scheme,
		Controller: crconfig.Controller{
			SkipNameValidation: to.Ptr(true),
		},

		// Suppress metrics in tests to avoid conflicts.
		Metrics: server.Options{
			BindAddress: "0",
		},
	})
	require.NoError(t, err)

	// Set up FluxController.
	fluxController := &FluxController{
		Client:         mgr.GetClient(),
		FileSystem:     options.filesystem,
		Bicep:          options.bicep,
		ArchiveFetcher: options.archiveFetcher,
	}
	err = (fluxController).SetupWithManager(mgr)
	require.NoError(t, err)

	go func() {
		err := mgr.Start(ctx)
		require.NoError(t, err)
	}()

	return mgr.GetClient()
}

type Step struct {
	Path            string
	BicepFiles      []string
	BicepParamFiles []string
}

func RunFluxControllerTest(t *testing.T, steps []Step) {
	ctx := testcontext.New(t)

	// Set up mocks
	mctrl := gomock.NewController(t)
	defer mctrl.Finish()

	fs := filesystem.NewMemMapFileSystem(nil)
	archiveFetcher := NewMockArchiveFetcher(mctrl)
	bicep := bicep.NewMockInterface(mctrl)

	options := setupFluxControllerTestOptions{
		archiveFetcher: archiveFetcher,
		filesystem:     fs,
		bicep:          bicep,
	}

	k8sClient := SetupFluxControllerTest(t, options)

	testGitRepoName := "example-repo"
	testGitRepoURL := "https://github.com/radius-project/example-repo.git"
	testGitRepoSHA := "sha256:1234"

	var archiveFetcherCalls []any
	var bicepCalls []any
	for _, step := range steps {
		s := step // capture loop variable
		call := archiveFetcher.EXPECT().
			Fetch(testGitRepoURL, testGitRepoSHA, gomock.Any()).
			Return(nil).
			Times(1).
			Do(func(archiveURL, digest, dir string) {
				err := filepath.WalkDir(s.Path, func(srcPath string, info os.DirEntry, err error) error {
					if err != nil {
						return err
					}
					relPath, err := filepath.Rel(s.Path, srcPath)
					if err != nil {
						return err
					}
					dstPath := filepath.Join(dir, relPath)
					if info.IsDir() {
						return fs.MkdirAll(dstPath, 0755)
					}
					// Read contents of srcFile and write to destination.
					data, err := os.ReadFile(srcPath)
					if err != nil {
						return err
					}
					return fs.WriteFile(dstPath, data, 0644)
				})
				require.NoError(t, err)
			})

		archiveFetcherCalls = append(archiveFetcherCalls, call)

		radiusConfig, err := ParseRadiusGitOpsConfig(path.Join(s.Path, "radius-gitops-config.yaml"))
		require.NoError(t, err)
		require.NotNil(t, radiusConfig)

		for _, configEntry := range radiusConfig.Config {
			ce := configEntry // capture loop variable
			bicepBuildCall := bicep.EXPECT().
				Build(gomock.Any(), "--outfile", gomock.Any()).
				Return(nil, nil).
				Times(1).
				Do(func(args ...string) {
					filePath := args[0]
					outFilePath := args[2]
					outFileName := filepath.Base(outFilePath)
					localFilePath := s.Path
					fileContent, err := os.ReadFile(path.Join(localFilePath, outFileName))
					require.NoError(t, err)
					err = fs.WriteFile(filePath, fileContent, 0644)
					require.NoError(t, err)
				})
			bicepCalls = append(bicepCalls, bicepBuildCall)

			if ce.Params != "" {
				bicepBuildParamsCall := bicep.EXPECT().
					BuildParams(gomock.Any(), "--outfile", gomock.Any()).
					Return(nil, nil).
					Times(1).
					Do(func(args ...string) {
						filePath := args[0]
						outFilePath := args[2]
						outFileName := filepath.Base(outFilePath)
						localFilePath := s.Path
						fileContent, err := os.ReadFile(path.Join(localFilePath, outFileName))
						require.NoError(t, err)
						err = fs.WriteFile(filePath, fileContent, 0644)
						require.NoError(t, err)
					})
				bicepCalls = append(bicepCalls, bicepBuildParamsCall)
			}
		}
	}
	gomock.InOrder(archiveFetcherCalls...)
	gomock.InOrder(bicepCalls...)

	for stepIndex, step := range steps {
		stepNumber := stepIndex + 1
		radiusConfig, err := ParseRadiusGitOpsConfig(path.Join(step.Path, "radius-gitops-config.yaml"))
		require.NoError(t, err)
		require.NotNil(t, radiusConfig)

		for _, configEntry := range radiusConfig.Config {
			name := configEntry.Name
			nameBase := strings.TrimSuffix(name, path.Ext(name))
			namespaceName := configEntry.Namespace
			if namespaceName == "" {
				namespaceName = nameBase
			}

			// Create a namespace if it does not exist
			if err := k8sClient.Get(ctx, types.NamespacedName{Name: namespaceName}, &corev1.Namespace{}); err != nil {
				if k8sclient.IgnoreNotFound(err) != nil {
					require.NoError(t, err)
				}
				namespace := &corev1.Namespace{
					ObjectMeta: metav1.ObjectMeta{
						Name: namespaceName,
					},
				}
				err := k8sClient.Create(ctx, namespace)
				require.NoError(t, err)
			}

			gitRepo := sourcev1.GitRepository{}
			gitRepoNamespacedName := types.NamespacedName{Name: testGitRepoName, Namespace: namespaceName}
			if stepNumber == 1 {
				// Create a GitRepository resource on the cluster
				gitRepo = makeGitRepository(gitRepoNamespacedName, testGitRepoURL)
				err = k8sClient.Create(ctx, &gitRepo)
				require.NoError(t, err)

				// Wait for the GitRepository to be created
				err = waitForGitRepositoryToExistWithGeneration(ctx, k8sClient, gitRepoNamespacedName, &gitRepo, int64(stepNumber))
				require.NoError(t, err, "GitRepository was not created successfully")
			}

			// Fetch the latest GitRepository object
			err = k8sClient.Get(ctx, gitRepoNamespacedName, &gitRepo)
			require.NoError(t, err)

			// Update the Status subresource
			updatedStatus := sourcev1.GitRepositoryStatus{
				ObservedGeneration: int64(stepNumber),
				Artifact: &sourcev1.Artifact{
					URL:      testGitRepoURL,
					Digest:   testGitRepoSHA,
					Revision: fmt.Sprintf("v%d", stepNumber),
					LastUpdateTime: metav1.Time{
						Time: time.Now(),
					},
				},
				Conditions: []metav1.Condition{
					{
						Type:    "Ready",
						Status:  metav1.ConditionTrue,
						Reason:  "Succeeded",
						Message: "Repository is ready",
						LastTransitionTime: metav1.Time{
							Time: time.Now(),
						},
					},
				},
			}
			gitRepo.Status = updatedStatus
			err = k8sClient.Status().Update(ctx, &gitRepo)
			require.NoError(t, err)

			// Now, the FluxController should reconcile the GitRepository and create the DeploymentTemplate resource.
			deploymentTemplateName := name
			deploymentTemplateNamespacedName := types.NamespacedName{Name: deploymentTemplateName, Namespace: namespaceName}
			deploymentTemplate := radappiov1alpha3.DeploymentTemplate{}
			err = waitForDeploymentTemplateToExistWithGeneration(ctx, k8sClient, deploymentTemplateNamespacedName, &deploymentTemplate, int64(stepNumber))
			require.NoError(t, err)

			// Fetch the latest DeploymentTemplate object
			err = k8sClient.Get(ctx, deploymentTemplateNamespacedName, &deploymentTemplate)
			require.NoError(t, err)

		}
	}
}

func Test_FluxController_Basic(t *testing.T) {
	steps := []Step{
		{
			Path: "testdata/flux-basic",
		},
	}

	RunFluxControllerTest(t, steps)
}

func Test_FluxController_Update(t *testing.T) {
	steps := []Step{
		{
			Path: "testdata/flux-update/step-1",
		},
		{
			Path: "testdata/flux-update/step-2",
		},
	}

	RunFluxControllerTest(t, steps)
}

func makeGitRepository(namespacedName types.NamespacedName, url string) sourcev1.GitRepository {
	return sourcev1.GitRepository{
		TypeMeta: metav1.TypeMeta{
			Kind:       "GitRepository",
			APIVersion: "source.toolkit.fluxcd.io/v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      namespacedName.Name,
			Namespace: namespacedName.Namespace,
		},
		Spec: sourcev1.GitRepositorySpec{
			URL: url,
		},
	}
}

func waitForGitRepositoryToExistWithGeneration(ctx context.Context, k8sClient k8sclient.Client, key k8sclient.ObjectKey, obj k8sclient.Object, generation int64) error {
	timeout := 10 * time.Second
	interval := 1 * time.Second
	deadlineCtx, deadlineCancel := context.WithTimeout(ctx, timeout)
	defer deadlineCancel()

	err := wait.PollUntilContextTimeout(deadlineCtx, interval, timeout, true, func(_ context.Context) (bool, error) {
		err := k8sClient.Get(ctx, key, obj)
		if err != nil {
			return false, nil // Continue polling
		}

		gitRepo, ok := obj.(*sourcev1.GitRepository)
		if !ok {
			return false, nil // Continue polling
		}

		if gitRepo.Generation != generation {
			return false, nil // Continue polling
		}

		return true, nil // Found
	})

	if err != nil {
		return fmt.Errorf("GitRepository %s/%s was not created successfully", key.Namespace, key.Name)
	}

	return nil
}

func waitForDeploymentTemplateToExistWithGeneration(ctx context.Context, k8sClient k8sclient.Client, key types.NamespacedName, obj k8sclient.Object, generation int64) error {
	timeout := 10 * time.Second
	interval := 1 * time.Second
	deadlineCtx, deadlineCancel := context.WithTimeout(ctx, timeout)
	defer deadlineCancel()

	err := wait.PollUntilContextTimeout(deadlineCtx, interval, timeout, true, func(_ context.Context) (bool, error) {
		err := k8sClient.Get(ctx, key, obj)
		if err != nil {
			return false, nil // Continue polling
		}

		deploymentTemplate, ok := obj.(*radappiov1alpha3.DeploymentTemplate)
		if !ok {
			return false, nil // Continue polling
		}

		if deploymentTemplate.Generation != generation {
			return false, nil // Continue polling
		}

		return true, nil // Found
	})

	if err != nil {
		return fmt.Errorf("DeploymentTemplate %s/%s was not created successfully", key.Namespace, key.Name)
	}

	return nil
}
