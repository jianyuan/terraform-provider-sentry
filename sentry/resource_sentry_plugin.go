package sentry

import (
	"errors"
	"log"
	"strings"

	"github.com/hashicorp/terraform/helper/schema"
)

func resourceSentryPlugin() *schema.Resource {
	return &schema.Resource{
		Create: resourceSentryPluginCreate,
		Read:   resourceSentryPluginRead,
		Update: resourceSentryPluginUpdate,
		Delete: resourceSentryPluginDelete,
		Importer: &schema.ResourceImporter{
			State: resourceSentryPluginImporter,
		},

		Schema: map[string]*schema.Schema{
			"organization": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "The slug of the organization the project belongs to",
			},
			"project": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "The slug of the project to create the plugin for",
			},
			"plugin": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "The id of the plugin",
			},
			"config": &schema.Schema{
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "Plugin config",
			},
		},
	}
}

func resourceSentryPluginCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Client)

	plugin := d.Get("plugin").(string)
	org := d.Get("organization").(string)
	project := d.Get("project").(string)

	log.Printf("%v, %v, %v", plugin, org, project)

	if _, err := client.EnablePlugin(org, project, plugin); err != nil {
		return err
	}

	d.SetId(plugin)

	params := d.Get("config").(map[string]interface{})
	if _, _, err := client.UpdatePlugin(org, project, plugin, params); err != nil {
		return err
	}

	return resourceSentryPluginRead(d, meta)
}

func resourceSentryPluginRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Client)

	id := d.Id()
	org := d.Get("organization").(string)
	project := d.Get("project").(string)

	plugin, _, err := client.GetPlugin(org, project, id)
	if err != nil {
		d.SetId("")
		return nil
	}

	d.SetId(plugin.ID)

	pluginConfig := make(map[string]string)
	for _, entry := range plugin.Config {
		pluginConfig[entry.Name] = entry.Value
	}

	config := make(map[string]string)
	for k := range d.Get("config").(map[string]interface{}) {
		config[k] = pluginConfig[k]
	}

	d.Set("config", config)

	return nil
}

func resourceSentryPluginUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Client)

	id := d.Id()
	org := d.Get("organization").(string)
	project := d.Get("project").(string)

	params := d.Get("config").(map[string]interface{})
	if _, _, err := client.UpdatePlugin(org, project, id, params); err != nil {
		return err
	}

	return resourceSentryPluginRead(d, meta)
}

func resourceSentryPluginDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Client)

	id := d.Id()
	org := d.Get("organization").(string)
	project := d.Get("project").(string)

	_, err := client.DisablePlugin(org, project, id)
	return err
}

func resourceSentryPluginImporter(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	addrID := d.Id()

	log.Printf("[DEBUG] Importing key using ADDR ID %s", addrID)

	parts := strings.Split(addrID, "/")

	if len(parts) != 3 {
		return nil, errors.New("Project import requires an ADDR ID of the following schema org-slug/project-slug/plugin-id")
	}

	d.Set("organization", parts[0])
	d.Set("project", parts[1])
	d.SetId(parts[2])

	return []*schema.ResourceData{d}, nil
}
