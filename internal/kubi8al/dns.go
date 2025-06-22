package kubi8al

import (
	"errors"
	"strconv"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"

	"github.com/thedevflex/kubi8p/internal/constants"
	"github.com/thedevflex/kubi8p/internal/k8utils"
)

type DNS struct {
	admin *k8utils.Admin
}

func NewDNS(admin *k8utils.Admin) *DNS {
	return &DNS{admin: admin}
}

func (d *DNS) CreateDNSService() error {
	port := int32(8080)
	err := d.admin.NewService(constants.Kubi8alDNSName, map[string]string{
		"app": constants.Kubi8alDNSName,
	}).Default().SetSpec(corev1.ServiceSpec{
		Type: corev1.ServiceTypeLoadBalancer,
		Selector: map[string]string{
			"app": constants.Kubi8alDNSName,
		},
		Ports: []corev1.ServicePort{
			{
				Port:       80,
				TargetPort: intstr.FromInt(int(port)),
				Protocol:   corev1.ProtocolTCP,
				NodePort:   constants.Kubi8alDNSNodePort,
			},
		},
	}).Apply()
	return err
}

func (d *DNS) CreateDNSDeployment() error {

	containerPort := int32(8080)

	envs := []corev1.EnvVar{
		{
			Name:  "PORT",
			Value: strconv.Itoa(int(containerPort)),
		},
	}

	err := d.admin.NewDeployment(constants.Kubi8alDNSName, map[string]string{
		"app": constants.Kubi8alDNSName,
	}).Default().SetSpec(appsv1.DeploymentSpec{
		Replicas: &[]int32{1}[0],
		Selector: &metav1.LabelSelector{
			MatchLabels: map[string]string{
				"app": constants.Kubi8alDNSName,
			},
		},
		Template: corev1.PodTemplateSpec{
			ObjectMeta: metav1.ObjectMeta{
				Labels: map[string]string{
					"app": constants.Kubi8alDNSName,
				},
			},
			Spec: corev1.PodSpec{
				Containers: []corev1.Container{
					{
						Name:  "dns",
						Image: constants.Kubi8alDNSPackgeName + ":" + constants.Kubi8alDNSTag,
						Ports: []corev1.ContainerPort{
							{
								ContainerPort: int32(containerPort),
							},
						},
						Env: envs,
					},
				},
			},
		},
	}).Apply()
	return err
}

var (
	ErrorDeploymentNotReady     = "deployment is not ready yet"
	ErrorServiceNotReady        = "service is not ready yet"
	ErrorServiceNotLoadBalancer = "service is not of LoadBalancer type"
	ErrorServiceNoExternalIP    = "service external IP is not available"
)

// GetExternalIP returns the external IP of the DNS service
// Returns error if service or deployment is not ready
func (d *DNS) GetExternalIP() (string, error) {
	// First check if deployment is ready
	deployment, err := d.admin.GetDeploymentStatus(constants.Kubi8alDNSName)
	if err != nil {
		return "", errors.New(ErrorDeploymentNotReady)
	}

	// Check if deployment is ready
	if deployment.Status.ReadyReplicas == 0 || deployment.Status.ReadyReplicas < *deployment.Spec.Replicas {
		return "", errors.New(ErrorDeploymentNotReady)
	}

	// Check if deployment has any failed conditions
	for _, condition := range deployment.Status.Conditions {
		if condition.Type == appsv1.DeploymentReplicaFailure && condition.Status == corev1.ConditionTrue {
			return "", errors.New(ErrorDeploymentNotReady)
		}
	}

	// Get service status
	service, err := d.admin.GetServiceStatus(constants.Kubi8alDNSName)
	if err != nil {
		return "", errors.New(ErrorServiceNotReady)
	}

	// Check if service is of LoadBalancer type
	if service.Spec.Type != corev1.ServiceTypeLoadBalancer {
		return "", errors.New(ErrorServiceNotLoadBalancer)
	}

	// Check if service has external IP
	if len(service.Status.LoadBalancer.Ingress) == 0 {
		return "", errors.New(ErrorServiceNoExternalIP)
	}

	// Return the external IP
	ingress := service.Status.LoadBalancer.Ingress[0]
	if ingress.IP != "" {
		return ingress.IP, nil
	} else if ingress.Hostname != "" {
		return ingress.Hostname, nil
	}

	return "", errors.New(ErrorServiceNoExternalIP)
}
