package sentry

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/jianyuan/go-sentry/v2/sentry"
)

func resourceSentryOrganizationIntegrationPagerduty() *schema.Resource {
	return &schema.Resource{
		Description: "Sentry PagerDuty Organization Integration resource.",

		CreateContext: resourceSentryOrganizationIntegrationPagerdutyCreate,
		ReadContext:   resourceSentryOrganizationIntegrationPagerdutyRead,
		UpdateContext: resourceSentryOrganizationIntegrationPagerdutyUpdate,
		DeleteContext: resourceSentryOrganizationIntegrationPagerdutyDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"organization": {
				Description: "The slug of the Sentry organization this resource belongs to.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"integration_id": {
				Description: "The organization integration ID for PagerDuty.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"service_name": {
				Description: "The name of the PagerDuty service.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"integration_key": {
				Description: "The integration key from PagerDuty to associate with service_name.",
				Type:        schema.TypeString,
				Required:    true,
				Sensitive:   true,
			},
			"internal_id": {
				Description: "The internal ID for this PagerDuty service integration.",
				Type:        schema.TypeString,
				Computed:    true,
			},
		},
	}
}

func resourceSentryOrganizationIntegrationPagerdutyCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*sentry.Client)

	org := d.Get("organization").(string)
	integrationId := d.Get("integration_id").(string)
	serviceName := d.Get("service_name").(string)
	integrationKey := d.Get("integration_key").(string)

	tflog.Debug(ctx, "Creating PagerDuty service integration", map[string]interface{}{
		"org":            org,
		"integration_id": integrationId,
		"service_name":   serviceName,
	})
	orgIntegration, _, err := client.OrganizationIntegrations.Get(ctx, org, integrationId)
	if err != nil {
		return diag.FromErr(err)
	}

	serviceTable, err := extractServiceTable(orgIntegration)
	if err != nil {
		return diag.FromErr(err)
	}

	serviceTable = append(serviceTable, map[string]interface{}{
		"service":         serviceName,
		"integration_key": integrationKey,
		"id":              "",
	})
	updatedConfigData := sentry.UpdateConfigOrganizationIntegrationsParams{
		"service_table": serviceTable,
	}
	_, err = client.OrganizationIntegrations.UpdateConfig(ctx, org, integrationId, &updatedConfigData)
	if err != nil {
		return diag.FromErr(err)
	}

	orgIntegration, _, err = client.OrganizationIntegrations.Get(ctx, org, integrationId)
	if err != nil {
		return diag.FromErr(err)
	}

	serviceTable, err = extractServiceTable(orgIntegration)
	if err != nil {
		return diag.FromErr(err)
	}

	_, foundServiceRow := findServiceRowByNameAndKey(serviceTable, serviceName, integrationKey)
	if foundServiceRow == nil {
		return diag.Errorf("Unable to find PagerDuty service %s", serviceName)
	}

	d.SetId(buildThreePartID(org, integrationId, string(foundServiceRow["id"].(json.Number))))

	return resourceSentryOrganizationIntegrationPagerdutyRead(ctx, d, meta)
}

func resourceSentryOrganizationIntegrationPagerdutyRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*sentry.Client)
	org, integrationId, internalId, err := splitThreePartID(d.Id(), "organization-slug", "integration-id", "service-id")
	if err != nil {
		diag.FromErr(err)
	}

	tflog.Debug(ctx, "Reading Sentry PagerDuty Organization Integration", map[string]interface{}{
		"org":            org,
		"integration_id": integrationId,
		"internal_id":    internalId,
	})
	orgIntegration, _, err := client.OrganizationIntegrations.Get(ctx, org, integrationId)
	if err != nil {
		return diag.FromErr(err)
	}

	serviceTable, err := extractServiceTable(orgIntegration)
	if err != nil {
		return diag.FromErr(err)
	}

	_, foundServiceRow := findServiceRowById(serviceTable, internalId)
	if foundServiceRow != nil {
		internalId := string(foundServiceRow["id"].(json.Number))
		d.SetId(buildThreePartID(org, integrationId, internalId))
		retErr := multierror.Append(
			d.Set("organization", org),
			d.Set("integration_id", integrationId),
			d.Set("service_name", foundServiceRow["service"]),
			d.Set("integration_key", foundServiceRow["integration_key"]),
			d.Set("internal_id", internalId),
		)
		return diag.FromErr(retErr.ErrorOrNil())
	}

	tflog.Info(ctx, "Removing PagerDuty service from state because it no longer exists in Sentry", map[string]interface{}{"id": d.Id()})
	d.SetId("")
	return nil
}

