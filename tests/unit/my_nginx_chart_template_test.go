package unit

import (
	"path/filepath"
	"strings"
	"testing"

	"github.com/gruntwork-io/terratest/modules/helm"
	"github.com/gruntwork-io/terratest/modules/k8s"
	"github.com/gruntwork-io/terratest/modules/random"
	"github.com/stretchr/testify/require"
	appsv1 "k8s.io/api/apps/v1"
)

const myChart = "../../charts/my-chart"

func TestTemplateRenderedDeployment(t *testing.T) {
	type args struct {
		kubeVersion   string
		namespace     string
		releaseName   string
		chartRelPath  string
		expectedImage string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "Kubernetes 1.23",
			args: args{
				kubeVersion:   "1.23",
				namespace:     "test-" + strings.ToLower(random.UniqueId()),
				releaseName:   "test-" + strings.ToLower(random.UniqueId()),
				chartRelPath:  myChart,
				expectedImage: "docker.io/nginx:1.21",
			},
		},
		{
			name: "Kubernetes 1.22",
			args: args{
				kubeVersion:   "1.22",
				namespace:     "test-" + strings.ToLower(random.UniqueId()),
				releaseName:   "test-" + strings.ToLower(random.UniqueId()),
				chartRelPath:  myChart,
				expectedImage: "docker.io/nginx:1.21",
			},
		},
		{
			name: "Kubernetes 1.21",
			args: args{
				kubeVersion:   "1.21",
				namespace:     "test-" + strings.ToLower(random.UniqueId()),
				releaseName:   "test-" + strings.ToLower(random.UniqueId()),
				chartRelPath:  myChart,
				expectedImage: "docker.io/nginx:1.21",
			},
		},
		{
			name: "Kubernetes 1.20",
			args: args{
				kubeVersion:   "1.20",
				namespace:     "test-" + strings.ToLower(random.UniqueId()),
				releaseName:   "test-" + strings.ToLower(random.UniqueId()),
				chartRelPath:  myChart,
				expectedImage: "docker.io/nginx:1.20.0",
			},
		},
		{
			name: "Kubernetes 1.19",
			args: args{
				kubeVersion:   "1.19",
				namespace:     "test-" + strings.ToLower(random.UniqueId()),
				releaseName:   "test-" + strings.ToLower(random.UniqueId()),
				chartRelPath:  myChart,
				expectedImage: "docker.io/nginx:1.19",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// arrange
			chartPath, err := filepath.Abs(tt.args.chartRelPath)
			require.NoError(t, err)

			options := &helm.Options{
				KubectlOptions: k8s.NewKubectlOptions("", "", tt.args.namespace),
			}

			// act
			output := helm.RenderTemplate(t, options, chartPath, tt.args.releaseName, []string{"templates/deployment.yaml"}, "--kube-version", tt.args.kubeVersion)

			var deployment appsv1.Deployment
			helm.UnmarshalK8SYaml(t, output, &deployment)

			// assert
			require.Equal(t, tt.args.namespace, deployment.Namespace)
			deploymentSetContainers := deployment.Spec.Template.Spec.Containers
			require.Equal(t, len(deploymentSetContainers), 1)
			require.Equal(t, tt.args.expectedImage, deploymentSetContainers[0].Image)
		})
	}
}
