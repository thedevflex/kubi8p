package k8utils

import (
	"context"
	"fmt"
	"os"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

// Admin provides methods to manage deployments, services, and secrets.
type Admin struct {
	clientset *kubernetes.Clientset
	namespace string
}

// NewAdmin initializes the Admin with in-cluster config
func NewAdmin() (*Admin, error) {

	config, err := parseKubeConfig(os.Getenv("KUBECONFIG_PATH"))

	if err != nil {
		return nil, err
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	return &Admin{
		clientset: clientset,
		namespace: "default", // can be set dynamically later if needed
	}, nil
}

func parseKubeConfig(path string) (*rest.Config, error) {
	if path != "" {
		config, err := clientcmd.BuildConfigFromFlags("", path)
		if err != nil {
			return nil, err
		}
		return config, nil
	}
	config, err := rest.InClusterConfig()
	if err != nil {
		return nil, err
	}
	return config, nil
}

// GetServiceStatus retrieves the service and returns its status
func (a *Admin) GetServiceStatus(serviceName string) (*corev1.Service, error) {
	servicesClient := a.clientset.CoreV1().Services(a.namespace)
	service, err := servicesClient.Get(context.TODO(), serviceName, metav1.GetOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to get service %s: %w", serviceName, err)
	}
	return service, nil
}

// GetDeploymentStatus retrieves the deployment and returns its status
func (a *Admin) GetDeploymentStatus(deploymentName string) (*appsv1.Deployment, error) {
	deploymentsClient := a.clientset.AppsV1().Deployments(a.namespace)
	deployment, err := deploymentsClient.Get(context.TODO(), deploymentName, metav1.GetOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to get deployment %s: %w", deploymentName, err)
	}
	return deployment, nil
}
