package constants

const (
	Kubi8alNamespace = "kubi8al"
	Kubi8alSecret    = "kubi8al-secret"
	// db
	Kubi8alDBName = "kubi8al-db"

	// Webhook
	Kubi8alWebhookName       = "kubi8al-webhook"
	Kubi8alWebhookTag        = "v2.0.0"
	Kubi8alWebhookPackgeName = "ghcr.io/thedevflex/kubi8al-webhook"

	// DNS
	Kubi8alDNSName       = "kubi8al-dns"
	Kubi8alDNSTag        = "v1.0.0"
	Kubi8alDNSPackgeName = "ghcr.io/thedevflex/kubi8al-dns"

	// Kubi8al
	Kubi8alVersion    = "v1.0.0"
	Kubi8alPackgeName = "ghcr.io/thedevflex/kubi8al-api"
	Kubi8alPort       = 8080
	Kubi8alInKubeIp   = "http://kubi8al:8080"
)
