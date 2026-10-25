package onepassword

import (
	"context"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

//// TABLE DEFINITION

func tableOnepasswordUser(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "onepassword_user",
		Description: "Users in the 1Password account, from the Users API.",
		List: &plugin.ListConfig{
			Hydrate: listUsers,
		},
		Columns: []*plugin.Column{
			{
				Name:        "id",
				Type:        proto.ColumnType_STRING,
				Description: "The unique identifier of the user.",
				Transform:   transform.FromField("ID"),
			},
			{
				Name:        "email",
				Type:        proto.ColumnType_STRING,
				Description: "The email address of the user.",
			},
			{
				Name:        "display_name",
				Type:        proto.ColumnType_STRING,
				Description: "The display name of the user.",
			},
			{
				Name:        "state",
				Type:        proto.ColumnType_STRING,
				Description: "The state of the user: ACTIVE or SUSPENDED. Null for users in other states, such as pending confirmation.",
				Transform:   transform.FromField("State").NullIfZero(),
			},
			{
				Name:        "create_time",
				Type:        proto.ColumnType_TIMESTAMP,
				Description: "When the user was created.",
			},
			{
				Name:        "path",
				Type:        proto.ColumnType_STRING,
				Description: "The resource path of the user, in the form accounts/{account}/users/{user}.",
			},

			/// Steampipe standard columns
			{
				Name:        "title",
				Type:        proto.ColumnType_STRING,
				Description: "The title of the user.",
				Transform:   transform.FromField("DisplayName"),
			},
		},
	}
}

//// LIST FUNCTION

func listUsers(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	client, err := getUsersClient(d)
	if err != nil {
		plugin.Logger(ctx).Error("onepassword_user.listUsers", "connection_error", err)
		return nil, err
	}
	err = client.listUsers(ctx, func(user apiUser) bool {
		d.StreamListItem(ctx, user)
		return d.RowsRemaining(ctx) != 0
	})
	if err != nil {
		plugin.Logger(ctx).Error("onepassword_user.listUsers", "api_error", err)
		return nil, err
	}
	return nil, nil
}
