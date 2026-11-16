package sentrytypes

// StringValueLike is satisfied by any string-based Terraform value, such as
// basetypes.StringValue and the custom string types in this package.
type StringValueLike interface {
	IsNull() bool
	IsUnknown() bool
	ValueString() string
}

// KnownStringPointer returns a pointer to the string value, or nil when the
// value is null or unknown. Use it for optional API fields that should be
// omitted from the request when not set.
func KnownStringPointer[T StringValueLike](v T) *string {
	if v.IsNull() || v.IsUnknown() {
		return nil
	}
	return new(v.ValueString())
}
