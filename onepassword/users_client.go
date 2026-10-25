package onepassword

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

const (
	defaultUsersAPIURL = "https://api.1password.com"
	usersPageSize      = 1000
	tokenRefreshMargin = 60 * time.Second
	usersUserAgent     = "steampipe-plugin-onepassword"

	usersTokenPath = "/v1/oauth/token"
)

var errMissingUsersConfig = errors.New("users_client_id, users_client_secret and users_account_id must be set in the connection configuration to query onepassword_user. Edit your connection configuration file and then restart Steampipe")

// usersTokenCache shares access tokens across queries so each one does not spend a token request
// against the Users API's 100 requests/min limit.
var usersTokenCache sync.Map

type usersClient struct {
	apiTransport
	baseURL      string
	accountID    string
	clientID     string
	clientSecret string

	mu      sync.Mutex
	token   string
	expires time.Time
}

func newUsersClient(cfg onepasswordConfig) (*usersClient, error) {
	if !nonEmpty(cfg.UsersClientID) || !nonEmpty(cfg.UsersClientSecret) || !nonEmpty(cfg.UsersAccountID) {
		return nil, errMissingUsersConfig
	}
	baseURL := defaultUsersAPIURL
	if nonEmpty(cfg.UsersAPIURL) {
		baseURL = *cfg.UsersAPIURL
	}
	baseURL = strings.TrimRight(baseURL, "/")

	key := strings.Join([]string{baseURL, *cfg.UsersClientID, *cfg.UsersClientSecret}, "\x00")
	if cached, ok := usersTokenCache.Load(key); ok {
		client := cached.(*usersClient)
		if client.accountID == *cfg.UsersAccountID {
			return client, nil
		}
	}
	client := &usersClient{
		apiTransport: newAPITransport(),
		baseURL:      baseURL,
		accountID:    *cfg.UsersAccountID,
		clientID:     *cfg.UsersClientID,
		clientSecret: *cfg.UsersClientSecret,
	}
	usersTokenCache.Store(key, client)
	return client, nil
}

func getUsersClient(d *plugin.QueryData) (*usersClient, error) {
	return newUsersClient(GetConfig(d.Connection))
}

// accessToken returns the cached token, requesting a new one shortly before it expires.
func (c *usersClient) accessToken(ctx context.Context) (string, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.token != "" && c.now().Before(c.expires.Add(-tokenRefreshMargin)) {
		return c.token, nil
	}

	form := url.Values{
		"grant_type":    {"client_credentials"},
		"client_id":     {c.clientID},
		"client_secret": {c.clientSecret},
	}.Encode()
	resp, err := c.do(ctx, func() (*http.Request, error) {
		req, err := http.NewRequest(http.MethodPost, c.baseURL+usersTokenPath, strings.NewReader(form))
		if err != nil {
			return nil, err
		}
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		req.Header.Set("User-Agent", usersUserAgent)
		return req, nil
	})
	if err != nil {
		return "", err
	}
	defer func() { _ = resp.Body.Close() }()
	if resp.StatusCode != http.StatusOK {
		return "", newAPIError("Users", resp)
	}
	var parsed struct {
		AccessToken string `json:"access_token"`
		ExpiresIn   int    `json:"expires_in"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&parsed); err != nil {
		return "", err
	}
	if parsed.AccessToken == "" {
		return "", errors.New("1Password Users API token response contained no access_token")
	}
	c.token = parsed.AccessToken
	c.expires = c.now().Add(time.Duration(parsed.ExpiresIn) * time.Second)
	return c.token, nil
}

type usersPage struct {
	Results       []apiUser `json:"results"`
	NextPageToken string    `json:"next_page_token"`
}

type apiUser struct {
	ID          string     `json:"id"`
	Email       string     `json:"email"`
	DisplayName string     `json:"display_name"`
	State       string     `json:"state"`
	CreateTime  *time.Time `json:"create_time"`
	Path        string     `json:"path"`
}

func (c *usersClient) listPage(ctx context.Context, pageToken string) (*usersPage, error) {
	token, err := c.accessToken(ctx)
	if err != nil {
		return nil, err
	}
	query := url.Values{"max_page_size": {strconv.Itoa(usersPageSize)}}
	if pageToken != "" {
		query.Set("page_token", pageToken)
	}
	endpoint := c.baseURL + "/v1/accounts/" + url.PathEscape(c.accountID) + "/users?" + query.Encode()

	resp, err := c.do(ctx, func() (*http.Request, error) {
		req, err := http.NewRequest(http.MethodGet, endpoint, nil)
		if err != nil {
			return nil, err
		}
		req.Header.Set("Authorization", "Bearer "+token)
		req.Header.Set("User-Agent", usersUserAgent)
		return req, nil
	})
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()
	if resp.StatusCode != http.StatusOK {
		return nil, newAPIError("Users", resp)
	}
	page := new(usersPage)
	if err := json.NewDecoder(resp.Body).Decode(page); err != nil {
		return nil, err
	}
	return page, nil
}

// listUsers pages through the account's users, stopping early when emit returns false.
func (c *usersClient) listUsers(ctx context.Context, emit func(apiUser) bool) error {
	pageToken := ""
	for {
		page, err := c.listPage(ctx, pageToken)
		if err != nil {
			return err
		}
		for _, user := range page.Results {
			if !emit(user) {
				return nil
			}
		}
		if page.NextPageToken == "" {
			return nil
		}
		pageToken = page.NextPageToken
	}
}
