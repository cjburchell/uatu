package settings

import "github.com/cjburchell/tools-go/env"

var (
	ConfigFile     = env.Get("CONFIG_FILE", "/config.json")
	SlackToken     = env.Get("SLACK_TOKEN", "")
	UseNats        = env.GetBool("USE_NATS", true)
	NatsURL        = env.Get("NATS_URL", "tcp://nats:4222")
	UseRest        = env.GetBool("USE_REST", false)
	RestPort       = env.GetInt("REST_PORT", 8081)
	PortalUsername = env.Get("ADMIN_USER", "admin")
	PortalPassword = env.Get("ADMIN_PASSWORD", "admin")
	PortalPort     = env.GetInt("PORTAL_PORT", 8080)
)
