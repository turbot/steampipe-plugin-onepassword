package onepassword

import (
	"context"
	"errors"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/1Password/connect-sdk-go/connect"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func getClient(ctx context.Context, d *plugin.QueryData) (connect.Client, error) {
	onepasswordConfig := GetConfig(d.Connection)

	token := os.Getenv("OP_CONNECT_TOKEN")
	url := os.Getenv("OP_CONNECT_HOST")

	if onepasswordConfig.Token != nil {
		token = *onepasswordConfig.Token
	}
	if onepasswordConfig.URL != nil {
		url = *onepasswordConfig.URL
	}

	if url != "" && token != "" {
		client := connect.NewClient(url, token)
		return client, nil
	} else if url == "" && token != "" { // set default url
		client := connect.NewClient("http://localhost:8080", token)
		return client, nil
	}

	return nil, errors.New("'token' must be set in the connection configuration. Edit your connection configuration file and then restart Steampipe.")
}

func convertTimestamp(_ context.Context, d *transform.TransformData) (interface{}, error) {
	epochTime := d.Value.(string)

	if epochTime != "" {
		unixtime, _ := strconv.Atoi(epochTime)
		unixTimestamp := time.Unix(int64(unixtime), 0)
		return unixTimestamp, nil
	}
	return nil, nil
}

func isNotFoundError(err error) bool {
	return strings.Contains(err.Error(), "not found")
}
