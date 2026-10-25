package onepassword

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"
)

type recordedSleeps struct {
	mu    sync.Mutex
	waits []time.Duration
}

func (r *recordedSleeps) sleep(_ context.Context, d time.Duration) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.waits = append(r.waits, d)
	return nil
}

func testEventsClient(baseURL string, sleeps *recordedSleeps) *eventsClient {
	return &eventsClient{
		apiTransport: apiTransport{httpClient: http.DefaultClient, sleep: sleeps.sleep, now: time.Now},
		baseURL:      baseURL,
		token:        "events-token",
	}
}

func TestEventWindow(t *testing.T) {
	now := time.Date(2026, 9, 30, 12, 0, 0, 0, time.UTC)
	day := func(d int) time.Time { return time.Date(2026, 9, d, 0, 0, 0, 0, time.UTC) }

	cases := []struct {
		name      string
		quals     []timeQual
		wantStart time.Time
		wantEnd   time.Time
	}{
		{"no quals defaults to the last 24 hours", nil, now.Add(-24 * time.Hour), time.Time{}},
		{"lower bound only", []timeQual{{">=", day(23)}}, day(23), time.Time{}},
		{"strict bounds map to the same instants", []timeQual{{">", day(23)}, {"<", day(25)}}, day(23), day(25)},
		{"tightest lower and upper bounds win", []timeQual{{">=", day(20)}, {">", day(22)}, {"<=", day(28)}, {"<", day(26)}}, day(22), day(26)},
		{"an upper bound in the past without a lower bound looks back 24 hours from it", []timeQual{{"<", day(10)}}, day(9), day(10)},
		{"an upper bound in the future keeps the default start", []timeQual{{"<", day(31)}}, now.Add(-24 * time.Hour), day(31)},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := eventWindow(tc.quals, now)
			if !got.start.Equal(tc.wantStart) || !got.end.Equal(tc.wantEnd) {
				t.Errorf("got start=%v end=%v, want start=%v end=%v", got.start, got.end, tc.wantStart, tc.wantEnd)
			}
		})
	}
}

func TestListEventsFollowsTheCursorAcrossPages(t *testing.T) {
	var bodies []map[string]any
	var auth []string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost || r.URL.Path != signInAttemptsPath {
			t.Errorf("unexpected %s %s", r.Method, r.URL.Path)
		}
		auth = append(auth, r.Header.Get("Authorization"))
		var body map[string]any
		_ = json.NewDecoder(r.Body).Decode(&body)
		bodies = append(bodies, body)
		if _, ok := body["cursor"]; !ok {
			_, _ = w.Write([]byte(`{"cursor":"c1","has_more":true,"items":[
				{"uuid":"a1","timestamp":"2026-09-29T10:00:00Z","category":"credentials_failed","type":"password_secret_bad",
				 "target_user":{"uuid":"u1","name":"Alice","email":"alice@example.com","user_type":"user"},
				 "client":{"ip_address":"203.0.113.9"}}]}`))
			return
		}
		_, _ = w.Write([]byte(`{"cursor":"c2","has_more":false,"items":[{"uuid":"a2","timestamp":"2026-09-29T11:00:00Z","category":"success"}]}`))
	}))
	defer server.Close()

	window := timeWindow{start: time.Date(2026, 9, 29, 0, 0, 0, 0, time.UTC), end: time.Date(2026, 9, 30, 0, 0, 0, 0, time.UTC)}
	var got []signInAttempt
	err := listEvents(context.Background(), testEventsClient(server.URL, &recordedSleeps{}), signInAttemptsPath, window, func(a signInAttempt) bool {
		got = append(got, a)
		return true
	})
	if err != nil {
		t.Fatalf("listEvents: %v", err)
	}

	if len(got) != 2 || got[0].UUID != "a1" || got[1].UUID != "a2" {
		t.Fatalf("got %+v, want a1 then a2", got)
	}
	if got[0].TargetUser.Email != "alice@example.com" || got[0].Client.IPAddress != "203.0.113.9" {
		t.Errorf("nested fields did not decode: %+v", got[0])
	}
	if got[1].Client != nil {
		t.Errorf("an absent client must stay nil so the column is null, got %+v", got[1].Client)
	}

	first := bodies[0]
	if first["limit"] != float64(eventsPageLimit) || first["start_time"] != "2026-09-29T00:00:00Z" || first["end_time"] != "2026-09-30T00:00:00Z" {
		t.Errorf("first request must reset the cursor with the window, got %v", first)
	}
	second := bodies[1]
	if len(second) != 1 || second["cursor"] != "c1" {
		t.Errorf("follow-up request must send only the cursor, got %v", second)
	}
	for _, a := range auth {
		if a != "Bearer events-token" {
			t.Errorf("Authorization = %q", a)
		}
	}
}

