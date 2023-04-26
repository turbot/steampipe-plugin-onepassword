package onepassword

import (
	"context"

	"github.com/1Password/connect-sdk-go/onepassword"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func tableOnepasswordItemLogin(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "onepassword_item_login",
		Description: "Retrieve information about your logins.",
		List: &plugin.ListConfig{
			ParentHydrate: listVaults,
			Hydrate:       listItemLogins,
			KeyColumns: []*plugin.KeyColumn{
				{
					Name:    "vault_id",
					Require: plugin.Optional,
				},
			},
		},
		Get: &plugin.GetConfig{
			Hydrate:    getItemLogin,
			KeyColumns: plugin.AllColumns([]string{"id", "vault_id"}),
		},
		Columns: []*plugin.Column{
			{
				Name:        "id",
				Type:        proto.ColumnType_STRING,
				Description: "The UUID of the item.",
				Transform:   transform.FromField("ID"),
			},
			{
				Name:        "vault_id",
				Type:        proto.ColumnType_STRING,
				Description: "The UUID of the vault the item is in.",
				Transform:   transform.FromField("Vault.ID"),
			},
			{
				Name:        "username",
				Type:        proto.ColumnType_STRING,
				Description: "The username stored in the item.",
				Hydrate:     getItemLogin,
			},
			{
				Name:        "password",
				Type:        proto.ColumnType_STRING,
				Description: "The password stored in the item.",
				Hydrate:     getItemLogin,
			},
			{
				Name:        "favorite",
				Type:        proto.ColumnType_BOOL,
				Description: "Whether the item is marked as a favorite.",
			},
			{
				Name:        "category",
				Type:        proto.ColumnType_STRING,
				Description: "The category of the item.",
			},
			{
				Name:        "created_at",
				Type:        proto.ColumnType_TIMESTAMP,
				Description: "Date and time when the item was created.",
			},
			{
				Name:        "last_edited_by",
				Type:        proto.ColumnType_STRING,
				Description: "UUID of the user that last edited the item.",
			},
			{
				Name:        "updated_at",
				Type:        proto.ColumnType_TIMESTAMP,
				Description: "Date and time when the item was last changed.",
			},
			{
				Name:        "version",
				Type:        proto.ColumnType_INT,
				Description: "The version of the item.",
			},
			{
				Name:        "sections",
				Type:        proto.ColumnType_JSON,
				Description: "The sections of the item.",
				Hydrate:     getItemLogin,
			},
			{
				Name:        "fields",
				Type:        proto.ColumnType_JSON,
				Description: "The fields of the item.",
				Hydrate:     getItemLogin,
			},
			{
				Name:        "files",
				Type:        proto.ColumnType_JSON,
				Description: "The files of the item.",
				Hydrate:     getItemLogin,
			},
			{
				Name:        "tags",
				Type:        proto.ColumnType_JSON,
				Description: "An array of strings of the tags assigned to the item.",
			},
			{
				Name:        "urls",
				Type:        proto.ColumnType_JSON,
				Description: "Array of URL objects containing URLs for the item.",
				Transform:   transform.FromField("URLs"),
			},

			/// Steampipe standard columns
			{
				Name:        "title",
				Description: "The title of the item.",
				Type:        proto.ColumnType_STRING,
			},
		},
	}
}

type ItemLogin struct {
	Username string
	Password string
	onepassword.Item
}

func listItemLogins(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	vault := h.Item.(onepassword.Vault)
	vault_id := d.EqualsQualString("vault_id")

	// check if the provided vault_id is not matching with the parentHydrate
	if vault_id != "" && vault_id != vault.ID {
		return nil, nil
	}

	client, err := getClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("onepassword_item_login.listItemLogins", "connection_error", err)
		return nil, err
	}

	items, err := client.GetItems(vault.ID)
	if err != nil {
		plugin.Logger(ctx).Error("onepassword_item_login.listItemLogins", "api_error", err)
		return nil, err
	}

	for _, item := range items {

		// restricting data based on the item category LOGIN
		if item.Category == "LOGIN" {
			d.StreamListItem(ctx, ItemLogin{"", "", item})
		}
		// Context can be cancelled due to manual cancellation or the limit has been hit
		if d.RowsRemaining(ctx) == 0 {
			return nil, nil
		}
	}

	return nil, nil
}

func getItemLogin(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	var id, vault_id string
	if h.Item != nil {
		id = h.Item.(ItemLogin).Item.ID
		vault_id = h.Item.(ItemLogin).Item.Vault.ID
	} else {
		id = d.EqualsQualString("id")
		vault_id = d.EqualsQualString("vault_id")
	}

	// Check if id or vault_id is empty
	if id == "" || vault_id == "" {
		return nil, nil
	}

	client, err := getClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("onepassword_item_login.getItemLogin", "connection_error", err)
		return nil, err
	}

	item, err := client.GetItem(id, vault_id)
	if err != nil {
		plugin.Logger(ctx).Error("onepassword_item_login.getItemLogin", "api_error", err)
		return nil, err
	}

	// restricting data based on the item category LOGIN
	var username, password string
	if item.Category == "LOGIN" {
		for _, field := range item.Fields {
			if field.ID == "username" && field.Purpose == "USERNAME" {
				username = field.Value
			}else if field.ID == "password" && field.Purpose == "PASSWORD" {
				password = field.Value
			}
		}
		return ItemLogin{username, password, *item}, nil
	}

	return nil, nil
}
