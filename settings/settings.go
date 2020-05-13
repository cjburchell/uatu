package settings

import (
	"github.com/cjburchell/settings-go"
)

//Get the application configuration
func Get(config settings.ISettings) AppConfig {
	return AppConfig{
		ConfigFile:     config.Get("CONFIG_FILE", "/config.json"),
		UseNats:        config.GetBool("USE_NATS", true),
		NatsURL:        config.Get("NATS_URL", "tcp://nats:4222"),
		UseRest:        config.GetBool("USE_REST", false),
		RestPort:       config.GetInt("REST_PORT", 8081),
		RestToken:      config.Get("REST_TOKEN", "token"),
		PortalEnable:   config.GetBool("PORTAL_ENABLE", false),
		PortalUsername: config.Get("ADMIN_USER", "admin"),
		PortalPassword: config.Get("ADMIN_PASSWORD", "admin"),
		PortalPort:     config.GetInt("PORTAL_PORT", 8080),
	}
}

// AppConfig object
type AppConfig struct {
	ConfigFile     string
	UseNats        bool
	NatsURL        string
	UseRest        bool
	RestPort       int
	RestToken      string
	PortalEnable   bool
	PortalUsername string
	PortalPassword string
	PortalPort     int
}
