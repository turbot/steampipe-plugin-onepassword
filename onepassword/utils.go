package onepassword

import (
	"context"
	"math"
	"time"

	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func convertTimestamp(_ context.Context, d *transform.TransformData) (interface{}, error) {
	epochTime := d.Value.(*int)

	if epochTime != nil {
		timeInSec := math.Floor(float64(*epochTime) / 1000)
		unixTimestamp := time.Unix(int64(timeInSec), 0)
		timestampRFC3339Format := unixTimestamp.Format(time.RFC3339)
		return timestampRFC3339Format, nil
	}
	return nil, nil
}
