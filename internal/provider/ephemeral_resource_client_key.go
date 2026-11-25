package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-validators/ephemeralvalidator"
	"github.com/hashicorp/terraform-plugin-framework/ephemeral"
	"github.com/hashicorp/terraform-plugin-framework/ephemeral/schema"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"
	supertypes "github.com/orange-cloudavenue/terraform-plugin-framework-supertypes"
)

var _ ephemeral.EphemeralResource = &ClientKeyEphemeralResource{}
var _ ephemeral.EphemeralResourceWithConfigure = &ClientKeyEphemeralResource{}
var _ ephemeral.EphemeralResourceWithConfigValidators = &ClientKeyEphemeralResource{}

func NewClientKeyEphemeralResource() ephemeral.EphemeralResource {
	return &ClientKeyEphemeralResource{}
}

type ClientKeyEphemeralResource struct {
	baseEphemeralResource
}

func (e *ClientKeyEphemeralResource) Metadata(ctx context.Context, req ephemeral.MetadataRequest, resp *ephemeral.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_key"
}

func (e *ClientKeyEphemeralResource) Schema(ctx context.Context, req ephemeral.SchemaRequest, resp *ephemeral.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Retrieve a Project's Client Key without storing it in Terraform state. The DSN and secret can be passed to write-only attributes, for example to store them in a secrets manager.",
		Attributes: map[string]schema.Attribute{
			"organization": schema.StringAttribute{
				MarkdownDescription: "The organization the resource belongs to.",
				Required:            true,
			},
			"project": schema.StringAttribute{
				MarkdownDescription: "The project the resource belongs to.",
				Required:            true,
			},
			"id": schema.StringAttribute{
				MarkdownDescription: "The ID of the client key.",
				Optional:            true,
				Computed:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "The name of the client key.",
				Optional:            true,
				Computed:            true,
			},
			"first": schema.BoolAttribute{
				MarkdownDescription: "Return the first key of the returned keys.",
				Optional:            true,
			},
			"project_id": schema.StringAttribute{
				MarkdownDescription: "The ID of the project that the key belongs to.",
				Computed:            true,
			},
			"public": schema.StringAttribute{
				MarkdownDescription: "The public key.",
				Computed:            true,
			},
			"secret": schema.StringAttribute{
				MarkdownDescription: "The secret key.",
				Computed:            true,
				Sensitive:           true,
			},
			"rate_limit_window": schema.Int64Attribute{
				MarkdownDescription: "Length of time in seconds that will be considered when checking the rate limit.",
				Computed:            true,
			},
			"rate_limit_count": schema.Int64Attribute{
				MarkdownDescription: "Number of events that can be reported within the rate limit window.",
				Computed:            true,
			},
			"javascript_loader_script": schema.SingleNestedAttribute{
				MarkdownDescription: "The JavaScript loader script configuration.",
				Computed:            true,
				CustomType:          supertypes.NewSingleNestedObjectTypeOf[ClientKeyJavascriptLoaderScriptModel](ctx),
				Attributes: map[string]schema.Attribute{
					"browser_sdk_version": schema.StringAttribute{
						MarkdownDescription: "The version of the browser SDK to load.",
						Computed:            true,
					},
					"performance_monitoring_enabled": schema.BoolAttribute{
						MarkdownDescription: "Whether performance monitoring is enabled for this key.",
						Computed:            true,
					},
					"session_replay_enabled": schema.BoolAttribute{
						MarkdownDescription: "Whether session replay is enabled for this key.",
						Computed:            true,
					},
					"debug_enabled": schema.BoolAttribute{
						MarkdownDescription: "Whether debug bundles & logging are enabled for this key.",
						Computed:            true,
					},
				},
			},
			"dsn": schema.MapAttribute{
				MarkdownDescription: "This is a map of DSN values. The keys include `public`, `secret`, `csp`, `security`, `minidump`, `nel`, `unreal`, `cdn`, and `crons`.",
				ElementType:         types.StringType,
				CustomType:          supertypes.NewMapTypeOf[string](ctx),
				Computed:            true,
				Sensitive:           true,
			},
		},
	}
}

func (e *ClientKeyEphemeralResource) ConfigValidators(ctx context.Context) []ephemeral.ConfigValidator {
	return []ephemeral.ConfigValidator{
		ephemeralvalidator.Conflicting(
			path.MatchRoot("id"),
			path.MatchRoot("name"),
			path.MatchRoot("first"),
		),
	}
}

func (e *ClientKeyEphemeralResource) Open(ctx context.Context, req ephemeral.OpenRequest, resp *ephemeral.OpenResponse) {
	var data ClientKeyLookupModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	foundKey, diags := findClientKey(ctx, e.apiClient, data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(data.Fill(ctx, *foundKey)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.Result.Set(ctx, &data)...)
}
