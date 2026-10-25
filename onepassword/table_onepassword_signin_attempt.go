package onepassword

import (
	"context"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

//// TABLE DEFINITION

func tableOnepasswordSigninAttempt(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "onepassword_signin_attempt",
		Description: "Sign-in attempts to the 1Password account, from the Events API. Defaults to the last 24 hours unless timestamp is constrained.",
		List: &plugin.ListConfig{
			Hydrate:    listSigninAttempts,
			KeyColumns: timestampKeyColumns(),
		},
		Columns: []*plugin.Column{
			{
				Name:        "uuid",
				Type:        proto.ColumnType_STRING,
				Description: "The UUID of the sign-in attempt.",
				Transform:   transform.FromField("UUID"),
			},
			{
				Name:        "session_uuid",
				Type:        proto.ColumnType_STRING,
				Description: "The UUID of the session that created the sign-in attempt.",
				Transform:   transform.FromField("SessionUUID").NullIfZero(),
			},
			{
				Name:        "timestamp",
				Type:        proto.ColumnType_TIMESTAMP,
				Description: "When the sign-in attempt happened.",
			},
			{
				Name:        "category",
				Type:        proto.ColumnType_STRING,
				Description: "The outcome: success, credentials_failed, mfa_failed, sso_failed, modern_version_failed, firewall_failed or firewall_reported_success.",
			},
			{
				Name:        "type",
				Type:        proto.ColumnType_STRING,
				Description: "The detailed reason, such as credentials_ok, password_secret_bad or country_blocked.",
			},
			{
				Name:        "country",
				Type:        proto.ColumnType_STRING,
				Description: "The ISO 3166 country code of the attempt.",
				Transform:   transform.FromField("Country").NullIfZero(),
			},
			{
				Name:        "details",
				Type:        proto.ColumnType_JSON,
				Description: "Additional information about the attempt.",
			},
			{
				Name:        "target_user_uuid",
				Type:        proto.ColumnType_STRING,
				Description: "The UUID of the user the attempt was for.",
				Transform:   transform.FromField("TargetUser.UUID").NullIfZero(),
			},
			{
				Name:        "target_user_name",
				Type:        proto.ColumnType_STRING,
				Description: "The name of the user the attempt was for.",
				Transform:   transform.FromField("TargetUser.Name").NullIfZero(),
			},
			{
				Name:        "target_user_email",
				Type:        proto.ColumnType_STRING,
				Description: "The email address of the user the attempt was for.",
				Transform:   transform.FromField("TargetUser.Email").NullIfZero(),
			},
			{
				Name:        "target_user_type",
				Type:        proto.ColumnType_STRING,
				Description: "The type of the user the attempt was for: user or external_user.",
				Transform:   transform.FromField("TargetUser.UserType").NullIfZero(),
			},
			{
				Name:        "client_ip_address",
				Type:        proto.ColumnType_IPADDR,
				Description: "The IP address the attempt came from.",
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
				Description: "The approximate location of the attempt.",
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

func listSigninAttempts(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	return streamEvents[signInAttempt](ctx, d, "onepassword_signin_attempt.listSigninAttempts", signInAttemptsPath)
}

func timestampKeyColumns() []*plugin.KeyColumn {
	return []*plugin.KeyColumn{
		{Name: "timestamp", Require: plugin.Optional, Operators: []string{">", ">=", "<", "<="}},
	}
}
