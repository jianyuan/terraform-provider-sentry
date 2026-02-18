package provider

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/jianyuan/terraform-provider-sentry/internal/diagutils"
	"github.com/oapi-codegen/nullable"
)

func nullableToInt64(value nullable.Nullable[int64], key string, diags *diag.Diagnostics) types.Int64 {
	return nullableToFrameworkValue(
		value,
		key,
		diags,
		types.Int64Null,
		types.Int64Value,
	)
}

func nullableToString(value nullable.Nullable[string], key string, diags *diag.Diagnostics) types.String {
	return nullableToFrameworkValue(
		value,
		key,
		diags,
		types.StringNull,
		types.StringValue,
	)
}

func nullableToFrameworkValue[T any, V any](
	value nullable.Nullable[T],
	key string,
	diags *diag.Diagnostics,
	nullValue func() V,
	typedValue func(T) V,
) V {
	if !value.IsSpecified() || value.IsNull() {
		return nullValue()
	}

	v, err := value.Get()
	if err != nil {
		diags.Append(diagutils.NewFillError(fmt.Errorf("invalid %s value: %w", key, err)))
		return nullValue()
	}

	return typedValue(v)
}

func setNullableInt64(value types.Int64, target *nullable.Nullable[int64]) {
	setNullableValue(value.IsUnknown(), value.IsNull(), value.ValueInt64, target)
}

func setNullableString(value types.String, target *nullable.Nullable[string]) {
	setNullableValue(value.IsUnknown(), value.IsNull(), value.ValueString, target)
}

func setNullableValue[T any](isUnknown bool, isNull bool, value func() T, target *nullable.Nullable[T]) {
	if isUnknown {
		return
	}

	if isNull {
		target.SetNull()
		return
	}

	target.Set(value())
}
