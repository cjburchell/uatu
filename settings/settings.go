package settings

import (
	"github.com/cjburchell/pubsub"
	pubSubSettings "github.com/cjburchell/pubsub/settings"
	"github.com/cjburchell/settings-go"
)

//Get the application configuration
func Get(config settings.ISettings) AppConfig {
	return AppConfig{
		UsePubSub:  config.GetSection("PubSub").GetBool("Enabled", true),
		PubSub:     pubSubSettings.Get(config.GetSection("PubSub")),
		UseRest:    config.GetSection("Rest").GetBool("Enabled", false),
		RestPort:   config.GetSection("Rest").GetInt("Port", 8081),
		RestToken:  config.GetSection("Rest").Get("Token", "token"),
	}
}

// AppConfig object
type AppConfig struct {
	UsePubSub  bool
	PubSub     pubsub.Settings
	UseRest    bool
	RestPort   int
	RestToken  string
}
