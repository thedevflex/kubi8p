package k8utils

import (
	"context"
	"fmt"
	"log"

	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// DeploymentBuilder provides a builder pattern for Kubernetes Deployments
type DeploymentBuilder struct {
	admin      *Admin
	deployment *appsv1.Deployment
}

// NewDeployment creates a new DeploymentBuilder with name and labels
func (a *Admin) NewDeployment(name string, labels map[string]string) *DeploymentBuilder {
	return &DeploymentBuilder{
		admin: a,
		deployment: &appsv1.Deployment{
			ObjectMeta: metav1.ObjectMeta{
				Name:      name,
				Labels:    labels,
				Namespace: a.namespace,
			},
		},
	}
}

// Default sets up basic metadata and labels for the deployment
func (d *DeploymentBuilder) Default() *DeploymentBuilder {
	d.deployment.APIVersion = "apps/v1"
	d.deployment.Kind = "Deployment"
	return d
}

// SetSpec sets the deployment spec (user-provided)
func (d *DeploymentBuilder) SetSpec(spec appsv1.DeploymentSpec) *DeploymentBuilder {
	d.deployment.Spec = spec
	return d
}

// Apply creates the deployment in the cluster
func (d *DeploymentBuilder) Apply() error {
	deploymentsClient := d.admin.clientset.AppsV1().Deployments(d.admin.namespace)
	_, err := deploymentsClient.Create(context.TODO(), d.deployment, metav1.CreateOptions{})
	if err != nil {
		return fmt.Errorf("failed to apply deployment: %w", err)
	}
	log.Printf("Deployment '%s' created in namespace '%s'", d.deployment.Name, d.admin.namespace)
	return nil
}
