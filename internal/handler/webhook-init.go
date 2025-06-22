package handler

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/thedevflex/kubi8p/internal/kubi8al"
)

func (h *Handler) InitiateWebhook(w http.ResponseWriter, r *http.Request) {

	var payload struct {
		WebhookSecret string `json:"webhook_secret"`
	}

	kubi8al.WithPort(8080)(h.webhook)

	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		log.Println("Error decoding webhook secret:", err)
		return
	}
	err := h.webhook.CreateWebhookDeployment(kubi8al.WebhookDeploymentEnvStruct{
		WEBHOOK_SECRET: payload.WebhookSecret,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Println("Error creating webhook deployment:", err)
		return
	}
	err = h.webhook.CreateWebhookService()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Println("Error creating webhook service:", err)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Webhook initiated successfully"))
}
