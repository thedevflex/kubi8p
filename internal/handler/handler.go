package handler

import (
	"github.com/thedevflex/kubi8p/internal/cache"
	"github.com/thedevflex/kubi8p/internal/k8utils"
)

type Handler struct {
	cache *cache.InstallerCache
	admin *k8utils.Admin
}

func NewHandler(cache *cache.InstallerCache, admin *k8utils.Admin) *Handler {
	return &Handler{cache: cache, admin: admin}
}
