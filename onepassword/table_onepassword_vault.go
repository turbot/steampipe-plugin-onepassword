package onepassword

import (
	"context"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func tableOnepasswordVault(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "onepassword_vault",
		Description: "Retrieve information about your vaults.",
		List: &plugin.ListConfig{
			Hydrate: listVaults,
		},
		Get: &plugin.GetConfig{
			Hydrate:    getVault,
			KeyColumns: plugin.SingleColumn("id"),
		},
		Columns: []*plugin.Column{
			{
				Name:        "id",
				Type:        proto.ColumnType_STRING,
				Description: "The UUID of the vault.",
				Transform:   transform.FromField("ID"),
			},
			{
				Name:        "name",
				Type:        proto.ColumnType_STRING,
				Description: "The name of the vault.",
			},
			{
				Name:        "attr_version",
				Type:        proto.ColumnType_INT,
				Description: "The version of the vault metadata.",
			},
			{
				Name:        "created_at",
				Type:        proto.ColumnType_TIMESTAMP,
				Description: "Date and time when the vault was created.",
			},
			{
				Name:        "content_version",
				Type:        proto.ColumnType_INT,
				Description: "The version of the vault contents.",
				Transform:   transform.FromField("ContentVersoin"),
				// NOTE: The typo 'ContentVersoin' is from the SDK, if and when the SDK is updated, we need to sync with it.
				// Ref: https://github.com/1Password/connect-sdk-go/blob/ac4a2c9c017fcfb76ac97b63b145121926e06fd3/onepassword/vaults.go#L15
			},
			{
				Name:        "description",
				Type:        proto.ColumnType_STRING,
				Description: "The description for the vault.",
			},
			{
				Name:        "items",
				Type:        proto.ColumnType_INT,
				Description: "Number of active items in the vault.",
			},
			{
				Name:        "type",
				Type:        proto.ColumnType_STRING,
				Description: "The type of vault. Possible values are EVERYONE, PERSONAL and USER_CREATED.",
			},
			{
				Name:        "updated_at",
				Type:        proto.ColumnType_TIMESTAMP,
				Description: "Date and time when the vault or its contents were last changed.",
			},

			/// Steampipe standard columns
			{
				Name:        "title",
				Description: "The title of the vault.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Name"),
			},
		},
	}
}

func listVaults(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	client, err := getClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("onepassword_vault.listVaults", "connection_error", err)
		return nil, err
	}

	vaults, err := client.GetVaults()
	if err != nil {
		plugin.Logger(ctx).Error("onepassword_vault.listVaults", "query_error", err)
		return nil, err
	}

	for _, vault := range vaults {
		d.StreamListItem(ctx, vault)

		// Context can be cancelled due to manual cancellation or the limit has been hit
		if d.RowsRemaining(ctx) == 0 {
			return nil, nil
		}
	}

	return nil, nil
}

func getVault(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	id := d.EqualsQuals["id"].GetStringValue()

	// Check if id is empty
	if id == "" {
		return nil, nil
	}

	client, err := getClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("onepassword_vault.getVault", "connection_error", err)
		return nil, err
	}

	vault, err := client.GetVault(id)
	if err != nil {
		plugin.Logger(ctx).Error("onepassword_vault.getVault", "api_error", err)
		return nil, err
	}

	return vault, nil
}
