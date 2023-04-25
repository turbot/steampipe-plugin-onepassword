package onepassword

import (
	"context"

	"github.com/1Password/connect-sdk-go/onepassword"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func tableOnepasswordItemFile(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "onepassword_item_file",
		Description: "Retrieve information about your files.",
		List: &plugin.ListConfig{
			ParentHydrate: listVaults,
			Hydrate:       listItemFiles,
			KeyColumns: []*plugin.KeyColumn{
				{
					Name:    "item_id",
					Require: plugin.Required,
				},
				{
					Name:    "vault_id",
					Require: plugin.Optional,
				},
			},
		},
		Columns: []*plugin.Column{
			{
				Name:        "id",
				Type:        proto.ColumnType_STRING,
				Description: "The ID of the Item.",
				Transform:   transform.FromField("ID"),
			},
			{
				Name:        "item_id",
				Type:        proto.ColumnType_STRING,
				Description: "The ID of the Item.",
				Transform:   transform.FromQual("item_id"),
			},
			{
				Name:        "name",
				Type:        proto.ColumnType_STRING,
				Description: "The title of this Item.",
			},
			{
				Name:        "content",
				Type:        proto.ColumnType_JSON,
				Description: "The parent vault ID of the Item.",
				Hydrate:     getFileContent,
				Transform:   transform.FromValue(),
			},
			{
				Name:        "vault_id",
				Type:        proto.ColumnType_STRING,
				Description: "The ID of the Item.",
				Transform:   transform.FromQual("vault_id"),
			},
			{
				Name:        "entropy",
				Type:        proto.ColumnType_DOUBLE,
				Description: "Is the item favorite.",
			},
			{
				Name:        "generate",
				Type:        proto.ColumnType_BOOL,
				Description: "The version of the item.",
			},
			{
				Name:        "purpose",
				Type:        proto.ColumnType_STRING,
				Description: "The category of the item.",
			},
			{
				Name:        "totp",
				Type:        proto.ColumnType_STRING,
				Description: "The category of the item.",
				Transform:   transform.FromField("TOTP"),
			},
			{
				Name:        "type",
				Type:        proto.ColumnType_STRING,
				Description: "The category of the item.",
			},
			{
				Name:        "recipe",
				Type:        proto.ColumnType_JSON,
				Description: "The category of the item.",
			},
			{
				Name:        "section",
				Type:        proto.ColumnType_JSON,
				Description: "UUID of the user that last edited the item.",
			},

			/// Steampipe standard columns
			{
				Name:        "title",
				Description: "Title of the resource.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ID"),
			},
		},
	}
}

func listItemFiles(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	item_id := d.EqualsQualString("item_id")
	vault := h.Item.(onepassword.Vault)
	vault_id := d.EqualsQuals["vault_id"].GetStringValue()

	// check if the item_id is empty or the provided vault_id is not matching with the parentHydrate
	if item_id == "" || (vault_id != "" && vault_id != vault.ID) {
		return nil, nil
	}

	client, err := getClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("onepassword_item_file.listItemFiles", "connection_error", err)
		return nil, err
	}

	files, err := client.GetFiles(item_id, vault.ID)
	if err != nil {
		plugin.Logger(ctx).Error("onepassword_item_file.listItemFiles", "api_error", err)
		return nil, err
	}

	for _, file := range files {
		d.StreamListItem(ctx, file)

		// Context can be cancelled due to manual cancellation or the limit has been hit
		if d.RowsRemaining(ctx) == 0 {
			return nil, nil
		}
	}

	return nil, nil
}

func getFileContent(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	file := h.Item.(onepassword.File)

	client, err := getClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("onepassword_item_file.getFileContent", "connection_error", err)
		return nil, err
	}

	content, err := client.GetFileContent(&file)
	if err != nil {
		plugin.Logger(ctx).Error("onepassword_item_file.getFileContent", "api_error", err)
		return nil, err
	}

	return string(content), nil
}
