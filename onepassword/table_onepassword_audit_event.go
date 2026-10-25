package onepassword

import (
	"context"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

//// TABLE DEFINITION

func tableOnepasswordAuditEvent(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "onepassword_audit_event",
		Description: "Actions team members performed in the 1Password account, from the Events API. Defaults to the last 24 hours unless timestamp is constrained.",
		List: &plugin.ListConfig{
			Hydrate:    listAuditEvents,
			KeyColumns: timestampKeyColumns(),
		},
		Columns: []*plugin.Column{
			{
				Name:        "uuid",
				Type:        proto.ColumnType_STRING,
				Description: "The UUID of the event.",
				Transform:   transform.FromField("UUID"),
			},
			{
				Name:        "timestamp",
				Type:        proto.ColumnType_TIMESTAMP,
				Description: "When the action was performed.",
			},
			{
				Name:        "action",
				Type:        proto.ColumnType_STRING,
				Description: "The action performed, such as create, suspend, disblmfa or disblsso.",
			},
			{
				Name:        "object_type",
				Type:        proto.ColumnType_STRING,
				Description: "The type of object the action was performed on, such as user, vault or account.",
			},
			{
				Name:        "object_uuid",
				Type:        proto.ColumnType_STRING,
				Description: "The UUID of the object the action was performed on.",
				Transform:   transform.FromField("ObjectUUID").NullIfZero(),
			},
			{
				Name:        "object_details",
				Type:        proto.ColumnType_JSON,
				Description: "The user the action was performed on, when the object is a user.",
			},
			{
				Name:        "actor_uuid",
				Type:        proto.ColumnType_STRING,
				Description: "The UUID of the user who performed the action.",
				Transform:   transform.FromField("ActorUUID").NullIfZero(),
			},
			{
				Name:        "actor_name",
				Type:        proto.ColumnType_STRING,
				Description: "The name of the user who performed the action.",
				Transform:   transform.FromField("ActorDetails.Name").NullIfZero(),
			},
			{
				Name:        "actor_email",
				Type:        proto.ColumnType_STRING,
				Description: "The email address of the user who performed the action.",
				Transform:   transform.FromField("ActorDetails.Email").NullIfZero(),
			},
			{
				Name:        "actor_type",
				Type:        proto.ColumnType_STRING,
				Description: "The type of the actor: user or external_user.",
				Transform:   transform.FromField("ActorType").NullIfZero(),
			},
			{
				Name:        "actor_account_uuid",
				Type:        proto.ColumnType_STRING,
				Description: "The account UUID of the actor, for managed service provider accounts.",
				Transform:   transform.FromField("ActorAccountUUID").NullIfZero(),
			},
			{
				Name:        "aux_id",
				Type:        proto.ColumnType_INT,
				Description: "The ID of an additional object related to the action.",
				Transform:   transform.FromField("AuxID"),
			},
			{
				Name:        "aux_uuid",
				Type:        proto.ColumnType_STRING,
				Description: "The UUID of an additional object related to the action.",
				Transform:   transform.FromField("AuxUUID").NullIfZero(),
			},
			{
				Name:        "aux_details",
				Type:        proto.ColumnType_JSON,
				Description: "The user related to the action, when the additional object is a user.",
			},
			{
				Name:        "aux_info",
				Type:        proto.ColumnType_STRING,
				Description: "Additional information about the action.",
				Transform:   transform.FromField("AuxInfo").NullIfZero(),
			},
			{
				Name:        "session",
				Type:        proto.ColumnType_JSON,
				Description: "The session the action was performed in: uuid, login_time, device_uuid and ip.",
			},
			{
				Name:        "location",
				Type:        proto.ColumnType_JSON,
				Description: "The approximate location the action was performed from.",
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

func listAuditEvents(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	return streamEvents[auditEvent](ctx, d, "onepassword_audit_event.listAuditEvents", auditEventsPath)
}
