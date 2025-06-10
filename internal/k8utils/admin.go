package k8utils

import (
	"os"

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
