package onepassword

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"
)

type usersFakeClock struct {
	mu  sync.Mutex
	cur time.Time
}

func (c *usersFakeClock) now() time.Time {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.cur
}

func (c *usersFakeClock) advance(d time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.cur = c.cur.Add(d)
}

func testUsersClient(baseURL string, clock *usersFakeClock, waits *[]time.Duration) *usersClient {
	return &usersClient{
		apiTransport: apiTransport{
			httpClient: http.DefaultClient,
			sleep:      func(_ context.Context, d time.Duration) error { *waits = append(*waits, d); return nil },
			now:        clock.now,
		},
		baseURL:      baseURL,
		accountID:    "ACCT1",
		clientID:     "client-id",
		clientSecret: "client-secret",
	}
}

func TestNewUsersClientRequiresAllThreeKeys(t *testing.T) {
	full := onepasswordConfig{UsersClientID: strPtr("id"), UsersClientSecret: strPtr("secret"), UsersAccountID: strPtr("acct")}
	if _, err := newUsersClient(full); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	cases := map[string]onepasswordConfig{
		"none":       {},
		"no id":      {UsersClientSecret: strPtr("secret"), UsersAccountID: strPtr("acct")},
		"no secret":  {UsersClientID: strPtr("id"), UsersAccountID: strPtr("acct")},
		"no account": {UsersClientID: strPtr("id"), UsersClientSecret: strPtr("secret")},
		"empty":      {UsersClientID: strPtr(""), UsersClientSecret: strPtr("secret"), UsersAccountID: strPtr("acct")},
	}
	for name, cfg := range cases {
		t.Run(name, func(t *testing.T) {
			_, err := newUsersClient(cfg)
			if err == nil || !strings.Contains(err.Error(), "users_client_id, users_client_secret and users_account_id must be set") {
				t.Errorf("got %v", err)
			}
		})
	}
}

func TestNewUsersClientDefaultsAndTrimsTheURL(t *testing.T) {
	cfg := onepasswordConfig{UsersClientID: strPtr("id-default"), UsersClientSecret: strPtr("s"), UsersAccountID: strPtr("a")}
	c, _ := newUsersClient(cfg)
	if c.baseURL != defaultUsersAPIURL {
		t.Errorf("got %q", c.baseURL)
	}
	cfg.UsersClientID = strPtr("id-custom")
	cfg.UsersAPIURL = strPtr("https://api.example.test/")
	c, _ = newUsersClient(cfg)
	if c.baseURL != "https://api.example.test" {
		t.Errorf("got %q", c.baseURL)
	}
}

func TestListUsersAuthenticatesAndFollowsPageTokens(t *testing.T) {
	var tokenForm map[string]string
	var tokenCalls int
	var userQueries []map[string]string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/v1/oauth/token":
			tokenCalls++
			_ = r.ParseForm()
			tokenForm = map[string]string{
				"grant_type":    r.PostForm.Get("grant_type"),
				"client_id":     r.PostForm.Get("client_id"),
				"client_secret": r.PostForm.Get("client_secret"),
				"content_type":  r.Header.Get("Content-Type"),
				"user_agent":    r.Header.Get("User-Agent"),
			}
			_, _ = w.Write([]byte(`{"access_token":"tok-1","token_type":"Bearer","expires_in":900}`))
		case "/v1/accounts/ACCT1/users":
			if got := r.Header.Get("Authorization"); got != "Bearer tok-1" {
				t.Errorf("authorization = %q", got)
			}
			if r.Header.Get("User-Agent") == "" {
				t.Error("missing User-Agent")
			}
			userQueries = append(userQueries, map[string]string{
				"max_page_size": r.URL.Query().Get("max_page_size"),
				"page_token":    r.URL.Query().Get("page_token"),
			})
			if r.URL.Query().Get("page_token") == "" {
				_, _ = w.Write([]byte(`{"results":[{"id":"u1","email":"a@example.com","display_name":"A","state":"ACTIVE","create_time":"2026-01-02T03:04:05Z"}],"next_page_token":"p2"}`))
				return
			}
			_, _ = w.Write([]byte(`{"results":[{"id":"u2","email":"b@example.com","display_name":"B","state":"SUSPENDED"}]}`))
		default:
			t.Errorf("unexpected path %s", r.URL.Path)
		}
	}))
	defer server.Close()

	var waits []time.Duration
	client := testUsersClient(server.URL, &usersFakeClock{cur: time.Unix(1_800_000_000, 0)}, &waits)
	var got []apiUser
	if err := client.listUsers(context.Background(), func(u apiUser) bool { got = append(got, u); return true }); err != nil {
		t.Fatal(err)
	}

	if len(got) != 2 || got[0].ID != "u1" || got[1].State != "SUSPENDED" {
		t.Fatalf("users = %+v", got)
	}
	if got[0].CreateTime == nil || got[1].CreateTime != nil {
		t.Errorf("create_time = %v, %v", got[0].CreateTime, got[1].CreateTime)
	}
	if tokenCalls != 1 {
		t.Errorf("token requested %d times, want 1", tokenCalls)
	}
	want := map[string]string{"grant_type": "client_credentials", "client_id": "client-id", "client_secret": "client-secret", "content_type": "application/x-www-form-urlencoded", "user_agent": usersUserAgent}
	for k, v := range want {
		if tokenForm[k] != v {
			t.Errorf("token request %s = %q, want %q", k, tokenForm[k], v)
		}
	}
	if len(userQueries) != 2 || userQueries[0]["max_page_size"] != "1000" || userQueries[1]["page_token"] != "p2" {
		t.Errorf("queries = %v", userQueries)
	}
}

