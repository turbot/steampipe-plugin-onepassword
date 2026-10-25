package onepassword

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

const (
	defaultEventsURL    = "https://events.1password.com"
	eventsPageLimit     = 1000
	defaultEventsWindow = 24 * time.Hour

	signInAttemptsPath = "/api/v2/signinattempts"
	auditEventsPath    = "/api/v2/auditevents"
	itemUsagesPath     = "/api/v2/itemusages"
)

var errMissingEventsToken = errors.New("'events_token' must be set in the connection configuration to query 1Password Events tables. Edit your connection configuration file and then restart Steampipe")

type eventsClient struct {
	apiTransport
	baseURL string
	token   string
}

func newEventsClient(cfg onepasswordConfig) (*eventsClient, error) {
	if !nonEmpty(cfg.EventsToken) {
		return nil, errMissingEventsToken
	}
	baseURL := defaultEventsURL
	if nonEmpty(cfg.EventsURL) {
		baseURL = *cfg.EventsURL
	}
	return &eventsClient{apiTransport: newAPITransport(), baseURL: strings.TrimRight(baseURL, "/"), token: *cfg.EventsToken}, nil
}

func getEventsClient(d *plugin.QueryData) (*eventsClient, error) {
	return newEventsClient(GetConfig(d.Connection))
}

type eventsPage struct {
	Cursor  string            `json:"cursor"`
	HasMore bool              `json:"has_more"`
	Items   []json.RawMessage `json:"items"`
}

type timeWindow struct {
	start time.Time
	end   time.Time
}

type timeQual struct {
	operator string
	value    time.Time
}

// eventWindow turns timestamp quals into the API's start_time/end_time. Both bounds are sent
// inclusive; Postgres re-applies the strict comparisons to the returned rows.
func eventWindow(quals []timeQual, now time.Time) timeWindow {
	var w timeWindow
	for _, q := range quals {
		switch q.operator {
		case ">", ">=":
			if w.start.IsZero() || q.value.After(w.start) {
				w.start = q.value
			}
		case "<", "<=":
			if w.end.IsZero() || q.value.Before(w.end) {
				w.end = q.value
			}
		}
	}
	if w.start.IsZero() {
		w.start = now.Add(-defaultEventsWindow)
		if !w.end.IsZero() && w.end.Before(now) {
			w.start = w.end.Add(-defaultEventsWindow)
		}
	}
	return w
}

func timestampQuals(d *plugin.QueryData) []timeQual {
	var quals []timeQual
	if d.Quals["timestamp"] == nil {
		return nil
	}
	for _, q := range d.Quals["timestamp"].Quals {
		if ts := q.Value.GetTimestampValue(); ts != nil {
			quals = append(quals, timeQual{operator: q.Operator, value: ts.AsTime()})
		}
	}
	return quals
}

func (c *eventsClient) post(ctx context.Context, path string, body map[string]any) (*eventsPage, error) {
	payload, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	resp, err := c.do(ctx, func() (*http.Request, error) {
		req, err := http.NewRequest(http.MethodPost, c.baseURL+path, bytes.NewReader(payload))
		if err != nil {
			return nil, err
		}
		req.Header.Set("Authorization", "Bearer "+c.token)
		req.Header.Set("Content-Type", "application/json")
		return req, nil
	})
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()
	if resp.StatusCode != http.StatusOK {
		return nil, newAPIError("Events", resp)
	}
	page := new(eventsPage)
	if err := json.NewDecoder(resp.Body).Decode(page); err != nil {
		return nil, err
	}
	return page, nil
}

// listEvents pages through one Events endpoint, stopping early when emit returns false.
func listEvents[T any](ctx context.Context, c *eventsClient, path string, w timeWindow, emit func(T) bool) error {
	body := map[string]any{"limit": eventsPageLimit, "start_time": w.start.UTC().Format(time.RFC3339)}
	if !w.end.IsZero() {
		body["end_time"] = w.end.UTC().Format(time.RFC3339)
	}
	for {
		page, err := c.post(ctx, path, body)
		if err != nil {
			return err
		}
		for _, raw := range page.Items {
			var event T
			if err := json.Unmarshal(raw, &event); err != nil {
				return err
			}
			if !emit(event) {
				return nil
			}
		}
		if !page.HasMore || page.Cursor == "" {
			return nil
		}
		body = map[string]any{"cursor": page.Cursor}
	}
}