func TestListEventsStopsWhenTheCallerIsDone(t *testing.T) {
	calls := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		calls++
		_, _ = w.Write([]byte(`{"cursor":"more","has_more":true,"items":[{"uuid":"x"},{"uuid":"y"}]}`))
	}))
	defer server.Close()

	seen := 0
	err := listEvents(context.Background(), testEventsClient(server.URL, &recordedSleeps{}), auditEventsPath, timeWindow{start: time.Now()}, func(auditEvent) bool {
		seen++
		return false
	})
	if err != nil || seen != 1 || calls != 1 {
		t.Fatalf("err=%v seen=%d calls=%d, want one row from one request", err, seen, calls)
	}
}

func TestListEventsWaitsOutRateLimitsUsingRetryAfter(t *testing.T) {
	var bodies []string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var body map[string]any
		_ = json.NewDecoder(r.Body).Decode(&body)
		encoded, _ := json.Marshal(body)
		bodies = append(bodies, string(encoded))
		if len(bodies) == 1 {
			w.Header().Set("Retry-After", "7")
			w.WriteHeader(http.StatusTooManyRequests)
			return
		}
		_, _ = w.Write([]byte(`{"has_more":false,"items":[{"uuid":"u","action":"disblmfa"}]}`))
	}))
	defer server.Close()

	sleeps := &recordedSleeps{}
	var got []auditEvent
	err := listEvents(context.Background(), testEventsClient(server.URL, sleeps), auditEventsPath, timeWindow{start: time.Now()}, func(e auditEvent) bool {
		got = append(got, e)
		return true
	})
	if err != nil {
		t.Fatalf("listEvents: %v", err)
	}
	if len(sleeps.waits) != 1 || sleeps.waits[0] != 7*time.Second {
		t.Errorf("waits = %v, want one 7s wait", sleeps.waits)
	}
	if len(bodies) != 2 || bodies[0] != bodies[1] {
		t.Errorf("the retry must resend the same body, got %v", bodies)
	}
	if len(got) != 1 || got[0].Action != "disblmfa" {
		t.Errorf("got %+v", got)
	}
}

func TestListEventsGivesUpAfterRepeatedRateLimits(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusTooManyRequests)
	}))
	defer server.Close()

	sleeps := &recordedSleeps{}
	err := listEvents(context.Background(), testEventsClient(server.URL, sleeps), itemUsagesPath, timeWindow{start: time.Now()}, func(itemUsage) bool { return true })
	var apiErr *apiError
	if !errors.As(err, &apiErr) || apiErr.StatusCode != http.StatusTooManyRequests {
		t.Fatalf("got %v, want a 429 apiError", err)
	}
	if len(sleeps.waits) != maxRateLimitRetries {
		t.Errorf("waited %d times, want %d", len(sleeps.waits), maxRateLimitRetries)
	}
}

func TestListEventsReportsTheAPIErrorMessage(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		_, _ = w.Write([]byte(`{"Error":{"Message":"Unauthorized"}}`))
	}))
	defer server.Close()

	err := listEvents(context.Background(), testEventsClient(server.URL, &recordedSleeps{}), signInAttemptsPath, timeWindow{start: time.Now()}, func(signInAttempt) bool { return true })
	var apiErr *apiError
	if !errors.As(err, &apiErr) || apiErr.StatusCode != http.StatusUnauthorized || apiErr.Message != "Unauthorized" || apiErr.API != "Events" {
		t.Fatalf("got %#v", err)
	}
}

func TestRetryAfter(t *testing.T) {
	now := time.Date(2026, 9, 30, 12, 0, 0, 0, time.UTC)
	cases := []struct {
		name    string
		header  string
		attempt int
		want    time.Duration
	}{
		{"seconds", "3", 0, 3 * time.Second},
		{"http date", now.Add(9 * time.Second).Format(http.TimeFormat), 0, 9 * time.Second},
		{"date in the past", now.Add(-time.Minute).Format(http.TimeFormat), 0, 0},
		{"missing header backs off exponentially", "", 3, 8 * time.Second},
		{"capped", "3600", 0, maxRetryWait},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := retryAfter(tc.header, now, tc.attempt); got != tc.want {
				t.Errorf("got %v, want %v", got, tc.want)
			}
		})
	}
}

func TestNewEventsClientConfig(t *testing.T) {
	if _, err := newEventsClient(onepasswordConfig{}); !errors.Is(err, errMissingEventsToken) {
		t.Errorf("no token: got %v, want errMissingEventsToken", err)
	}
	c, err := newEventsClient(onepasswordConfig{EventsToken: strPtr("t")})
	if err != nil || c.baseURL != defaultEventsURL {
		t.Errorf("default url: got %+v, %v", c, err)
	}
	c, err = newEventsClient(onepasswordConfig{EventsToken: strPtr("t"), EventsURL: strPtr("https://events.1password.eu/")})
	if err != nil || c.baseURL != "https://events.1password.eu" {
		t.Errorf("custom url: got %+v, %v", c, err)
	}
}