func TestListUsersStopsWhenEmitReturnsFalse(t *testing.T) {
	pages := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == usersTokenPath {
			_, _ = w.Write([]byte(`{"access_token":"t","expires_in":900}`))
			return
		}
		pages++
		_, _ = w.Write([]byte(`{"results":[{"id":"u1"},{"id":"u2"}],"next_page_token":"more"}`))
	}))
	defer server.Close()

	var waits []time.Duration
	client := testUsersClient(server.URL, &usersFakeClock{cur: time.Now()}, &waits)
	count := 0
	_ = client.listUsers(context.Background(), func(apiUser) bool { count++; return false })
	if count != 1 || pages != 1 {
		t.Errorf("emitted %d users over %d pages", count, pages)
	}
}

func TestAccessTokenIsRequestedAgainBeforeItExpires(t *testing.T) {
	tokenCalls := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenCalls++
		_ = json.NewEncoder(w).Encode(map[string]any{"access_token": "tok-" + string(rune('0'+tokenCalls)), "expires_in": 900})
	}))
	defer server.Close()

	var waits []time.Duration
	clock := &usersFakeClock{cur: time.Unix(1_800_000_000, 0)}
	client := testUsersClient(server.URL, clock, &waits)
	ctx := context.Background()

	first, _ := client.accessToken(ctx)
	clock.advance(800 * time.Second)
	again, _ := client.accessToken(ctx)
	if first != again || tokenCalls != 1 {
		t.Fatalf("token should be reused while fresh: %q %q, %d calls", first, again, tokenCalls)
	}
	clock.advance(50 * time.Second) // 850s: inside the 60s refresh margin of the 900s lifetime
	renewed, _ := client.accessToken(ctx)
	if renewed == first || tokenCalls != 2 {
		t.Errorf("token should be renewed: %q, %d calls", renewed, tokenCalls)
	}
}

func TestUsersRequestsWaitOutRateLimits(t *testing.T) {
	limited := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == usersTokenPath {
			_, _ = w.Write([]byte(`{"access_token":"t","expires_in":900}`))
			return
		}
		if limited == 0 {
			limited++
			w.Header().Set("Retry-After", "7")
			w.WriteHeader(http.StatusTooManyRequests)
			return
		}
		_, _ = w.Write([]byte(`{"results":[{"id":"u1"}]}`))
	}))
	defer server.Close()

	var waits []time.Duration
	client := testUsersClient(server.URL, &usersFakeClock{cur: time.Now()}, &waits)
	count := 0
	if err := client.listUsers(context.Background(), func(apiUser) bool { count++; return true }); err != nil {
		t.Fatal(err)
	}
	if count != 1 || len(waits) != 1 || waits[0] != 7*time.Second {
		t.Errorf("count=%d waits=%v", count, waits)
	}
}

func TestUsersAPIErrorsCarryTheStatusAndMessage(t *testing.T) {
	cases := []struct {
		name      string
		tokenCode int
		tokenBody string
		userCode  int
		wantParts []string
	}{
		{"bad client", http.StatusUnauthorized, `{"error":"invalid_client","error_description":"client_id or client_secret is invalid"}`, 0, []string{"HTTP 401", "client_id or client_secret is invalid"}},
		{"forbidden", http.StatusOK, `{"access_token":"t","expires_in":900}`, http.StatusForbidden, []string{"HTTP 403"}},
		{"unknown account", http.StatusOK, `{"access_token":"t","expires_in":900}`, http.StatusNotFound, []string{"HTTP 404"}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Path == usersTokenPath {
					w.WriteHeader(tc.tokenCode)
					_, _ = w.Write([]byte(tc.tokenBody))
					return
				}
				w.WriteHeader(tc.userCode)
			}))
			defer server.Close()

			var waits []time.Duration
			client := testUsersClient(server.URL, &usersFakeClock{cur: time.Now()}, &waits)
			err := client.listUsers(context.Background(), func(apiUser) bool { return true })
			if err == nil {
				t.Fatal("expected an error")
			}
			for _, part := range tc.wantParts {
				if !strings.Contains(err.Error(), part) {
					t.Errorf("error %q lacks %q", err, part)
				}
			}
		})
	}
}
