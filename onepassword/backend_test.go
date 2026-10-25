package onepassword

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/1Password/connect-sdk-go/connect"
	"github.com/1Password/connect-sdk-go/onepassword"
	opsdk "github.com/1password/onepassword-sdk-go"
)

func envOf(vars map[string]string) func(string) string {
	return func(k string) string { return vars[k] }
}

func TestResolveAuthPrecedence(t *testing.T) {
	cases := []struct {
		name      string
		cfg       onepasswordConfig
		env       map[string]string
		wantMode  authMode
		wantToken string
		wantURL   string
	}{
		{
			name:      "config service account token beats config connect token",
			cfg:       onepasswordConfig{ServiceAccountToken: strPtr("ops_cfg"), Token: strPtr("connect_cfg")},
			wantMode:  authServiceAccount,
			wantToken: "ops_cfg",
		},
		{
			name:      "config connect token beats env service account token",
			cfg:       onepasswordConfig{Token: strPtr("connect_cfg"), URL: strPtr("https://connect.example")},
			env:       map[string]string{"OP_SERVICE_ACCOUNT_TOKEN": "ops_env"},
			wantMode:  authConnect,
			wantToken: "connect_cfg",
			wantURL:   "https://connect.example",
		},
		{
			name:      "env service account token beats env connect token",
			env:       map[string]string{"OP_SERVICE_ACCOUNT_TOKEN": "ops_env", "OP_CONNECT_TOKEN": "connect_env"},
			wantMode:  authServiceAccount,
			wantToken: "ops_env",
		},
		{
			name:      "env connect token with env host",
			env:       map[string]string{"OP_CONNECT_TOKEN": "connect_env", "OP_CONNECT_HOST": "https://env.example"},
			wantMode:  authConnect,
			wantToken: "connect_env",
			wantURL:   "https://env.example",
		},
		{
			name:      "config url beats env host",
			cfg:       onepasswordConfig{URL: strPtr("https://cfg.example")},
			env:       map[string]string{"OP_CONNECT_TOKEN": "connect_env", "OP_CONNECT_HOST": "https://env.example"},
			wantMode:  authConnect,
			wantToken: "connect_env",
			wantURL:   "https://cfg.example",
		},
		{
			name:      "connect url defaults to localhost",
			cfg:       onepasswordConfig{Token: strPtr("connect_cfg")},
			wantMode:  authConnect,
			wantToken: "connect_cfg",
			wantURL:   defaultConnectURL,
		},
		{
			name:      "empty config strings fall through to env",
			cfg:       onepasswordConfig{ServiceAccountToken: strPtr(""), Token: strPtr("")},
			env:       map[string]string{"OP_SERVICE_ACCOUNT_TOKEN": "ops_env"},
			wantMode:  authServiceAccount,
			wantToken: "ops_env",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := resolveAuth(tc.cfg, envOf(tc.env))
			if err != nil {
				t.Fatalf("resolveAuth: %v", err)
			}
			if got.mode != tc.wantMode || got.token != tc.wantToken || got.url != tc.wantURL {
				t.Errorf("got %+v, want mode=%v token=%q url=%q", got, tc.wantMode, tc.wantToken, tc.wantURL)
			}
		})
	}
}

func TestResolveAuthWithoutAnyTokenFails(t *testing.T) {
	_, err := resolveAuth(onepasswordConfig{URL: strPtr("https://connect.example")}, envOf(nil))
	if !errors.Is(err, errMissingCredentials) {
		t.Fatalf("got %v, want errMissingCredentials", err)
	}
}

type fakeVaults struct {
	overviews []opsdk.VaultOverview
}

func (f fakeVaults) List(context.Context, ...opsdk.VaultListParams) ([]opsdk.VaultOverview, error) {
	return f.overviews, nil
}

func (f fakeVaults) GetOverview(_ context.Context, id string) (opsdk.VaultOverview, error) {
	for _, v := range f.overviews {
		if v.ID == id {
			return v, nil
		}
	}
	return opsdk.VaultOverview{}, errors.New("vault not found")
}

type fakeItems struct {
	byVault map[string][]opsdk.ItemOverview
	full    map[string]opsdk.Item
}

func (f fakeItems) List(_ context.Context, vaultID string, _ ...opsdk.ItemListFilter) ([]opsdk.ItemOverview, error) {
	return f.byVault[vaultID], nil
}

func (f fakeItems) Get(_ context.Context, vaultID, itemID string) (opsdk.Item, error) {
	item, ok := f.full[vaultID+"/"+itemID]
	if !ok {
		return opsdk.Item{}, errors.New("item not found")
	}
	return item, nil
}

func TestSDKBackendMapsVaultsToConnectShape(t *testing.T) {
	created := time.Date(2026, 1, 2, 3, 4, 5, 0, time.UTC)
	backend := sdkBackend{vaults: fakeVaults{overviews: []opsdk.VaultOverview{
		{ID: "v1", Title: "Engineering", Description: "shared", VaultType: opsdk.VaultTypeUserCreated, ActiveItemCount: 7, ContentVersion: 12, AttributeVersion: 3, CreatedAt: created, UpdatedAt: created},
		{ID: "v2", Title: "Mystery", VaultType: opsdk.VaultTypeUnsupported},
	}}}

	vaults, err := backend.ListVaults(context.Background())
	if err != nil {
		t.Fatalf("ListVaults: %v", err)
	}
	want := onepassword.Vault{ID: "v1", Name: "Engineering", Description: "shared", Items: 7, ContentVersoin: 12, AttrVersion: 3, Type: onepassword.UserCreatedVault, CreatedAt: created, UpdatedAt: created}
	if vaults[0] != want {
		t.Errorf("got %+v, want %+v", vaults[0], want)
	}
	if vaults[1].Type != onepassword.UnknownVault {
		t.Errorf("unsupported vault type mapped to %q, want UNKNOWN", vaults[1].Type)
	}

	one, err := backend.GetVault(context.Background(), "v1")
	if err != nil || one.Name != "Engineering" {
		t.Errorf("GetVault = %+v, %v", one, err)
	}
}

