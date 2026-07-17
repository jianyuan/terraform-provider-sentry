package tfutils

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/jianyuan/terraform-provider-sentry/internal/diagutils"
)

func ImportStateTwoPart(ctx context.Context, part1 string, part2 string, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	part1Value, part2Value, err := SplitTwoPartId(req.ID, part1, part2)
	if err != nil {
		resp.Diagnostics.Append(diagutils.NewImportError(err))
		return
	}
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root(part1), part1Value)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root(part2), part2Value)...)
}

func ImportStateTwoPartId(ctx context.Context, part1 string, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	ImportStateTwoPart(ctx, part1, "id", req, resp)
}

func ImportStateThreePartId(ctx context.Context, part1 string, part2 string, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	part1Value, part2Value, id, err := SplitThreePartId(req.ID, part1, part2, "id")
	if err != nil {
		resp.Diagnostics.Append(diagutils.NewImportError(err))
		return
	}
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root(part1), part1Value)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root(part2), part2Value)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), id)...)
}

func ImportStateFourPartId(ctx context.Context, part1 string, part2 string, part3 string, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	part1Value, part2Value, part3Value, id, err := SplitFourPartId(req.ID, part1, part2, part3, "id")
	if err != nil {
		resp.Diagnostics.Append(diagutils.NewImportError(err))
		return
	}
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root(part1), part1Value)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root(part2), part2Value)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root(part3), part3Value)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), id)...)
}
