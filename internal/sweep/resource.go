package sweep

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/jianyuan/terraform-provider-sentry/internal/diagutils"
	"github.com/jianyuan/terraform-provider-sentry/internal/providerdata"
)

var _ Sweepable = (*sweepResource)(nil)

type sweepResource struct {
	factory    func() resource.Resource
	pd         *providerdata.ProviderData
	attributes map[string]any
}

func NewSweepResource(factory func() resource.Resource, pd *providerdata.ProviderData, attributes map[string]any) *sweepResource {
	return &sweepResource{
		factory:    factory,
		pd:         pd,
		attributes: attributes,
	}
}

func (sr *sweepResource) Delete(ctx context.Context) error {
	res := sr.factory()

	if res, ok := res.(resource.ResourceWithConfigure); ok {
		var configureResp resource.ConfigureResponse
		res.Configure(ctx, resource.ConfigureRequest{ProviderData: sr.pd}, &configureResp)
		if configureResp.Diagnostics.HasError() {
			return diagutils.DiagnosticsError(configureResp.Diagnostics)
		}
	}

	var schemaResp resource.SchemaResponse
	res.Schema(ctx, resource.SchemaRequest{}, &schemaResp)
	if schemaResp.Diagnostics.HasError() {
		return diagutils.DiagnosticsError(schemaResp.Diagnostics)
	}

	state := tfsdk.State{
		Raw:    tftypes.NewValue(schemaResp.Schema.Type().TerraformType(ctx), nil),
		Schema: schemaResp.Schema,
	}

	for k, v := range sr.attributes {
		if diags := state.SetAttribute(ctx, path.Root(k), v); diags.HasError() {
			return diagutils.DiagnosticsError(diags)
		}
	}

	log.Printf("[INFO] Deleting resource: %v", sr.attributes)
	var deleteResp resource.DeleteResponse
	res.Delete(ctx, resource.DeleteRequest{State: state}, &deleteResp)
	if deleteResp.Diagnostics.HasError() {
		return diagutils.DiagnosticsError(deleteResp.Diagnostics)
	}

	return nil
}