func TestSDKBackendMapsItemOverviewsToConnectCategories(t *testing.T) {
	backend := sdkBackend{items: fakeItems{byVault: map[string][]opsdk.ItemOverview{
		"v1": {
			{ID: "i1", Title: "GitHub", Category: opsdk.ItemCategoryLogin, VaultID: "v1", Tags: []string{"prod"}, Websites: []opsdk.Website{{URL: "https://github.com", Label: "website"}}},
			{ID: "i2", Title: "Deploy key", Category: opsdk.ItemCategoryAPICredentials, VaultID: "v1"},
			{ID: "i3", Title: "Box", Category: opsdk.ItemCategorySSHKey, VaultID: "v1"},
			{ID: "i4", Title: "Odd", Category: opsdk.ItemCategoryUnsupported, VaultID: "v1"},
		},
	}}}

	items, err := backend.ListItems(context.Background(), "v1")
	if err != nil {
		t.Fatalf("ListItems: %v", err)
	}
	wantCategories := []onepassword.ItemCategory{onepassword.Login, onepassword.ApiCredential, onepassword.SSHKey, onepassword.Custom}
	for i, want := range wantCategories {
		if items[i].Category != want {
			t.Errorf("item %s category = %q, want %q", items[i].ID, items[i].Category, want)
		}
		if items[i].Vault.ID != "v1" {
			t.Errorf("item %s vault = %q, want v1", items[i].ID, items[i].Vault.ID)
		}
	}
	if len(items[0].URLs) != 1 || items[0].URLs[0].URL != "https://github.com" || items[0].Tags[0] != "prod" {
		t.Errorf("urls/tags did not map: %+v", items[0])
	}
	if items[0].Fields != nil {
		t.Errorf("an overview must not carry fields, got %+v", items[0].Fields)
	}
}

func TestSDKBackendGetItemMarksBuiltInLoginFields(t *testing.T) {
	section := "sec1"
	backend := sdkBackend{items: fakeItems{full: map[string]opsdk.Item{
		"v1/i1": {
			ID: "i1", VaultID: "v1", Category: opsdk.ItemCategoryLogin,
			Sections: []opsdk.ItemSection{{ID: "sec1", Title: "Extra"}},
			Fields: []opsdk.ItemField{
				{ID: "username", Title: "username", FieldType: opsdk.ItemFieldTypeText, Value: "alice"},
				{ID: "password", Title: "password", FieldType: opsdk.ItemFieldTypeConcealed, Value: "pw"},
				{ID: "password", Title: "old password", FieldType: opsdk.ItemFieldTypeConcealed, Value: "old", SectionID: &section},
			},
			Notes: "note",
		},
	}}}

	item, err := backend.GetItem(context.Background(), "i1", "v1")
	if err != nil {
		t.Fatalf("GetItem: %v", err)
	}
	if item.Fields[0].Purpose != "USERNAME" || item.Fields[1].Purpose != "PASSWORD" {
		t.Errorf("built-in fields lost their purpose: %+v %+v", item.Fields[0], item.Fields[1])
	}
	if item.Fields[2].Purpose != "" || item.Fields[2].Section == nil || item.Fields[2].Section.Label != "Extra" {
		t.Errorf("a sectioned field must keep its section and no purpose: %+v", item.Fields[2])
	}
	if item.Fields[1].Type != "CONCEALED" {
		t.Errorf("field type = %q, want CONCEALED", item.Fields[1].Type)
	}
	notes := item.Fields[3]
	if notes.ID != "notesPlain" || notes.Purpose != "NOTES" || notes.Value != "note" {
		t.Errorf("notes did not map to the Connect notesPlain field: %+v", notes)
	}
}

func TestSDKBackendRefusesFileOperations(t *testing.T) {
	_, err := sdkBackend{}.ListFiles(context.Background(), "i1", "v1")
	if !errors.Is(err, errFilesNeedConnect) {
		t.Fatalf("got %v, want errFilesNeedConnect", err)
	}
}

func TestConnectBackendCallsConnectServer(t *testing.T) {
	const vaultID = "abcdefghijklmnopqrstuvwxyz"
	var auth []string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		auth = append(auth, r.Header.Get("Authorization"))
		w.Header().Set("Content-Type", "application/json")
		switch r.URL.Path {
		case "/v1/vaults":
			_, _ = w.Write([]byte(`[{"id":"` + vaultID + `","name":"Engineering","type":"USER_CREATED","items":2}]`))
		case "/v1/vaults/" + vaultID + "/items":
			_, _ = w.Write([]byte(`[{"id":"item1","title":"GitHub","category":"LOGIN","vault":{"id":"` + vaultID + `"}}]`))
		default:
			http.NotFound(w, r)
		}
	}))
	defer server.Close()

	backend := connectBackend{client: connect.NewClient(server.URL, "connect-token")}
	vaults, err := backend.ListVaults(context.Background())
	if err != nil || len(vaults) != 1 || vaults[0].Name != "Engineering" {
		t.Fatalf("ListVaults = %+v, %v", vaults, err)
	}
	items, err := backend.ListItems(context.Background(), vaultID)
	if err != nil || len(items) != 1 || items[0].Category != onepassword.Login {
		t.Fatalf("ListItems = %+v, %v", items, err)
	}
	for _, a := range auth {
		if a != "Bearer connect-token" {
			t.Errorf("Authorization = %q, want the Connect token", a)
		}
	}
}
