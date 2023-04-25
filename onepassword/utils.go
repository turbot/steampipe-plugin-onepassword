package onepassword

import (
	"context"
	"strconv"
	"time"

	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func convertTimestamp(_ context.Context, d *transform.TransformData) (interface{}, error) {
	epochTime := d.Value.(string)

	if epochTime != "" {
		unixtime, _ := strconv.Atoi(epochTime)
		unixTimestamp := time.Unix(int64(unixtime), 0)
		return unixTimestamp, nil
	}
	return nil, nil
}
