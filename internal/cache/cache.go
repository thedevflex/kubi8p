package cache

import (
	"log"
	"sync"

	"github.com/go-playground/validator/v10"
)

type DBConnectionPayload struct {
	Type             string `json:"type" validate:"required,oneof=new external"`
	ConnectionType   string `json:"connectionType" validate:"omitempty,required_if=Type external,oneof=secret string"`
	SecretName       string `json:"secretName" validate:"omitempty,required_if=ConnectionType secret"`
	SecretPath       string `json:"secretPath" validate:"omitempty,required_if=ConnectionType secret"`
	ConnectionString string `json:"connectionString" validate:"omitempty,required_if=ConnectionType string"`
}

type DNSPayload struct {
	Prefix string `json:"prefix" validate:"required"`
	Domain string `json:"domain" validate:"required"`
}

type InstallerCache struct {
	mutex      sync.RWMutex
	dbPayload  DBConnectionPayload
	dnsPayload DNSPayload
}

func NewInstallerCache() *InstallerCache {
	return &InstallerCache{}
}

func (c *InstallerCache) SetDBConnectionPayload(payload DBConnectionPayload) error {

	if ok, err := c.validateDBConnectionPayload(payload); !ok {
		log.Printf("Invalid DB connection payload: %v", err)
		return err
	}

	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.dbPayload = payload

	return nil
}

func (c *InstallerCache) GetDBConnectionPayload() DBConnectionPayload {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	return c.dbPayload
}

func (c *InstallerCache) validateDBConnectionPayload(payload DBConnectionPayload) (bool, error) {
	validate := validator.New()
	err := validate.Struct(payload)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (c *InstallerCache) validateDNSPayload(payload DNSPayload) (bool, error) {
	validate := validator.New()
	err := validate.Struct(payload)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (c *InstallerCache) GetDNSPayload() DNSPayload {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	return c.dnsPayload
}

func (c *InstallerCache) SetDNSPayload(payload DNSPayload) error {

	if ok, err := c.validateDNSPayload(payload); !ok {
		log.Printf("Invalid DNS payload: %v", err)
		return err
	}

	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.dnsPayload = payload

	return nil
}
