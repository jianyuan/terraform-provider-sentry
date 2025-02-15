package provider

import (
	schemaD "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	schemaR "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	superschema "github.com/orange-cloudavenue/terraform-plugin-framework-superschema"
)

func clientKeySchema() superschema.Schema {
	return superschema.Schema{
		Resource: superschema.SchemaDetails{
			MarkdownDescription: "Return a client key bound to a project.",
		},
		DataSource: superschema.SchemaDetails{
			MarkdownDescription: "Retrieve a Project's Client Key.",
		},
		Attributes: map[string]superschema.Attribute{
			"id": superschema.StringAttribute{
				Common: &schemaR.StringAttribute{
					MarkdownDescription: "The ID of this resource.",
				},
				Resource: &schemaR.StringAttribute{
					Computed: true,
					PlanModifiers: []planmodifier.String{
						stringplanmodifier.UseStateForUnknown(),
					},
				},
				DataSource: &schemaD.StringAttribute{
					Optional: true,
				},
			},
			"organization": superschema.StringAttribute{
				Common: &schemaR.StringAttribute{
					MarkdownDescription: "The organization of this resource.",
					Required:            true,
				},
				Resource: &schemaR.StringAttribute{
					PlanModifiers: []planmodifier.String{
						stringplanmodifier.UseStateForUnknown(),
					},
				},
			},
			"project": superschema.StringAttribute{
				Common: &schemaR.StringAttribute{
					MarkdownDescription: "The project of this resource.",
					Required:            true,
				},
				Resource: &schemaR.StringAttribute{
					PlanModifiers: []planmodifier.String{
						stringplanmodifier.UseStateForUnknown(),
					},
				},
			},
			"name": superschema.StringAttribute{
				Common: &schemaR.StringAttribute{
					MarkdownDescription: "The name of the client key.",
				},
				Resource: &schemaR.StringAttribute{
					Required: true,
				},
				DataSource: &schemaD.StringAttribute{
					Optional: true,
				},
			},
			"rate_limit_window": superschema.Int64Attribute{
				Common: &schemaR.Int64Attribute{
					MarkdownDescription: "Length of time in seconds that will be considered when checking the rate limit.",
				},
				Resource: &schemaR.Int64Attribute{
					Optional: true,
					Computed: true,
					PlanModifiers: []planmodifier.Int64{
						int64planmodifier.UseStateForUnknown(),
					},
				},
				DataSource: &schemaD.Int64Attribute{
					Computed: true,
				},
			},
			"rate_limit_count": superschema.Int64Attribute{
				Common: &schemaR.Int64Attribute{
					MarkdownDescription: "Number of events that can be reported within the rate limit window.",
				},
				Resource: &schemaR.Int64Attribute{
					Optional: true,
					Computed: true,
					PlanModifiers: []planmodifier.Int64{
						int64planmodifier.UseStateForUnknown(),
					},
				},
				DataSource: &schemaD.Int64Attribute{
					Computed: true,
				},
			},
			"javascript_loader_script": superschema.SuperSingleNestedAttributeOf[ClientKeyJavascriptLoaderScriptModel]{
				Common: &schemaR.SingleNestedAttribute{
					MarkdownDescription: "The JavaScript loader script configuration.",
				},
				Resource: &schemaR.SingleNestedAttribute{
					Optional: true,
					Computed: true,
					PlanModifiers: []planmodifier.Object{
						objectplanmodifier.UseStateForUnknown(),
					},
				},
				DataSource: &schemaD.SingleNestedAttribute{
					Computed: true,
				},
				Attributes: map[string]superschema.Attribute{
					"browser_sdk_version": superschema.StringAttribute{
						Common: &schemaR.StringAttribute{
							MarkdownDescription: "The version of the browser SDK to load.",
						},
						Resource: &schemaR.StringAttribute{
							Optional: true,
							Computed: true,
							PlanModifiers: []planmodifier.String{
								stringplanmodifier.UseStateForUnknown(),
							},
						},
						DataSource: &schemaD.StringAttribute{
							Computed: true,
						},
					},
					"performance_monitoring_enabled": superschema.BoolAttribute{
						Common: &schemaR.BoolAttribute{
							MarkdownDescription: "Whether performance monitoring is enabled for this key.",
						},
						Resource: &schemaR.BoolAttribute{
							Optional: true,
							Computed: true,
							PlanModifiers: []planmodifier.Bool{
								boolplanmodifier.UseStateForUnknown(),
							},
						},
						DataSource: &schemaD.BoolAttribute{
							Computed: true,
						},
					},
					"session_replay_enabled": superschema.BoolAttribute{
						Common: &schemaR.BoolAttribute{
							MarkdownDescription: "Whether session replay is enabled for this key.",
						},
						Resource: &schemaR.BoolAttribute{
							Optional: true,
							Computed: true,
							PlanModifiers: []planmodifier.Bool{
								boolplanmodifier.UseStateForUnknown(),
							},
						},
						DataSource: &schemaD.BoolAttribute{
							Computed: true,
						},
					},
					"debug_enabled": superschema.BoolAttribute{
						Common: &schemaR.BoolAttribute{
							MarkdownDescription: "Whether debug bundles & logging are enabled for this key.",
						},
						Resource: &schemaR.BoolAttribute{
							Optional: true,
							Computed: true,
							PlanModifiers: []planmodifier.Bool{
								boolplanmodifier.UseStateForUnknown(),
							},
						},
						DataSource: &schemaD.BoolAttribute{
							Computed: true,
						},
					},
				},
			},
			"project_id": superschema.StringAttribute{
				Common: &schemaR.StringAttribute{
					MarkdownDescription: "The ID of the project that the key belongs to.",
					Computed:            true,
				},
			},
			"public": superschema.StringAttribute{
				Common: &schemaR.StringAttribute{
					MarkdownDescription: "The public key.",
					Computed:            true,
				},
			},
			"secret": superschema.StringAttribute{
				Common: &schemaR.StringAttribute{
					MarkdownDescription: "The secret key.",
					Computed:            true,
					Sensitive:           true,
				},
			},
			"dsn": superschema.MapAttribute{
				Common: &schemaR.MapAttribute{
					MarkdownDescription: "This is a map of DSN values. The keys include `public`, `secret`, `csp`, `security`, `minidump`, `nel`, `unreal`, `cdn`, and `crons`.",
					ElementType:         types.StringType,
					Computed:            true,
					Sensitive:           true,
				},
			},
			"dsn_public": superschema.StringAttribute{
				Deprecated: &superschema.Deprecated{
					DeprecationMessage:                "This field is deprecated and will be removed in a future version. Use `dsn[\"public\"]` instead.",
					ComputeMarkdownDeprecationMessage: true,
					Removed:                           true,
					TargetAttributeName:               "dsn[\"public\"]",
				},
				Common: &schemaR.StringAttribute{
					MarkdownDescription: "The DSN tells the SDK where to send the events to.",
					Computed:            true,
				},
			},
			"dsn_secret": superschema.StringAttribute{
				Deprecated: &superschema.Deprecated{
					DeprecationMessage:                "This field is deprecated and will be removed in a future version. Use `dsn[\"secret\"]` instead.",
					ComputeMarkdownDeprecationMessage: true,
					Removed:                           true,
					TargetAttributeName:               "dsn[\"secret\"]",
				},
				Common: &schemaR.StringAttribute{
					MarkdownDescription: "Deprecated DSN includes a secret which is no longer required by newer SDK versions. If you are unsure which to use, follow installation instructions for your language.",
					Computed:            true,
					Sensitive:           true,
				},
			},
			"dsn_csp": superschema.StringAttribute{
				Deprecated: &superschema.Deprecated{
					DeprecationMessage:                "This field is deprecated and will be removed in a future version. Use `dsn[\"csp\"]` instead.",
					ComputeMarkdownDeprecationMessage: true,
					Removed:                           true,
					TargetAttributeName:               "dsn[\"csp\"]",
				},
				Common: &schemaR.StringAttribute{
					MarkdownDescription: "Security header endpoint for features like CSP and Expect-CT reports.",
					Computed:            true,
				},
			},
			"first": superschema.BoolAttribute{
				DataSource: &schemaD.BoolAttribute{
					MarkdownDescription: "Return the first key of the returned keys.",
					Optional:            true,
				},
			},
		},
	}
}