// streamEvents is the shared List hydrate body of the three Events tables.
func streamEvents[T any](ctx context.Context, d *plugin.QueryData, table, path string) (interface{}, error) {
	client, err := getEventsClient(d)
	if err != nil {
		plugin.Logger(ctx).Error(table, "connection_error", err)
		return nil, err
	}
	window := eventWindow(timestampQuals(d), time.Now())
	err = listEvents(ctx, client, path, window, func(event T) bool {
		d.StreamListItem(ctx, event)
		return d.RowsRemaining(ctx) != 0
	})
	if err != nil {
		plugin.Logger(ctx).Error(table, "api_error", err)
		return nil, err
	}
	return nil, nil
}

type eventUser struct {
	UUID            string `json:"uuid"`
	Name            string `json:"name"`
	Email           string `json:"email"`
	UserType        string `json:"user_type"`
	UserAccountUUID string `json:"user_account_uuid"`
}

type eventClientInfo struct {
	AppName         string `json:"app_name"`
	AppVersion      string `json:"app_version"`
	PlatformName    string `json:"platform_name"`
	PlatformVersion string `json:"platform_version"`
	OSName          string `json:"os_name"`
	OSVersion       string `json:"os_version"`
	IPAddress       string `json:"ip_address"`
}

type eventLocation struct {
	Country   string  `json:"country"`
	Region    string  `json:"region"`
	City      string  `json:"city"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

type eventSession struct {
	UUID       string     `json:"uuid"`
	LoginTime  *time.Time `json:"login_time"`
	DeviceUUID string     `json:"device_uuid"`
	IP         string     `json:"ip"`
}

type signInAttempt struct {
	UUID        string           `json:"uuid"`
	SessionUUID string           `json:"session_uuid"`
	Timestamp   time.Time        `json:"timestamp"`
	Category    string           `json:"category"`
	Type        string           `json:"type"`
	Country     string           `json:"country"`
	Details     *json.RawMessage `json:"details"`
	TargetUser  eventUser        `json:"target_user"`
	Client      *eventClientInfo `json:"client"`
	Location    *eventLocation   `json:"location"`
	AccountUUID string           `json:"account_uuid"`
}

type auditEvent struct {
	UUID             string           `json:"uuid"`
	Timestamp        time.Time        `json:"timestamp"`
	ActorUUID        string           `json:"actor_uuid"`
	ActorDetails     eventUser        `json:"actor_details"`
	ActorType        string           `json:"actor_type"`
	ActorAccountUUID string           `json:"actor_account_uuid"`
	AccountUUID      string           `json:"account_uuid"`
	Action           string           `json:"action"`
	ObjectType       string           `json:"object_type"`
	ObjectUUID       string           `json:"object_uuid"`
	ObjectDetails    *json.RawMessage `json:"object_details"`
	AuxID            *int64           `json:"aux_id"`
	AuxUUID          string           `json:"aux_uuid"`
	AuxDetails       *json.RawMessage `json:"aux_details"`
	AuxInfo          string           `json:"aux_info"`
	Session          *eventSession    `json:"session"`
	Location         *eventLocation   `json:"location"`
}

type itemUsage struct {
	UUID        string           `json:"uuid"`
	Timestamp   time.Time        `json:"timestamp"`
	UsedVersion *int64           `json:"used_version"`
	VaultUUID   string           `json:"vault_uuid"`
	ItemUUID    string           `json:"item_uuid"`
	Action      string           `json:"action"`
	User        eventUser        `json:"user"`
	Client      *eventClientInfo `json:"client"`
	Location    *eventLocation   `json:"location"`
	AccountUUID string           `json:"account_uuid"`
}
