package onepassword

import (
	"context"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func tableOnepasswordItemField(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "onepassword_item_field",
		Description: "Retrieve information about your item fields.",
		List: &plugin.ListConfig{
			Hydrate: listItemFields,
			KeyColumns: []*plugin.KeyColumn{
				{
					Name:    "item_id",
					Require: plugin.Required,
				},
				{
					Name:    "vault_id",
					Require: plugin.Required,
				},
			},
		},
		Columns: []*plugin.Column{
			{
				Name:        "id",
				Type:        proto.ColumnType_STRING,
				Description: "The ID of the Item.",
			},
			{
				Name:        "item_id",
				Type:        proto.ColumnType_STRING,
				Description: "The ID of the Item.",
				Transform:   transform.FromQual("item_id"),
			},
			{
				Name:        "label",
				Type:        proto.ColumnType_STRING,
				Description: "The title of this Item.",
			},
			{
				Name:        "value",
				Type:        proto.ColumnType_STRING,
				Description: "The parent vault ID of the Item.",
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

func listItemFields(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	item_id := d.EqualsQualString("item_id")
	vault_id := d.EqualsQualString("vault_id")

	client, err := getClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("onepassword_Item.listItemFields", "connection_error", err)
		return nil, err
	}

	item, err := client.GetItem(item_id, vault_id)
	if err != nil {
		plugin.Logger(ctx).Error("onepassword_item.listItemFields", "api_error", err)
		return nil, err
	}

	for _, field := range item.Fields {
		d.StreamListItem(ctx, field)

		// Context can be cancelled due to manual cancellation or the limit has been hit
		if d.RowsRemaining(ctx) == 0 {
			return nil, nil
		}
	}

	return nil, nil
}
