package resource

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/jianyuan/terraform-provider-sentry/internal/resourceid"
)

// ImportState1PartPassthrough imports a single-part ID (e.g. "my-org" or "12345").
func ImportState1PartPassthrough(
	attrPathA string,
) func(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	return func(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
		id := strings.TrimSpace(req.ID)
		if id == "" {
			resp.Diagnostics.AddError(
				"Invalid Resource Import ID",
				"Import ID cannot be empty",
			)
			return
		}

		resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root(attrPathA), id)...)
	}
}

// ImportState2PartPath imports a slash-separated two-part ID (e.g. "my-org/12345").
func ImportState2PartPath(
	attrPathA, attrPathB string,
) func(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	return func(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
		parts := strings.Split(strings.TrimSpace(req.ID), "/")
		if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
			resp.Diagnostics.AddError(
				"Invalid Resource Import ID",
				fmt.Sprintf("Unexpected import ID format %q, expected %s/%s", req.ID, attrPathA, attrPathB),
			)
			return
		}

		resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root(attrPathA), parts[0])...)
		resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root(attrPathB), parts[1])...)
	}
}

// ImportState3PartPath imports a slash-separated three-part ID (e.g. "my-org/web-app/12345").
func ImportState3PartPath(
	attrPathA, attrPathB, attrPathC string,
) func(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	return func(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
		parts := strings.Split(strings.TrimSpace(req.ID), "/")
		if len(parts) != 3 || parts[0] == "" || parts[1] == "" || parts[2] == "" {
			resp.Diagnostics.AddError(
				"Invalid Resource Import ID",
				fmt.Sprintf("Unexpected import ID format %q, expected %s/%s/%s", req.ID, attrPathA, attrPathB, attrPathC),
			)
			return
		}

		resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root(attrPathA), parts[0])...)
		resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root(attrPathB), parts[1])...)
		resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root(attrPathC), parts[2])...)
	}
}

// ImportState4PartPath imports a slash-separated four-part ID (e.g. "my-org/web-app/12345/events/").
func ImportState4PartPath(
	attrPathA, attrPathB, attrPathC, attrPathD string,
) func(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	return func(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
		parts := strings.Split(strings.TrimSpace(req.ID), "/")
		if len(parts) != 4 || parts[0] == "" || parts[1] == "" || parts[2] == "" || parts[3] == "" {
			resp.Diagnostics.AddError(
				"Invalid Resource Import ID",
				fmt.Sprintf("Unexpected import ID format %q, expected %s/%s/%s/%s", req.ID, attrPathA, attrPathB, attrPathC, attrPathD),
			)
			return
		}

		resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root(attrPathA), parts[0])...)
		resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root(attrPathB), parts[1])...)
		resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root(attrPathC), parts[2])...)
		resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root(attrPathD), parts[3])...)
	}
}

// ImportState1Part returns a resource.Resource ImportState handler for 1-part identifiers.
// Example URL: "https://{organization}.sentry.io/settings/"
// Example Key: "my-org"
func ImportState1Part(
	urlTemplate string,
	labelA, attrPathA string,
) func(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	return func(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
		valA, err := resourceid.Parse(req.ID, urlTemplate, labelA)
		if err != nil {
			resp.Diagnostics.AddError(
				"Invalid Resource Import ID",
				fmt.Sprintf("Could not parse import ID %q: %s", req.ID, err.Error()),
			)
			return
		}

		// Set the extracted value in Terraform state schema
		resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root(attrPathA), valA)...)
	}
}

// ImportState2Part returns a resource.Resource ImportState handler for 2-part identifiers.
// Example URL: "https://{organization}.sentry.io/monitors/{id}/"
// Example Key: "my-org/12345"
func ImportState2Part(
	urlTemplate string,
	labelA, attrPathA string,
	labelB, attrPathB string,
) func(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	return func(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
		valA, valB, err := resourceid.Split2(req.ID, urlTemplate, labelA, labelB)
		if err != nil {
			resp.Diagnostics.AddError(
				"Invalid Resource Import ID",
				fmt.Sprintf("Could not parse import ID %q: %s", req.ID, err.Error()),
			)
			return
		}

		// Set extracted values in Terraform state schema
		resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root(attrPathA), valA)...)
		resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root(attrPathB), valB)...)
	}
}

// ImportState3Part returns a resource.Resource ImportState handler for 3-part identifiers.
func ImportState3Part(
	urlTemplate string,
	labelA, attrPathA string,
	labelB, attrPathB string,
	labelC, attrPathC string,
) func(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	return func(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
		valA, valB, valC, err := resourceid.Split3(req.ID, urlTemplate, labelA, labelB, labelC)
		if err != nil {
			resp.Diagnostics.AddError(
				"Invalid Resource Import ID",
				fmt.Sprintf("Could not parse import ID %q: %s", req.ID, err.Error()),
			)
			return
		}

		resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root(attrPathA), valA)...)
		resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root(attrPathB), valB)...)
		resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root(attrPathC), valC)...)
	}
}

// ImportState4Part returns a resource.Resource ImportState handler for 4-part identifiers.
func ImportState4Part(
	urlTemplate string,
	labelA, attrPathA string,
	labelB, attrPathB string,
	labelC, attrPathC string,
	labelD, attrPathD string,
) func(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	return func(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
		valA, valB, valC, valD, err := resourceid.Split4(req.ID, urlTemplate, labelA, labelB, labelC, labelD)
		if err != nil {
			resp.Diagnostics.AddError(
				"Invalid Resource Import ID",
				fmt.Sprintf("Could not parse import ID %q: %s", req.ID, err.Error()),
			)
			return
		}

		resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root(attrPathA), valA)...)
		resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root(attrPathB), valB)...)
		resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root(attrPathC), valC)...)
		resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root(attrPathD), valD)...)
	}
}
