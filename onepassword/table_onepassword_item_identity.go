package onepassword

import (
	"context"

	"github.com/1Password/connect-sdk-go/onepassword"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

//// TABLE DEFINITION

func tableOnepasswordItemIdentity(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "onepassword_item_identity",
		Description: "Retrieve information about your identities.",
		List: &plugin.ListConfig{
			ParentHydrate: listVaults,
			Hydrate:       listItemIdentities,
			KeyColumns: []*plugin.KeyColumn{
				{
					Name:    "vault_id",
					Require: plugin.Optional,
				},
			},
		},
		Get: &plugin.GetConfig{
			Hydrate:    getItemIdentity,
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
				Name:        "first_name",
				Type:        proto.ColumnType_STRING,
				Description: "The first name of the identity.",
				Hydrate:     getItemIdentity,
			},
			{
				Name:        "initial",
				Type:        proto.ColumnType_STRING,
				Description: "The initial of the identity.",
				Hydrate:     getItemIdentity,
			},
			{
				Name:        "last_name",
				Type:        proto.ColumnType_STRING,
				Description: "The last name of the identity.",
				Hydrate:     getItemIdentity,
			},
			{
				Name:        "birth_date",
				Type:        proto.ColumnType_TIMESTAMP,
				Description: "The birth date of the identity.",
				Hydrate:     getItemIdentity,
				Transform:   transform.FromField("BirthDate").Transform(transform.UnixToTimestamp),
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
				Name:        "company",
				Type:        proto.ColumnType_STRING,
				Description: "The company of the identity.",
				Hydrate:     getItemIdentity,
			},
			{
				Name:        "department",
				Type:        proto.ColumnType_STRING,
				Description: "The department of the identity.",
				Hydrate:     getItemIdentity,
			},
			{
				Name:        "favorite",
				Type:        proto.ColumnType_BOOL,
				Description: "Whether the item is marked as a favorite.",
			},
			{
				Name:        "gender",
				Type:        proto.ColumnType_STRING,
				Description: "The gender of the identity.",
				Hydrate:     getItemIdentity,
			},
			{
				Name:        "job_title",
				Type:        proto.ColumnType_STRING,
				Description: "The job title of the identity.",
				Hydrate:     getItemIdentity,
			},
			{
				Name:        "last_edited_by",
				Type:        proto.ColumnType_STRING,
				Description: "UUID of the user that last edited the item.",
			},
			{
				Name:        "occupation",
				Type:        proto.ColumnType_STRING,
				Description: "The occupation of the identity.",
				Hydrate:     getItemIdentity,
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
				Name:        "trashed",
				Type:        proto.ColumnType_BOOL,
				Description: "Checks if the item is trashed.",
			},
			{
				Name:        "sections",
				Type:        proto.ColumnType_JSON,
				Description: "The sections of the item.",
				Hydrate:     getItemIdentity,
			},
			{
				Name:        "fields",
				Type:        proto.ColumnType_JSON,
				Description: "The fields of the item.",
				Hydrate:     getItemIdentity,
			},
			{
				Name:        "files",
				Type:        proto.ColumnType_JSON,
				Description: "The files of the item.",
				Hydrate:     getItemIdentity,
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

type ItemIdentity struct {
	FirstName  string
	Initial    string
	LastName   string
	Gender     string
	BirthDate  string
	Occupation string
	Company    string
	Department string
	JobTitle   string
	onepassword.Item
}

/// LIST FUNCTION

func listItemIdentities(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	vault := h.Item.(onepassword.Vault)
	vault_id := d.EqualsQualString("vault_id")

	// check if the provided vault_id is not matching with the parentHydrate
	if vault_id != "" && vault_id != vault.ID {
		return nil, nil
	}

	client, err := getClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("onepassword_item_identity.listItemIdentities", "connection_error", err)
		return nil, err
	}

	items, err := client.GetItems(vault.ID)
	if err != nil {
		plugin.Logger(ctx).Error("onepassword_item_identity.listItemIdentities", "api_error", err)
		return nil, err
	}

	for _, item := range items {

		// restricting data based on the item category IDENTITY
		if item.Category == "IDENTITY" {
			d.StreamListItem(ctx, ItemIdentity{"", "", "", "", "", "", "", "", "", item})
		}
		// Context can be cancelled due to manual cancellation or the limit has been hit
		if d.RowsRemaining(ctx) == 0 {
			return nil, nil
		}
	}

	return nil, nil
}

/// HYDRATE FUNCTION

func getItemIdentity(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	var id, vault_id string
	if h.Item != nil {
		id = h.Item.(ItemIdentity).Item.ID
		vault_id = h.Item.(ItemIdentity).Item.Vault.ID
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
		plugin.Logger(ctx).Error("onepassword_item_identity.getItemIdentity", "connection_error", err)
		return nil, err
	}

	item, err := client.GetItem(id, vault_id)
	if err != nil {
		plugin.Logger(ctx).Error("onepassword_item_identity.getItemIdentity", "api_error", err)
		return nil, err
	}

	// restricting data based on the item category IDENTITY
	var firstname, initial, lastname, birthdate, gender, occupation, company, department, jobtitle string
	if item.Category == "IDENTITY" {
		for _, field := range item.Fields {
			if field.ID == "firstname" {
				firstname = field.Value
			} else if field.ID == "initial" {
				initial = field.Value
			} else if field.ID == "lastname" {
				lastname = field.Value
			} else if field.ID == "gender" {
				gender = field.Value
			} else if field.ID == "birthdate" {
				birthdate = field.Value
			} else if field.ID == "occupation" {
				occupation = field.Value
			} else if field.ID == "company" {
				company = field.Value
			} else if field.ID == "department" {
				department = field.Value
			} else if field.ID == "jobtitle" {
				jobtitle = field.Value
			}
		}
		return ItemIdentity{firstname, initial, lastname, gender, birthdate, occupation, company, department, jobtitle, *item}, nil
	}

	return nil, nil
}
