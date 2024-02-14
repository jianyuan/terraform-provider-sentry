package sentry

import (
	"context"

	"github.com/getkevin/terraform-provider-sentry/sentry/lib"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceSentryPlugin() *schema.Resource {
	return &schema.Resource{
		Description: "Sentry Plugin resource.",

		CreateContext: resourceSentryPluginCreate,
		ReadContext:   resourceSentryPluginRead,
		UpdateContext: resourceSentryPluginUpdate,
		DeleteContext: resourceSentryPluginDelete,
		Importer: &schema.ResourceImporter{
			StateContext: importOrganizationProjectAndID,
		},

		Schema: map[string]*schema.Schema{
			"organization": {
				Description: "The slug of the organization the project belongs to.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"project": {
				Description: "The slug of the project to create the plugin for.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"plugin": {
				Description: "Plugin ID.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"config": {
				Description: "Plugin config.",
				Type:        schema.TypeMap,
				Optional:    true,
			},
		},
	}
}

func resourceSentryPluginCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*sentry.Client)

	plugin := d.Get("plugin").(string)
	org := d.Get("organization").(string)
	project := d.Get("project").(string)

	tflog.Debug(ctx, "Creating Sentry plugin", map[string]interface{}{
		"pluginName": plugin,
		"org":        org,
		"project":    project,
	})
	_, err := client.ProjectPlugins.Enable(ctx, org, project, plugin)
	if err != nil {
		return diag.FromErr(err)
	}
	tflog.Debug(ctx, "Created Sentry plugin", map[string]interface{}{
		"pluginName": plugin,
		"org":        org,
		"project":    project,
	})

	d.SetId(plugin)

	params := d.Get("config").(map[string]interface{})
	if _, _, err := client.ProjectPlugins.Update(ctx, org, project, plugin, params); err != nil {
		return diag.FromErr(err)
	}

	return resourceSentryPluginRead(ctx, d, meta)
}

func resourceSentryPluginRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*sentry.Client)

	id := d.Id()
	org := d.Get("organization").(string)
	project := d.Get("project").(string)

	tflog.Debug(ctx, "Reading Sentry plugin", map[string]interface{}{
		"pluginID": id,
		"org":      org,
		"project":  project,
	})
	plugin, resp, err := client.ProjectPlugins.Get(ctx, org, project, id)
	if found, err := checkClientGet(resp, err, d); !found {
		return diag.FromErr(err)
	}
	tflog.Debug(ctx, "Read Sentry plugin", map[string]interface{}{
		"pluginID": plugin.ID,
		"org":      org,
		"project":  project,
	})

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

	retErr := multierror.Append(
		d.Set("config", config),
	)
	return diag.FromErr(retErr.ErrorOrNil())
}

func resourceSentryPluginUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*sentry.Client)

	id := d.Id()
	org := d.Get("organization").(string)
	project := d.Get("project").(string)

	tflog.Debug(ctx, "Updating Sentry plugin", map[string]interface{}{
		"pluginID": id,
		"org":      org,
		"project":  project,
	})
	params := d.Get("config").(map[string]interface{})
	plugin, _, err := client.ProjectPlugins.Update(ctx, org, project, id, params)
	if err != nil {
		return diag.FromErr(err)
	}
	tflog.Debug(ctx, "Updated Sentry plugin", map[string]interface{}{
		"pluginID": plugin.ID,
		"org":      org,
		"project":  project,
	})

	return resourceSentryPluginRead(ctx, d, meta)
}

func resourceSentryPluginDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*sentry.Client)

	id := d.Id()
	org := d.Get("organization").(string)
	project := d.Get("project").(string)

	tflog.Debug(ctx, "Deleting Sentry plugin", map[string]interface{}{
		"pluginID": id,
		"org":      org,
		"project":  project,
	})
	_, err := client.ProjectPlugins.Disable(ctx, org, project, id)
	tflog.Debug(ctx, "Deleted Sentry plugin", map[string]interface{}{
		"pluginID": id,
		"org":      org,
		"project":  project,
	})

	return diag.FromErr(err)
}
