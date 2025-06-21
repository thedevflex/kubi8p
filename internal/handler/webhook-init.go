package handler

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/thedevflex/kubi8p/internal/kubi8al"
)

func (h *Handler) InitiateWebhook(w http.ResponseWriter, r *http.Request) {
	webhook := kubi8al.NewWebhook(h.admin)
	var payload struct {
		WebhookSecret string `json:"webhook_secret"`
	}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		log.Println("Error decoding webhook secret:", err)
		return
	}
	err := webhook.CreateWebhookDeployment(kubi8al.WebhookDeploymentEnvStruct{
		WEBHOOK_SECRET: payload.WebhookSecret,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Println("Error creating webhook deployment:", err)
		return
	}
	err = webhook.CreateWebhookService()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Println("Error creating webhook service:", err)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Webhook initiated successfully"))
}
