package kubi8al

import (
	"strconv"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"

	"github.com/thedevflex/kubi8p/internal/constants"
	"github.com/thedevflex/kubi8p/internal/k8utils"
)

type Webhook struct {
	admin *k8utils.Admin
	port  int32
}

func NewWebhook(admin *k8utils.Admin, port ...int32) *Webhook {
	webhookPort := int32(8080)
	if len(port) > 0 {
		webhookPort = port[0]
	}
	return &Webhook{
		admin: admin,
		port:  webhookPort,
	}
}

func WithPort(port int32) func(*Webhook) {
	return func(w *Webhook) {
		w.port = port
	}
}

func (w *Webhook) CreateWebhookService() error {
	err := w.admin.NewService(constants.Kubi8alWebhookName, map[string]string{
		"app": constants.Kubi8alWebhookName,
	}).Default().SetSpec(corev1.ServiceSpec{
		Type: corev1.ServiceTypeLoadBalancer,
		Selector: map[string]string{
			"app": constants.Kubi8alWebhookName,
		},
		Ports: []corev1.ServicePort{
			{
				Port:       80,
				TargetPort: intstr.FromInt(int(w.port)),
				Protocol:   corev1.ProtocolTCP,
				NodePort:   constants.Kubi8alWebhookNodePort,
			},
		},
	}).Apply()
	return err
}

type WebhookDeploymentEnvStruct struct {
	WEBHOOK_SECRET string
	WEBHOOK_PORT   string
}

func (w *Webhook) CreateWebhookDeployment(env WebhookDeploymentEnvStruct) error {
	if env.WEBHOOK_PORT != "" {
		port, err := strconv.Atoi(env.WEBHOOK_PORT)
		if err != nil {
			return err
		}
		w.port = int32(port)
	}

	envs := []corev1.EnvVar{
		{
			Name:  "EMMITER_API_ADDRESS",
			Value: constants.Kubi8alInKubeIp,
		},
		{
			Name:  "WEBHOOK_SECRET",
			Value: env.WEBHOOK_SECRET,
		},
		{
			Name:  "WEBHOOK_PORT",
			Value: strconv.Itoa(int(w.port)),
		},
	}

	err := w.admin.NewDeployment(constants.Kubi8alWebhookName, map[string]string{
		"app": constants.Kubi8alWebhookName,
	}).Default().SetSpec(appsv1.DeploymentSpec{
		Replicas: &[]int32{1}[0],
		Selector: &metav1.LabelSelector{
			MatchLabels: map[string]string{
				"app": constants.Kubi8alWebhookName,
			},
		},
		Template: corev1.PodTemplateSpec{
			ObjectMeta: metav1.ObjectMeta{
				Labels: map[string]string{
					"app": constants.Kubi8alWebhookName,
				},
			},
			Spec: corev1.PodSpec{
				Containers: []corev1.Container{
					{
						Name:  "webhook",
						Image: constants.Kubi8alWebhookPackgeName + ":" + constants.Kubi8alWebhookTag,
						Ports: []corev1.ContainerPort{
							{
								ContainerPort: w.port,
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
