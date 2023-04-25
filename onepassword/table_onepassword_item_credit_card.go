package onepassword

import (
	"context"

	"github.com/1Password/connect-sdk-go/onepassword"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func tableOnepasswordItemCreditCard(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "onepassword_item_credit_card",
		Description: "Retrieve information about your credit cards.",
		List: &plugin.ListConfig{
			ParentHydrate: listVaults,
			Hydrate:       listItemCreditCards,
			KeyColumns: []*plugin.KeyColumn{
				{
					Name:    "vault_id",
					Require: plugin.Optional,
				},
			},
		},
		Get: &plugin.GetConfig{
			Hydrate:    getItemCreditCard,
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
				Name:        "favorite",
				Type:        proto.ColumnType_BOOL,
				Description: "Whether the item is marked as a favorite.",
			},
			{
				Name:        "card_holder",
				Type:        proto.ColumnType_STRING,
				Description: "The card holder’s name for the credit card.",
				Hydrate:     getItemCreditCard,
			},
			{
				Name:        "category",
				Type:        proto.ColumnType_STRING,
				Description: "The category of the item.",
			},
			{
				Name:        "credit_card_number",
				Type:        proto.ColumnType_STRING,
				Description: "The credit card number.",
				Hydrate:     getItemCreditCard,
			},
			{
				Name:        "created_at",
				Type:        proto.ColumnType_TIMESTAMP,
				Description: "Date and time when the item was created.",
			},
			{
				Name:        "expiry_date",
				Type:        proto.ColumnType_TIMESTAMP,
				Description: "The expiry date of the credit card.",
				Hydrate:     getItemCreditCard,
			},
			{
				Name:        "last_edited_by",
				Type:        proto.ColumnType_STRING,
				Description: "UUID of the user that last edited the item.",
			},
			{
				Name:        "type",
				Type:        proto.ColumnType_STRING,
				Description: "The type of the credit card.",
				Hydrate:     getItemCreditCard,
			},
			{
				Name:        "updated_at",
				Type:        proto.ColumnType_TIMESTAMP,
				Description: "Date and time when the vault or its contents were last changed.",
			},
			{
				Name:        "valid_from",
				Type:        proto.ColumnType_TIMESTAMP,
				Description: "The valid from date of the credit card.",
				Hydrate:     getItemCreditCard,
			},
			{
				Name:        "verification_number",
				Type:        proto.ColumnType_STRING,
				Description: "The cvv number of the credit card.",
				Hydrate:     getItemCreditCard,
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
				Hydrate:     getItemCreditCard,
			},
			{
				Name:        "fields",
				Type:        proto.ColumnType_JSON,
				Description: "The fields of the item.",
				Hydrate:     getItemCreditCard,
			},
			{
				Name:        "files",
				Type:        proto.ColumnType_JSON,
				Description: "The files of the item.",
				Hydrate:     getItemCreditCard,
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

type ItemCreditCard struct {
	CardHolder         string
	Type               string
	CreditCardNumber   string
	VerificationNumber string
	ExpiryDate         string
	ValidFrom          string
	onepassword.Item
}

func listItemCreditCards(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	vault := h.Item.(onepassword.Vault)
	vault_id := d.EqualsQualString("vault_id")

	// check if the provided vault_id is not matching with the parentHydrate
	if vault_id != "" && vault_id != vault.ID {
		return nil, nil
	}

	client, err := getClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("onepassword_item_credit_card.listItemCreditCards", "connection_error", err)
		return nil, err
	}

	items, err := client.GetItems(vault.ID)
	if err != nil {
		plugin.Logger(ctx).Error("onepassword_item_credit_card.listItemCreditCards", "api_error", err)
		return nil, err
	}

	for _, item := range items {

		// restricting data based on the item category CREDIT_CARD
		if item.Category == "CREDIT_CARD" {
			d.StreamListItem(ctx, ItemCreditCard{"", "", "", "", "", "", item})
		}
		// Context can be cancelled due to manual cancellation or the limit has been hit
		if d.RowsRemaining(ctx) == 0 {
			return nil, nil
		}
	}

	return nil, nil
}

func getItemCreditCard(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	var id, vault_id string
	if h.Item != nil {
		id = h.Item.(ItemCreditCard).Item.ID
		vault_id = h.Item.(ItemCreditCard).Item.Vault.ID
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
		plugin.Logger(ctx).Error("onepassword_item.getItem", "connection_error", err)
		return nil, err
	}

	item, err := client.GetItem(id, vault_id)
	if err != nil {
		plugin.Logger(ctx).Error("onepassword_item.getItem", "api_error", err)
		return nil, err
	}

	// restricting data based on the item category CREDIT_CARD
	var cardholder, cctype, ccnum, cvv, expiry, valid_from string
	if item.Category == "CREDIT_CARD" {
		for _, field := range item.Fields {
			if field.ID == "cardholder" {
				cardholder = field.Value
			}
			if field.ID == "type" {
				cctype = field.Value
			}
			if field.ID == "ccnum" {
				ccnum = field.Value
			}
			if field.ID == "cvv" {
				cvv = field.Value
			}
			if field.ID == "expiry" {
				expiry = field.Value
			}
			if field.ID == "validFrom" {
				valid_from = field.Value
			}
		}
		return ItemCreditCard{cardholder, cctype, ccnum, cvv, expiry, valid_from, *item}, nil
	}

	return nil, nil
}
