package onepassword

import (
	"context"

	"github.com/1Password/connect-sdk-go/onepassword"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func tableOnepasswordItemMedicalRecord(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "onepassword_item_medical_record",
		Description: "Retrieve information about your medical records.",
		List: &plugin.ListConfig{
			ParentHydrate: listVaults,
			Hydrate:       listItemMedicalRecords,
			KeyColumns: []*plugin.KeyColumn{
				{
					Name:    "vault_id",
					Require: plugin.Optional,
				},
			},
		},
		Get: &plugin.GetConfig{
			Hydrate:    getItemMedicalRecord,
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
				Name:        "date",
				Type:        proto.ColumnType_TIMESTAMP,
				Description: "The medical record date.",
				Hydrate:     getItemMedicalRecord,
				Transform:   transform.FromField("Date").Transform(transform.UnixToTimestamp),
			},
			{
				Name:        "dosage",
				Type:        proto.ColumnType_STRING,
				Description: "Dosage specified in the medical record.",
				Hydrate:     getItemMedicalRecord,
			},
			{
				Name:        "healthcare_professional",
				Type:        proto.ColumnType_STRING,
				Description: "The healthcare professional of the medical record.",
				Hydrate:     getItemMedicalRecord,
			},
			{
				Name:        "last_edited_by",
				Type:        proto.ColumnType_STRING,
				Description: "UUID of the user that last edited the item.",
			},
			{
				Name:        "location",
				Type:        proto.ColumnType_STRING,
				Description: "Location specified in the medical record.",
				Hydrate:     getItemMedicalRecord,
			},
			{
				Name:        "medication",
				Type:        proto.ColumnType_STRING,
				Description: "Medication specified in the medical record.",
				Hydrate:     getItemMedicalRecord,
			},
			{
				Name:        "medication_notes",
				Type:        proto.ColumnType_STRING,
				Description: "Medication notes specified in the medical record.",
				Hydrate:     getItemMedicalRecord,
			},
			{
				Name:        "patient",
				Type:        proto.ColumnType_STRING,
				Description: "The patient details are specified in the medical record.",
				Hydrate:     getItemMedicalRecord,
			},
			{
				Name:        "reason",
				Type:        proto.ColumnType_STRING,
				Description: "The reason specified in the medical record.",
				Hydrate:     getItemMedicalRecord,
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
				Hydrate:     getItemMedicalRecord,
			},
			{
				Name:        "fields",
				Type:        proto.ColumnType_JSON,
				Description: "The fields of the item.",
				Hydrate:     getItemMedicalRecord,
			},
			{
				Name:        "files",
				Type:        proto.ColumnType_JSON,
				Description: "The files of the item.",
				Hydrate:     getItemMedicalRecord,
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

type ItemMedicalRecord struct {
	Date                   string
	Dosage                 string
	HealthcareProfessional string
	Location               string
	Medication             string
	MedicationNotes        string
	Patient                string
	Reason                 string
	onepassword.Item
}

func listItemMedicalRecords(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	vault := h.Item.(onepassword.Vault)
	vault_id := d.EqualsQualString("vault_id")

	// check if the provided vault_id is not matching with the parentHydrate
	if vault_id != "" && vault_id != vault.ID {
		return nil, nil
	}

	client, err := getClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("onepassword_item_medical_record.listItemMedicalRecords", "connection_error", err)
		return nil, err
	}

	items, err := client.GetItems(vault.ID)
	if err != nil {
		plugin.Logger(ctx).Error("onepassword_item_medical_record.listItemMedicalRecords", "api_error", err)
		return nil, err
	}

	for _, item := range items {

		// restricting data based on the item category MEDICAL_RECORD
		if item.Category == "MEDICAL_RECORD" {
			d.StreamListItem(ctx, ItemMedicalRecord{"", "", "", "", "", "", "", "", item})
		}
		// Context can be cancelled due to manual cancellation or the limit has been hit
		if d.RowsRemaining(ctx) == 0 {
			return nil, nil
		}
	}

	return nil, nil
}

func getItemMedicalRecord(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	var id, vault_id string
	if h.Item != nil {
		id = h.Item.(ItemMedicalRecord).Item.ID
		vault_id = h.Item.(ItemMedicalRecord).Item.Vault.ID
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
		plugin.Logger(ctx).Error("onepassword_item_medical_record.getItemMedicalRecord", "connection_error", err)
		return nil, err
	}

	item, err := client.GetItem(id, vault_id)
	if err != nil {
		plugin.Logger(ctx).Error("onepassword_item_medical_record.getItemMedicalRecord", "api_error", err)
		return nil, err
	}

	// restricting data based on the item category MEDICAL_RECORD
	var date, location, healthcareprofessional, patient, reason, medication, dosage, notes string
	if item.Category == "MEDICAL_RECORD" {
		for _, field := range item.Fields {
			if field.ID == "date" {
				date = field.Value
			} else if field.ID == "location" {
				location = field.Value
			} else if field.ID == "healthcareprofessional" {
				healthcareprofessional = field.Value
			} else if field.ID == "patient" {
				patient = field.Value
			} else if field.ID == "reason" {
				reason = field.Value
			} else if field.ID == "medication" {
				medication = field.Value
			} else if field.ID == "dosage" {
				dosage = field.Value
			} else if field.ID == "notes" && field.Label == "medication notes" {
				notes = field.Value
			}
		}

		return ItemMedicalRecord{date, dosage, healthcareprofessional, location, medication, notes, patient, reason, *item}, nil
	}

	return nil, nil
}
