package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/thedevflex/kubi8p/internal/cache"
)

func (h *Handler) InitiateDB(w http.ResponseWriter, r *http.Request) {
	payload := cache.DBConnectionPayload{}
	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.cache.SetDBConnectionPayload(payload); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Println("Payload: ", payload)

	if payload.Type == "new" {
		if err := h.CreatePostgresDeployment(); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("DB connection payload set"))
}

func (h *Handler) CreatePostgresDeployment() error {
	fmt.Println("Creating Postgres deployment")
	postgresSpec := appsv1.DeploymentSpec{
		Replicas: int32Ptr(1),
		Selector: &metav1.LabelSelector{
			MatchLabels: map[string]string{"app": "postgres"},
		},
		Template: corev1.PodTemplateSpec{
			ObjectMeta: metav1.ObjectMeta{
				Labels: map[string]string{"app": "postgres"},
			},
			Spec: corev1.PodSpec{
				Containers: []corev1.Container{
					{
						Name:  "postgres",
						Image: "postgres:14",
						Env: []corev1.EnvVar{
							{Name: "POSTGRES_USER", Value: "admin"},
							{Name: "POSTGRES_PASSWORD", Value: "admin"},
							{Name: "POSTGRES_DB", Value: "kubi8al"},
						},
						Ports: []corev1.ContainerPort{
							{ContainerPort: 5432},
						},
					},
				},
			},
		},
	}

	return h.admin.NewDeployment("kubi8al-db", map[string]string{"app": "kubi8al"}).Default().SetSpec(postgresSpec).Apply()
}

func int32Ptr(i int32) *int32 {
	return &i
}
