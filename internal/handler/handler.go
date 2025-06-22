package handler

import (
	"github.com/thedevflex/kubi8p/internal/cache"
	"github.com/thedevflex/kubi8p/internal/k8utils"
	"github.com/thedevflex/kubi8p/internal/kubi8al"
)

type Handler struct {
	cache   *cache.InstallerCache
	admin   *k8utils.Admin
	dns     *kubi8al.DNS
	webhook *kubi8al.Webhook
}

func NewHandler(cache *cache.InstallerCache, admin *k8utils.Admin) *Handler {
	return &Handler{cache: cache, admin: admin, dns: kubi8al.NewDNS(admin), webhook: kubi8al.NewWebhook(admin)}
}
