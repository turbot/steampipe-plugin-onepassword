package onepassword

import (
	"context"

	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func Plugin(ctx context.Context) *plugin.Plugin {
	p := &plugin.Plugin{
		Name: "steampipe-plugin-onepassword",
		ConnectionConfigSchema: &plugin.ConnectionConfigSchema{
			NewInstance: ConfigInstance,
		},
		DefaultTransform: transform.FromCamel(),
		TableMap: map[string]*plugin.Table{
			"onepassword_item":                  tableOnepasswordItem(ctx),
			"onepassword_item_api_credential":   tableOnepasswordItemAPICredential(ctx),
			"onepassword_item_credit_card":      tableOnepasswordItemCreditCard(ctx),
			"onepassword_item_file":             tableOnepasswordItemFile(ctx),
			"onepassword_item_identity":         tableOnepasswordItemIdentity(ctx),
			"onepassword_item_login":            tableOnepasswordItemLogin(ctx),
			"onepassword_item_medical_record":   tableOnepasswordItemMedicalRecord(ctx),
			"onepassword_item_password":         tableOnepasswordItemPassword(ctx),
			"onepassword_item_secure_note":      tableOnepasswordItemSecureNote(ctx),
			"onepassword_item_software_license": tableOnepasswordItemSoftwareLicense(ctx),
			"onepassword_vault":                 tableOnepasswordVault(ctx),
		},
	}
	return p
}
