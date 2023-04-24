package main

import (
	"github.com/turbot/steampipe-plugin-1password/onepassword"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{PluginFunc: onepassword.Plugin})
}
