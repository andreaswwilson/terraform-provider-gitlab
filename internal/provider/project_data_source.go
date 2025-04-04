// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	gitlab "gitlab.com/gitlab-org/api/client-go"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &ProjectDataSource{}

func NewProjectDataSource() datasource.DataSource {
	return &ProjectDataSource{}
}

// ProjectDataSource defines the data source implementation.
type ProjectDataSource struct {
	client *gitlab.Client
}

// ProjectDataSourceModel describes the data source data model.
type ProjectDataSourceModel struct {
	Description       types.String `tfsdk:"description"`
	PathWithNameSpace types.String `tfsdk:"path_with_namespace"`
}

func (d *ProjectDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_project"
}

func (d *ProjectDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Read gitlab project",

		Attributes: map[string]schema.Attribute{
			"path_with_namespace": schema.StringAttribute{
				MarkdownDescription: "Path with namespace of project",
				Required:            true,
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "Project description",
				Computed:            true,
			},
		},
	}
}

func (d *ProjectDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *ProjectDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data ProjectDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}
	project, _, err := d.client.Projects.GetProject(data.PathWithNameSpace.ValueString(), &gitlab.GetProjectOptions{})
	if err != nil {
		resp.Diagnostics.AddError("Unable to read project", err.Error())
		return
	}

	data.Description = types.StringValue(project.Description)

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
