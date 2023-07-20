package onepassword

import (
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/schema"
)

type onepasswordConfig struct {
	Token *string `cty:"token"`
	URL   *string `cty:"url"`
}

var ConfigSchema = map[string]*schema.Attribute{
	"token": {
		Type: schema.TypeString,
	},
	"url": {
		Type: schema.TypeString,
	},
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
