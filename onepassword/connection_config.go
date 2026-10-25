package onepassword

import (
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

type onepasswordConfig struct {
	Token       *string `hcl:"token"`
	URL         *string `hcl:"url"`
	EventsToken *string `hcl:"events_token"`
	EventsURL   *string `hcl:"events_url"`
}

func ConfigInstance() interface{} {
	return &onepasswordConfig{}
}

// GetConfig :: retrieve and cast connection config from query data
func GetConfig(connection *plugin.Connection) onepasswordConfig {
	if connection == nil || connection.Config == nil {
		return onepasswordConfig{}
	}
	config, _ := connection.Config.(onepasswordConfig)
	return config
}