func resourceSentryOrganizationIntegrationPagerdutyUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*sentry.Client)
	serviceName := d.Get("service_name")
	integrationKey := d.Get("integration_key")
	org, integrationId, internalId, err := splitThreePartID(d.Id(), "organization-slug", "integration-id", "service-id")
	if err != nil {
		diag.FromErr(err)
	}

	tflog.Debug(ctx, "Updating Sentry PagerDuty Organization Integration", map[string]interface{}{
		"org":            org,
		"integration_id": integrationId,
		"internal_id":    internalId,
	})
	orgIntegration, _, err := client.OrganizationIntegrations.Get(ctx, org, integrationId)
	if err != nil {
		return diag.FromErr(err)
	}

	serviceTable, err := extractServiceTable(orgIntegration)
	if err != nil {
		return diag.FromErr(err)
	}

	foundIndex, _ := findServiceRowById(serviceTable, internalId)
	if foundIndex >= 0 {
		serviceTable[foundIndex] = map[string]interface{}{
			"service":         serviceName,
			"integration_key": integrationKey,
			"id":              json.Number(internalId),
		}
	} else {
		return diag.Errorf("Unable to find PagerDuty service with id %s.", internalId)
	}

	updatedConfigData := sentry.UpdateConfigOrganizationIntegrationsParams{
		"service_table": serviceTable,
	}
	_, err = client.OrganizationIntegrations.UpdateConfig(ctx, org, integrationId, &updatedConfigData)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(buildThreePartID(org, integrationId, internalId))

	return resourceSentryOrganizationIntegrationPagerdutyRead(ctx, d, meta)
}

func resourceSentryOrganizationIntegrationPagerdutyDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*sentry.Client)
	org, integrationId, internalId, err := splitThreePartID(d.Id(), "organization-slug", "integration-id", "service-id")
	if err != nil {
		diag.FromErr(err)
	}

	tflog.Debug(ctx, "Deleting Sentry PagerDuty Organization Integration", map[string]interface{}{
		"org":            org,
		"integration_id": integrationId,
		"internal_id":    internalId,
	})
	orgIntegration, _, err := client.OrganizationIntegrations.Get(ctx, org, integrationId)
	if err != nil {
		return diag.FromErr(err)
	}

	serviceTable, err := extractServiceTable(orgIntegration)
	if err != nil {
		return diag.FromErr(err)
	}

	foundIndex, _ := findServiceRowById(serviceTable, internalId)
	if foundIndex < 0 {
		return diag.Errorf("Unable to find PagerDuty service with id %s.", internalId)
	}

	updatedServiceTable := append(serviceTable[:foundIndex], serviceTable[foundIndex+1:]...)
	updatedConfigData := sentry.UpdateConfigOrganizationIntegrationsParams{
		"service_table": updatedServiceTable,
	}
	_, err = client.OrganizationIntegrations.UpdateConfig(ctx, org, integrationId, &updatedConfigData)
	return diag.FromErr(err)
}

func extractServiceTable(orgIntegration *sentry.OrganizationIntegration) ([]interface{}, error) {
	configData := *orgIntegration.ConfigData
	serviceTable, ok := configData["service_table"].([]interface{})
	if !ok {
		return nil, fmt.Errorf("unable to find service_table in orgIntegration configData")
	}
	return serviceTable, nil
}

func findServiceRowById(serviceTable []interface{}, id string) (int, map[string]interface{}) {
	foundIndex := -1
	var foundServiceRow map[string]interface{}
	var serviceRow map[string]interface{}
	for index, row := range serviceTable {
		serviceRow = row.(map[string]interface{})
		if string(serviceRow["id"].(json.Number)) == id {
			foundServiceRow = serviceRow
			foundIndex = index
			break
		}
	}
	return foundIndex, foundServiceRow
}

func findServiceRowByNameAndKey(serviceTable []interface{}, serviceName string, integrationKey string) (int, map[string]interface{}) {
	foundIndex := -1
	var foundServiceRow map[string]interface{}
	var serviceRow map[string]interface{}
	for index, row := range serviceTable {
		serviceRow = row.(map[string]interface{})
		if serviceRow["service"] == serviceName && serviceRow["integration_key"] == integrationKey {
			foundServiceRow = serviceRow
			foundIndex = index
			break
		}
	}
	return foundIndex, foundServiceRow
}
