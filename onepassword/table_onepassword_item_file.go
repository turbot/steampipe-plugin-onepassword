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
		Get: &plugin.GetConfig{
			Hydrate:    getItemFile,
			KeyColumns: plugin.AllColumns([]string{"id", "item_id", "vault_id"}),
		},
		Columns: []*plugin.Column{
			{
				Name:        "id",
				Type:        proto.ColumnType_STRING,
				Description: "The UUID of the file.",
				Transform:   transform.FromField("ID"),
			},
			{
				Name:        "name",
				Type:        proto.ColumnType_STRING,
				Description: "The name of the file.",
			},
			{
				Name:        "item_id",
				Type:        proto.ColumnType_STRING,
				Description: "The UUID of the item.",
				Transform:   transform.FromQual("item_id"),
			},
			{
				Name:        "vault_id",
				Type:        proto.ColumnType_STRING,
				Description: "The UUID of the vault the file is in.",
			},
			{
				Name:        "content_path",
				Type:        proto.ColumnType_STRING,
				Description: "The path to download the contents of the file.",
			},
			{
				Name:        "size",
				Type:        proto.ColumnType_INT,
				Description: "The size of the file in bytes.",
			},
			{
				Name:        "content",
				Type:        proto.ColumnType_JSON,
				Description: "The Base64-encoded contents of the file, if inline_files is set to true.",
				Hydrate:     getFileContent,
				Transform:   transform.FromValue(),
			},
			{
				Name:        "section",
				Type:        proto.ColumnType_JSON,
				Description: "An object containing the UUID of a section in the item.",
			},

			/// Steampipe standard columns
			{
				Name:        "title",
				Description: "The title of the file.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Name"),
			},
		},
	}
}

type ItemFile struct {
	VaultId string
	onepassword.File
}

func listItemFiles(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	item_id := d.EqualsQualString("item_id")
	vault := h.Item.(onepassword.Vault)
	vault_id := d.EqualsQualString("vault_id")

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
		d.StreamListItem(ctx, ItemFile{vault.ID, file})

		// Context can be cancelled due to manual cancellation or the limit has been hit
		if d.RowsRemaining(ctx) == 0 {
			return nil, nil
		}
	}

	return nil, nil
}

func getItemFile(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	id := d.EqualsQualString("id")
	item_id := d.EqualsQualString("item_id")
	vault_id := d.EqualsQualString("vault_id")

	// Check if id, item_id or vault_id is empty
	if id == "" || item_id == "" || vault_id == "" {
		return nil, nil
	}

	client, err := getClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("onepassword_item_file.getItemFile", "connection_error", err)
		return nil, err
	}

	file, err := client.GetFile(id, item_id, vault_id)
	if err != nil {
		plugin.Logger(ctx).Error("onepassword_item_file.getItemFile", "api_error", err)
		return nil, err
	}

	return ItemFile{vault_id, *file}, nil
}

func getFileContent(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	file := h.Item.(ItemFile).File

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
