package tfutils

// Float64 variants of orange-cloudavenue's Null/RequireIfAttributeIsOneOf,
// which the library only ships for Int64.

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

// NullIfAttributeIsOneOfFloat64 ensures the float64 attribute is null when the
// referenced attribute has one of the given values.
func NullIfAttributeIsOneOfFloat64(p path.Expression, exceptedValues []attr.Value) validator.Float64 {
	return nullIfOneOfFloat64{path: p, exceptedValues: exceptedValues}
}

type nullIfOneOfFloat64 struct {
	path           path.Expression
	exceptedValues []attr.Value
}

func (v nullIfOneOfFloat64) Description(_ context.Context) string {
	return fmt.Sprintf("must be null when %s is one of %v", v.path, v.exceptedValues)
}

func (v nullIfOneOfFloat64) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

func (v nullIfOneOfFloat64) ValidateFloat64(ctx context.Context, req validator.Float64Request, resp *validator.Float64Response) {
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}
	paths, diags := req.Config.PathMatches(ctx, req.PathExpression.Merge(v.path))
	resp.Diagnostics.Append(diags...)
	if diags.HasError() {
		return
	}
	for _, p := range paths {
		var other attr.Value
		d := req.Config.GetAttribute(ctx, p, &other)
		resp.Diagnostics.Append(d...)
		if d.HasError() {
			return
		}
		if other.IsNull() || other.IsUnknown() {
			continue
		}
		for _, ex := range v.exceptedValues {
			if other.Equal(ex) {
				resp.Diagnostics.AddAttributeError(
					req.Path,
					"Invalid attribute combination",
					fmt.Sprintf("Attribute must be null when %s is %s.", p, ex),
				)
				return
			}
		}
	}
}

// RequireIfAttributeIsOneOfFloat64 ensures the float64 attribute is set when the
// referenced attribute has one of the given values.
func RequireIfAttributeIsOneOfFloat64(p path.Expression, requiredValues []attr.Value) validator.Float64 {
	return requireIfOneOfFloat64{path: p, requiredValues: requiredValues}
}

type requireIfOneOfFloat64 struct {
	path           path.Expression
	requiredValues []attr.Value
}

func (v requireIfOneOfFloat64) Description(_ context.Context) string {
	return fmt.Sprintf("must be set when %s is one of %v", v.path, v.requiredValues)
}

func (v requireIfOneOfFloat64) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

func (v requireIfOneOfFloat64) ValidateFloat64(ctx context.Context, req validator.Float64Request, resp *validator.Float64Response) {
	if !req.ConfigValue.IsNull() {
		return
	}
	paths, ds := req.Config.PathMatches(ctx, req.PathExpression.Merge(v.path))
	resp.Diagnostics.Append(ds...)
	if ds.HasError() {
		return
	}
	for _, p := range paths {
		var other attr.Value
		d := req.Config.GetAttribute(ctx, p, &other)
		resp.Diagnostics.Append(d...)
		if d.HasError() {
			return
		}
		if other.IsNull() || other.IsUnknown() {
			continue
		}
		for _, want := range v.requiredValues {
			if other.Equal(want) {
				resp.Diagnostics.AddAttributeError(
					req.Path,
					"Missing required attribute",
					fmt.Sprintf("Attribute must be set when %s is %s.", p, want),
				)
				return
			}
		}
	}
}

var (
	_ validator.Float64 = nullIfOneOfFloat64{}
	_ validator.Float64 = requireIfOneOfFloat64{}
)
