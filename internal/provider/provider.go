// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"os"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/ephemeral"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	gitlab "gitlab.com/gitlab-org/api/client-go"
)

// Ensure ScaffoldingProvider satisfies various provider interfaces.
var (
	_ provider.Provider                       = &GitlabProvider{}
	_ provider.ProviderWithFunctions          = &GitlabProvider{}
	_ provider.ProviderWithEphemeralResources = &GitlabProvider{}
)

// GitlabProvider defines the provider implementation.
type GitlabProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

// GitlabProviderModel describes the provider data model.
type GitlabProviderModel struct {
	Token types.String `tfsdk:"token"`
}

func (p *GitlabProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "gitlab"
	resp.Version = p.version
}

func (p *GitlabProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"token": schema.StringAttribute{
				Optional:  true,
				Sensitive: true,
			},
		},
	}
}

func (p *GitlabProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var data GitlabProviderModel
	// les data fra konfigurasjon
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if data.Token.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("token"), "Unkown token", "The provider cannot create the client as there is an unkown configuration value for the token. "+
				"Either target apply the source of the value first, set the value statically in the configration, or use the GITLAB_TOKEN environmental variable.",
		)
	}
	if resp.Diagnostics.HasError() {
		return
	}

	// Hent data fra miljøvariabel
	token := os.Getenv("GITLAB_TOKEN")

	// Sjekk om token har blitt satt i konfigurasjon, i så fall bruker vi denne
	if !data.Token.IsNull() {
		token = data.Token.ValueString()
	}
	if token == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("token"),
			"Missing Gitlab API Token",
			"The provider cannot create the API client as there is a missing or empty value for the API token. "+
				"Set the token value in the configuration or use the GITLAB_TOKEN environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}
	if resp.Diagnostics.HasError() {
		return
	}

	// Sett opp gitlab-client
	client, err := gitlab.NewClient(token)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Create Gitlab Client",
			"An unexpected error occurred when creating the gitlab client. "+
				"Gitlab Client Error: "+err.Error(),
		)
		return
	}
	resp.DataSourceData = client
	resp.ResourceData = client
}

func (p *GitlabProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		// NewExampleResource,
	}
}

func (p *GitlabProvider) EphemeralResources(ctx context.Context) []func() ephemeral.EphemeralResource {
	return []func() ephemeral.EphemeralResource{
		// NewExampleEphemeralResource,
	}
}

func (p *GitlabProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewCurrentUserDataSource,
	}
}

func (p *GitlabProvider) Functions(ctx context.Context) []func() function.Function {
	return []func() function.Function{
		// NewExampleFunction,
	}
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &GitlabProvider{
			version: version,
		}
	}
}
