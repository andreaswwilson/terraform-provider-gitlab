// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	gitlab "gitlab.com/gitlab-org/api/client-go"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &CurrentUserDataSource{}

func NewCurrentUserDataSource() datasource.DataSource {
	return &CurrentUserDataSource{}
}

type CurrentUserDataSource struct {
	client *gitlab.Client
}

// currentUserDataSourceModel describes the data source data model.
type currentUserDataSourceModel struct {
	Id       types.Int64  `tfsdk:"id"`
	Username types.String `tfsdk:"username"`
}

func (d *CurrentUserDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_current_user"
}

func (d *CurrentUserDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Current user data source",

		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				MarkdownDescription: "Gitlab user identifier",
				Computed:            true,
			},
			"username": schema.StringAttribute{
				MarkdownDescription: "Gitlab username",
				Computed:            true,
			},
		},
	}
}

func (d *CurrentUserDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*gitlab.Client)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *gitlab.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	d.client = client
}

func (d *CurrentUserDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data currentUserDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Les data fra gitlab ved hjelp av gitlab.client
	user, _, err := d.client.Users.CurrentUser()
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to read current user",
			err.Error(),
		)
		return
	}

	// Bruk data hentet fra gitlab og konverter til riktig type for terraform provider
	data.Id = types.Int64Value(int64(user.ID))
	data.Username = types.StringValue(user.Username)

	// Write logs using the tflog package
	// Documentation: https://terraform.io/plugin/log
	tflog.Trace(ctx, "read current user source")

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
