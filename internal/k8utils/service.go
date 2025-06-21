package k8utils

import (
	"context"
	"fmt"
	"log"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// ServiceBuilder provides a builder pattern for Kubernetes Services
type ServiceBuilder struct {
	admin   *Admin
	service *corev1.Service
}

// NewService creates a new ServiceBuilder with name and labels
func (a *Admin) NewService(name string, labels map[string]string) *ServiceBuilder {
	return &ServiceBuilder{
		admin: a,
		service: &corev1.Service{
			ObjectMeta: metav1.ObjectMeta{
				Name:      name,
				Labels:    labels,
				Namespace: a.namespace,
			},
		},
	}
}

// Default sets up basic metadata and labels for the service
func (s *ServiceBuilder) Default() *ServiceBuilder {
	s.service.APIVersion = "v1"
	s.service.Kind = "Service"
	return s
}

// SetSpec sets the service spec (user-provided)
func (s *ServiceBuilder) SetSpec(spec corev1.ServiceSpec) *ServiceBuilder {
	s.service.Spec = spec
	return s
}

// Apply creates the service in the cluster
func (s *ServiceBuilder) Apply() error {
	servicesClient := s.admin.clientset.CoreV1().Services(s.admin.namespace)
	_, err := servicesClient.Create(context.TODO(), s.service, metav1.CreateOptions{})
	if err != nil {
		return fmt.Errorf("failed to apply service: %w", err)
	}
	log.Printf("Service '%s' created in namespace '%s'", s.service.Name, s.admin.namespace)
	return nil
}
