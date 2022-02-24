package sentry

import (
	"context"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/jianyuan/go-sentry/sentry"
)

func resourceSentryPlugin() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceSentryPluginCreate,
		ReadContext:   resourceSentryPluginRead,
		UpdateContext: resourceSentryPluginUpdate,
		DeleteContext: resourceSentryPluginDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceSentryPluginImporter,
		},

		Schema: map[string]*schema.Schema{
			"organization": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The slug of the organization the project belongs to",
			},
			"project": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The slug of the project to create the plugin for",
			},
			"plugin": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Plugin ID",
			},
			"config": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "Plugin config",
			},
		},
	}
}

func resourceSentryPluginCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*sentry.Client)

	plugin := d.Get("plugin").(string)
	org := d.Get("organization").(string)
	project := d.Get("project").(string)

	tflog.Debug(ctx, "Creating Sentry plugin", "pluginName", plugin, "org", org, "project", project)
	_, err := client.ProjectPlugins.Enable(org, project, plugin)
	if err != nil {
		return diag.FromErr(err)
	}
	tflog.Debug(ctx, "Created Sentry plugin", "pluginName", plugin, "org", org, "project", project)

	d.SetId(plugin)

	params := d.Get("config").(map[string]interface{})
	if _, _, err := client.ProjectPlugins.Update(org, project, plugin, params); err != nil {
		return diag.FromErr(err)
	}

	return resourceSentryPluginRead(ctx, d, meta)
}

func resourceSentryPluginRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*sentry.Client)

	id := d.Id()
	org := d.Get("organization").(string)
	project := d.Get("project").(string)

	tflog.Debug(ctx, "Reading Sentry plugin", "pluginID", id, "org", org, "project", project)
	plugin, resp, err := client.ProjectPlugins.Get(org, project, id)
	if found, err := checkClientGet(resp, err, d); !found {
		return diag.FromErr(err)
	}
	tflog.Debug(ctx, "Read Sentry plugin", "pluginID", plugin.ID, "org", org, "project", project)

	d.SetId(plugin.ID)

	pluginConfig := make(map[string]string)
	for _, entry := range plugin.Config {
		if v, ok := entry.Value.(string); ok {
			pluginConfig[entry.Name] = v
		}
	}

	config := make(map[string]string)
	for k := range d.Get("config").(map[string]interface{}) {
		config[k] = pluginConfig[k]
	}

	d.Set("config", config)

	return nil
}

func resourceSentryPluginUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*sentry.Client)

	id := d.Id()
	org := d.Get("organization").(string)
	project := d.Get("project").(string)

	tflog.Debug(ctx, "Updating Sentry plugin", "pluginID", id, "org", org, "project", project)
	params := d.Get("config").(map[string]interface{})
	plugin, _, err := client.ProjectPlugins.Update(org, project, id, params)
	if err != nil {
		return diag.FromErr(err)
	}
	tflog.Debug(ctx, "Updated Sentry plugin", "pluginID", plugin.ID, "org", org, "project", project)

	return resourceSentryPluginRead(ctx, d, meta)
}

func resourceSentryPluginDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*sentry.Client)

	id := d.Id()
	org := d.Get("organization").(string)
	project := d.Get("project").(string)

	tflog.Debug(ctx, "Deleting Sentry plugin", "pluginID", id, "org", org, "project", project)
	_, err := client.ProjectPlugins.Disable(org, project, id)
	tflog.Debug(ctx, "Deleted Sentry plugin", "pluginID", id, "org", org, "project", project)

	return diag.FromErr(err)
}
