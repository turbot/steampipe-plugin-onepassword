package onepassword

import (
	"context"
	"errors"
	"os"
	"strings"
	"sync"

	"github.com/1Password/connect-sdk-go/connect"
	"github.com/1Password/connect-sdk-go/onepassword"
	opsdk "github.com/1password/onepassword-sdk-go"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

// vaultBackend is the vault and item source behind the tables, either a Connect server or the
// service-account SDK. Rows keep the Connect types so the table definitions work with both.
type vaultBackend interface {
	ListVaults(ctx context.Context) ([]onepassword.Vault, error)
	GetVault(ctx context.Context, id string) (*onepassword.Vault, error)
	ListItems(ctx context.Context, vaultID string) ([]onepassword.Item, error)
	GetItem(ctx context.Context, itemID, vaultID string) (*onepassword.Item, error)
	ListFiles(ctx context.Context, itemID, vaultID string) ([]onepassword.File, error)
	GetFile(ctx context.Context, fileID, itemID, vaultID string) (*onepassword.File, error)
	GetFileContent(ctx context.Context, file *onepassword.File) ([]byte, error)
}

const defaultConnectURL = "http://localhost:8080"

var errMissingCredentials = errors.New("either 'service_account_token' or 'token' must be set in the connection configuration. Edit your connection configuration file and then restart Steampipe.")

var errFilesNeedConnect = errors.New("onepassword_item_file is only supported with a Connect server ('token' and 'url'), not a service account")

type authMode int

const (
	authConnect authMode = iota
	authServiceAccount
)

type authSettings struct {
	mode  authMode
	token string
	url   string
}

// resolveAuth picks the credential: connection config beats environment variables, and within
// each source a service account token beats a Connect token.
func resolveAuth(cfg onepasswordConfig, getenv func(string) string) (authSettings, error) {
	switch {
	case nonEmpty(cfg.ServiceAccountToken):
		return authSettings{mode: authServiceAccount, token: *cfg.ServiceAccountToken}, nil
	case nonEmpty(cfg.Token):
		return connectSettings(*cfg.Token, cfg, getenv), nil
	case getenv("OP_SERVICE_ACCOUNT_TOKEN") != "":
		return authSettings{mode: authServiceAccount, token: getenv("OP_SERVICE_ACCOUNT_TOKEN")}, nil
	case getenv("OP_CONNECT_TOKEN") != "":
		return connectSettings(getenv("OP_CONNECT_TOKEN"), cfg, getenv), nil
	}
	return authSettings{}, errMissingCredentials
}

func connectSettings(token string, cfg onepasswordConfig, getenv func(string) string) authSettings {
	url := defaultConnectURL
	if nonEmpty(cfg.URL) {
		url = *cfg.URL
	} else if env := getenv("OP_CONNECT_HOST"); env != "" {
		url = env
	}
	return authSettings{mode: authConnect, token: token, url: url}
}

var sdkClientMu sync.Mutex

const sdkBackendCacheKey = "onepassword_sdk_backend"

func getBackend(ctx context.Context, d *plugin.QueryData) (vaultBackend, error) {
	auth, err := resolveAuth(GetConfig(d.Connection), os.Getenv)
	if err != nil {
		return nil, err
	}
	if auth.mode == authConnect {
		return connectBackend{client: connect.NewClient(auth.url, auth.token)}, nil
	}

	// SDK client start-up loads a WASM core and authenticates, so one client serves the connection.
	sdkClientMu.Lock()
	defer sdkClientMu.Unlock()
	if cached, ok := d.ConnectionCache.Get(ctx, sdkBackendCacheKey); ok {
		return cached.(sdkBackend), nil
	}
	client, err := opsdk.NewClient(ctx,
		opsdk.WithServiceAccountToken(auth.token),
		opsdk.WithIntegrationInfo("Steampipe 1Password plugin", "v1"),
	)
	if err != nil {
		return nil, err
	}
	backend := sdkBackend{vaults: client.Vaults(), items: client.Items()}
	if err := d.ConnectionCache.Set(ctx, sdkBackendCacheKey, backend); err != nil {
		plugin.Logger(ctx).Warn("onepassword.getBackend", "cache_error", err)
	}
	return backend, nil
}

type connectBackend struct {
	client connect.Client
}

func (b connectBackend) ListVaults(_ context.Context) ([]onepassword.Vault, error) {
	return b.client.GetVaults()
}

func (b connectBackend) GetVault(_ context.Context, id string) (*onepassword.Vault, error) {
	return b.client.GetVault(id)
}

func (b connectBackend) ListItems(_ context.Context, vaultID string) ([]onepassword.Item, error) {
	return b.client.GetItems(vaultID)
}

func (b connectBackend) GetItem(_ context.Context, itemID, vaultID string) (*onepassword.Item, error) {
	return b.client.GetItem(itemID, vaultID)
}

func (b connectBackend) ListFiles(_ context.Context, itemID, vaultID string) ([]onepassword.File, error) {
	return b.client.GetFiles(itemID, vaultID)
}

func (b connectBackend) GetFile(_ context.Context, fileID, itemID, vaultID string) (*onepassword.File, error) {
	return b.client.GetFile(fileID, itemID, vaultID)
}

func (b connectBackend) GetFileContent(_ context.Context, file *onepassword.File) ([]byte, error) {
	return b.client.GetFileContent(file)
}

type sdkVaultsAPI interface {
	List(ctx context.Context, params ...opsdk.VaultListParams) ([]opsdk.VaultOverview, error)
	GetOverview(ctx context.Context, vaultID string) (opsdk.VaultOverview, error)
}

type sdkItemsAPI interface {
	List(ctx context.Context, vaultID string, filters ...opsdk.ItemListFilter) ([]opsdk.ItemOverview, error)
	Get(ctx context.Context, vaultID string, itemID string) (opsdk.Item, error)
}

type sdkBackend struct {
	vaults sdkVaultsAPI
	items  sdkItemsAPI
}

func (b sdkBackend) ListVaults(ctx context.Context) ([]onepassword.Vault, error) {
	overviews, err := b.vaults.List(ctx)
	if err != nil {
		return nil, err
	}
	vaults := make([]onepassword.Vault, 0, len(overviews))
	for _, v := range overviews {
		vaults = append(vaults, vaultFromSDK(v))
	}
	return vaults, nil
}

func (b sdkBackend) GetVault(ctx context.Context, id string) (*onepassword.Vault, error) {
	overview, err := b.vaults.GetOverview(ctx, id)
	if err != nil {
		return nil, err
	}
	vault := vaultFromSDK(overview)
	return &vault, nil
}

func (b sdkBackend) ListItems(ctx context.Context, vaultID string) ([]onepassword.Item, error) {
	overviews, err := b.items.List(ctx, vaultID)
	if err != nil {
		return nil, err
	}
	items := make([]onepassword.Item, 0, len(overviews))
	for _, o := range overviews {
		items = append(items, itemFromSDKOverview(o))
	}
	return items, nil
}

func (b sdkBackend) GetItem(ctx context.Context, itemID, vaultID string) (*onepassword.Item, error) {
	item, err := b.items.Get(ctx, vaultID, itemID)
	if err != nil {
		return nil, err
	}
	converted := itemFromSDK(item)
	return &converted, nil
}

func (b sdkBackend) ListFiles(context.Context, string, string) ([]onepassword.File, error) {
	return nil, errFilesNeedConnect
}

func (b sdkBackend) GetFile(context.Context, string, string, string) (*onepassword.File, error) {
	return nil, errFilesNeedConnect
}

func (b sdkBackend) GetFileContent(context.Context, *onepassword.File) ([]byte, error) {
	return nil, errFilesNeedConnect
}

var sdkVaultTypes = map[opsdk.VaultType]onepassword.VaultType{
	opsdk.VaultTypePersonal:    onepassword.PersonalVault,
	opsdk.VaultTypeEveryone:    onepassword.EveryoneVault,
	opsdk.VaultTypeTransfer:    onepassword.TransferVault,
	opsdk.VaultTypeUserCreated: onepassword.UserCreatedVault,
}

func vaultFromSDK(v opsdk.VaultOverview) onepassword.Vault {
	vaultType, ok := sdkVaultTypes[v.VaultType]
	if !ok {
		vaultType = onepassword.UnknownVault
	}
	return onepassword.Vault{
		ID:             v.ID,
		Name:           v.Title,
		Description:    v.Description,
		AttrVersion:    int(v.AttributeVersion),
		ContentVersoin: int(v.ContentVersion),
		Items:          int(v.ActiveItemCount),
		Type:           vaultType,
		CreatedAt:      v.CreatedAt,
		UpdatedAt:      v.UpdatedAt,
	}
}

// The SDK and Connect name categories differently; rows use the Connect names so category
// filters behave the same whichever credential is configured.
var sdkItemCategories = map[opsdk.ItemCategory]onepassword.ItemCategory{
	opsdk.ItemCategoryLogin:                onepassword.Login,
	opsdk.ItemCategoryPassword:             onepassword.Password,
	opsdk.ItemCategoryAPICredentials:       onepassword.ApiCredential,
	opsdk.ItemCategoryServer:               onepassword.Server,
	opsdk.ItemCategoryDatabase:             onepassword.Database,
	opsdk.ItemCategoryCreditCard:           onepassword.CreditCard,
	opsdk.ItemCategoryMembership:           onepassword.Membership,
	opsdk.ItemCategoryPassport:             onepassword.Passport,
	opsdk.ItemCategorySoftwareLicense:      onepassword.SoftwareLicense,
	opsdk.ItemCategoryOutdoorLicense:       onepassword.OutdoorLicense,
	opsdk.ItemCategorySecureNote:           onepassword.SecureNote,
	opsdk.ItemCategoryRouter:               onepassword.WirelessRouter,
	opsdk.ItemCategoryBankAccount:          onepassword.BankAccount,
	opsdk.ItemCategoryDriverLicense:        onepassword.DriverLicense,
	opsdk.ItemCategoryIdentity:             onepassword.Identity,
	opsdk.ItemCategoryRewards:              onepassword.RewardProgram,
	opsdk.ItemCategoryDocument:             onepassword.Document,
	opsdk.ItemCategoryEmail:                onepassword.EmailAccount,
	opsdk.ItemCategorySocialSecurityNumber: onepassword.SocialSecurityNumber,
	opsdk.ItemCategoryMedicalRecord:        onepassword.MedicalRecord,
	opsdk.ItemCategorySSHKey:               onepassword.SSHKey,
	opsdk.ItemCategoryCryptoWallet:         "CRYPTO_WALLET",
	opsdk.ItemCategoryPerson:               "PERSON",
}

func categoryFromSDK(c opsdk.ItemCategory) onepassword.ItemCategory {
	if mapped, ok := sdkItemCategories[c]; ok {
		return mapped
	}
	return onepassword.Custom
}

func urlsFromSDK(websites []opsdk.Website) []onepassword.ItemURL {
	if len(websites) == 0 {
		return nil
	}
	urls := make([]onepassword.ItemURL, 0, len(websites))
	for _, w := range websites {
		urls = append(urls, onepassword.ItemURL{URL: w.URL, Label: w.Label})
	}
	return urls
}

func itemFromSDKOverview(o opsdk.ItemOverview) onepassword.Item {
	return onepassword.Item{
		ID:        o.ID,
		Title:     o.Title,
		Vault:     onepassword.ItemVault{ID: o.VaultID},
		Category:  categoryFromSDK(o.Category),
		Tags:      o.Tags,
		URLs:      urlsFromSDK(o.Websites),
		CreatedAt: o.CreatedAt,
		UpdatedAt: o.UpdatedAt,
	}
}

var sdkFieldTypes = map[opsdk.ItemFieldType]string{
	opsdk.ItemFieldTypeText:             "STRING",
	opsdk.ItemFieldTypeConcealed:        "CONCEALED",
	opsdk.ItemFieldTypeEmail:            "EMAIL",
	opsdk.ItemFieldTypeURL:              "URL",
	opsdk.ItemFieldTypeTOTP:             "OTP",
	opsdk.ItemFieldTypeDate:             "DATE",
	opsdk.ItemFieldTypeMonthYear:        "MONTH_YEAR",
	opsdk.ItemFieldTypeMenu:             "MENU",
	opsdk.ItemFieldTypePhone:            "PHONE",
	opsdk.ItemFieldTypeAddress:          "ADDRESS",
	opsdk.ItemFieldTypeCreditCardType:   "CREDIT_CARD_TYPE",
	opsdk.ItemFieldTypeCreditCardNumber: "CREDIT_CARD_NUMBER",
	opsdk.ItemFieldTypeReference:        "REFERENCE",
	opsdk.ItemFieldTypeSSHKey:           "SSHKEY",
}

func itemFromSDK(i opsdk.Item) onepassword.Item {
	item := onepassword.Item{
		ID:        i.ID,
		Title:     i.Title,
		Vault:     onepassword.ItemVault{ID: i.VaultID},
		Category:  categoryFromSDK(i.Category),
		Tags:      i.Tags,
		URLs:      urlsFromSDK(i.Websites),
		Version:   int(i.Version),
		CreatedAt: i.CreatedAt,
		UpdatedAt: i.UpdatedAt,
	}

	sections := map[string]*onepassword.ItemSection{}
	for _, s := range i.Sections {
		section := &onepassword.ItemSection{ID: s.ID, Label: s.Title}
		sections[s.ID] = section
		item.Sections = append(item.Sections, section)
	}

	for _, f := range i.Fields {
		fieldType, ok := sdkFieldTypes[f.FieldType]
		if !ok {
			fieldType = "UNKNOWN"
		}
		field := &onepassword.ItemField{ID: f.ID, Label: f.Title, Value: f.Value, Type: fieldType}
		if f.SectionID != nil {
			field.Section = sections[*f.SectionID]
		}
		// The SDK has no field purpose; Connect marks the built-in login fields with one, and the
		// category tables read username and password by it.
		if f.SectionID == nil && strings.EqualFold(f.ID, "username") {
			field.Purpose = "USERNAME"
		} else if f.SectionID == nil && strings.EqualFold(f.ID, "password") {
			field.Purpose = "PASSWORD"
		}
		item.Fields = append(item.Fields, field)
	}

	if i.Notes != "" {
		item.Fields = append(item.Fields, &onepassword.ItemField{ID: "notesPlain", Label: "notesPlain", Purpose: "NOTES", Type: "STRING", Value: i.Notes})
	}

	for _, f := range i.Files {
		item.Files = append(item.Files, &onepassword.File{ID: f.Attributes.ID, Name: f.Attributes.Name, Size: int(f.Attributes.Size), Section: sections[f.SectionID]})
	}
	if i.Document != nil {
		item.Files = append(item.Files, &onepassword.File{ID: i.Document.ID, Name: i.Document.Name, Size: int(i.Document.Size)})
	}
	return item
}
