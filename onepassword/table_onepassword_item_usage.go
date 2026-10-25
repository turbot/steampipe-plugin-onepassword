package onepassword

import (
	"context"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

//// TABLE DEFINITION

func tableOnepasswordItemUsage(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "onepassword_item_usage",
		Description: "Uses of items in shared vaults, from the Events API. Defaults to the last 24 hours unless timestamp is constrained.",
		List: &plugin.ListConfig{
			Hydrate:    listItemUsages,
			KeyColumns: timestampKeyColumns(),
		},
		Columns: []*plugin.Column{
			{
				Name:        "uuid",
				Type:        proto.ColumnType_STRING,
				Description: "The UUID of the usage event.",
				Transform:   transform.FromField("UUID"),
			},
			{
				Name:        "timestamp",
				Type:        proto.ColumnType_TIMESTAMP,
				Description: "When the item was used.",
			},
			{
				Name:        "action",
				Type:        proto.ColumnType_STRING,
				Description: "How the item was used, such as fill, reveal, secure-copy, export or share.",
			},
			{
				Name:        "item_uuid",
				Type:        proto.ColumnType_STRING,
				Description: "The UUID of the item used.",
				Transform:   transform.FromField("ItemUUID"),
			},
			{
				Name:        "vault_uuid",
				Type:        proto.ColumnType_STRING,
				Description: "The UUID of the vault the item is in.",
				Transform:   transform.FromField("VaultUUID"),
			},
			{
				Name:        "used_version",
				Type:        proto.ColumnType_INT,
				Description: "The version of the item that was used.",
			},
			{
				Name:        "user_uuid",
				Type:        proto.ColumnType_STRING,
				Description: "The UUID of the user who used the item.",
				Transform:   transform.FromField("User.UUID").NullIfZero(),
			},
			{
				Name:        "user_name",
				Type:        proto.ColumnType_STRING,
				Description: "The name of the user who used the item.",
				Transform:   transform.FromField("User.Name").NullIfZero(),
			},
			{
				Name:        "user_email",
				Type:        proto.ColumnType_STRING,
				Description: "The email address of the user who used the item.",
				Transform:   transform.FromField("User.Email").NullIfZero(),
			},
			{
				Name:        "client_ip_address",
				Type:        proto.ColumnType_IPADDR,
				Description: "The IP address the item was used from.",
				Transform:   transform.FromField("Client.IPAddress").NullIfZero(),
			},
			{
				Name:        "client",
				Type:        proto.ColumnType_JSON,
				Description: "The 1Password app, platform and operating system used.",
			},
			{
				Name:        "location",
				Type:        proto.ColumnType_JSON,
				Description: "The approximate location the item was used from.",
			},
			{
				Name:        "account_uuid",
				Type:        proto.ColumnType_STRING,
				Description: "The UUID of the account, for managed service provider accounts.",
				Transform:   transform.FromField("AccountUUID").NullIfZero(),
			},
		},
	}
}

//// LIST FUNCTION

func listItemUsages(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	return streamEvents[itemUsage](ctx, d, "onepassword_item_usage.listItemUsages", itemUsagesPath)
}
