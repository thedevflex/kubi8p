package handler

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/thedevflex/kubi8p/internal/cache"
	"github.com/thedevflex/kubi8p/internal/kubi8al"
)

func (h *Handler) InitiateDNS(w http.ResponseWriter, r *http.Request) {
	err := h.dns.CreateDNSDeployment()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Println("Error creating dns deployment:", err)
		return
	}
	err = h.dns.CreateDNSService()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Println("Error creating dns service:", err)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("DNS initiated successfully"))
}

func (h *Handler) GetDNSStatus(w http.ResponseWriter, r *http.Request) {
	// reject all other methods
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	externalIP, err := h.dns.GetExternalIP()
	if err != nil {
		switch err.Error() {
		case kubi8al.ErrorDeploymentNotReady:
			http.Error(w, err.Error(), http.StatusInternalServerError)
			log.Println("Error getting dns status:", err)
			return
		case kubi8al.ErrorServiceNotReady:
			http.Error(w, err.Error(), http.StatusInternalServerError)
			log.Println("Error getting dns status:", err)
			return
		case kubi8al.ErrorServiceNotLoadBalancer:
			http.Error(w, err.Error(), http.StatusInternalServerError)
			log.Println("Error getting dns status:", err)
			return
		case kubi8al.ErrorServiceNoExternalIP:
			http.Error(w, err.Error(), http.StatusInternalServerError)
			log.Println("Error getting dns status:", err)
			return
		}
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"status":      "ready",
		"external_ip": externalIP,
		"message":     "DNS service is ready! External IP: " + externalIP,
	})

}

func (h *Handler) ConfigureDNS(w http.ResponseWriter, r *http.Request) {
	payload := cache.DNSPayload{}
	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.cache.SetDNSPayload(payload); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("DNS configured successfully"))
}

func (h *Handler) VerifyDNS(w http.ResponseWriter, r *http.Request) {
	payload := h.cache.GetDNSPayload()
	log.Println("payload", payload)
	url := "https://kubi8al-dns-health-check." + payload.Prefix + "." + payload.Domain + "/health"
	log.Println("url", url)

	resp, err := http.Get(url)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()
}
