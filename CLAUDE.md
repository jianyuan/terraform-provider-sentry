# CLAUDE.md

## Terraform Plugin Framework types: use `supertypes`

Use `github.com/orange-cloudavenue/terraform-plugin-framework-supertypes`, imported as `supertypes`, for collection and nested object model fields and schema `CustomType`s. Don't use the untyped `types.Set`, `types.List`, `types.Map` or `types.Object` from `terraform-plugin-framework/types` in new code. Primitives (string, int64, bool, float64) stay on the plain `types` package.

```go
import supertypes "github.com/orange-cloudavenue/terraform-plugin-framework-supertypes"
```

The generic supertypes carry the element or nested model type, so you get typed `Get`/`Set` helpers instead of hand-written `ElementsAs`/`ObjectValueFrom` conversions and `attr.Type` maps.

| Schema shape | Model field type | Schema `CustomType` |
| --- | --- | --- |
| string | `types.String` | none |
| int64 | `types.Int64` | none |
| bool | `types.Bool` | none |
| float64 | `types.Float64` | none |
| list of primitives | `supertypes.ListValueOf[T]` | `supertypes.NewListTypeOf[T](ctx)` |
| set of primitives | `supertypes.SetValueOf[T]` | `supertypes.NewSetTypeOf[T](ctx)` |
| map of primitives | `supertypes.MapValueOf[T]` | `supertypes.NewMapTypeOf[T](ctx)` |
| list nested attribute | `supertypes.ListNestedObjectValueOf[M]` | `supertypes.NewListNestedObjectTypeOf[M](ctx)` |
| set nested attribute | `supertypes.SetNestedObjectValueOf[M]` | `supertypes.NewSetNestedObjectTypeOf[M](ctx)` |
| single nested attribute | `supertypes.SingleNestedObjectValueOf[M]` | `supertypes.NewSingleNestedObjectTypeOf[M](ctx)` |

`T` is a Go primitive (`string`, `int64`, `bool`) for list, set and map elements; `M` is the nested model struct.

Constructors to reach for:

- Collections: `supertypes.NewListValueOfSlice`, `NewSetValueOfSlice`, `NewSetValueOfNull`, `NewMapValueOfMap`, `NewMapValueOfNull`.
- Nested objects: `supertypes.NewListNestedObjectValueOfValueSlice`, `NewSetNestedObjectValueOfValueSlice`, `NewListNestedObjectValueOfNull`, `NewSingleNestedObjectValueOf`, `NewSingleNestedObjectValueOfNull`.

Read and write values with the typed helpers (`Get(ctx)`, `MustGet(ctx)`, `Set(ctx, v)`) rather than converting through `attr.Value`.

Primitives use `types.String`, `types.Int64`, `types.Bool` and `types.Float64` with no `CustomType`. The supertypes wrappers add nothing for them, and the float wrapper is actively broken: it inherits `Float64SemanticEquals` from basetypes without overriding it, so the type assertion rejects the wrapper on every plan (hashicorp/terraform-plugin-framework#786).

Some older resources (e.g. `resource_project.go`, `resource_client_key.go`, `resource_issue_alert.go`) still use `types.Set`/`types.List`/`types.Map`/`types.Object`. Don't copy that pattern. Migrate a field to supertypes when you're already changing it, as long as the change doesn't alter the schema or state.

## Custom types: `internal/sentrytypes`

Use a `sentrytypes` custom type when Sentry hands back a value in a different form from the one it was sent. Its semantic equality keeps the configured value in state, which avoids "Provider produced inconsistent result after apply" and diffs that never go away. Set the model field to the value type and the schema `CustomType` to the matching type.

| Value / `CustomType` | Use for | Used in |
| --- | --- | --- |
| `TrimmedString` / `TrimmedStringType{}` | Text that Sentry trims, such as ownership rules, CODEOWNERS files and other multi-line text often loaded with `file()`. Equal if equal after `strings.TrimSpace`. | `project_ownership.raw`, `project_codeowners.raw`, `project` and `uptime_monitor` fields |
| `LossyJson` / `LossyJsonType{IgnoreKeys: ...}` | JSON strings that Sentry rewrites. Equal if the JSON matches after numbers and their string forms are treated as the same value, and keys listed in `IgnoreKeys` (for example ones Sentry adds) are skipped. The type also checks that the value is valid JSON. | `issue_alert` conditions, filters and actions |
| `SlackChannel` / `SlackChannelType{}` | Slack channel names. Equal if they match once a leading `#` is removed. | `issue_alert`, `alert` Slack actions |
| `MsTeamsTeamId` / `MsTeamsTeamIdType{}` | Microsoft Teams team IDs, which Sentry replaces with an unrelated thread ID. Any two known values count as equal, so a read or apply that returns the thread ID keeps the configured value in state. | `alert` Microsoft Teams actions |
| `StringSet` / `StringSetType{}` | A Terraform set of strings that the API takes as one comma-separated string. Build it with `StringSetPointerValue`; read it back with `ValueString(ctx)` or `ValueStringPointer(ctx)`. | `issue_alert` |

`sentrytypes.KnownStringPointer(v)` is a helper, not a type. It returns `*string` for any string-based value, or `nil` when the value is null or unknown, for optional request fields that should be left out.

Don't add a `sentrytypes` type just to wrap a plain string; use `types.String`.

## Known checks: use `fwtypes.IsKnown`

To check that a value is neither null nor unknown, use `fwtypes.IsKnown(value)` from `github.com/jianyuan/terraform-plugin-framework-utils/fwtypes` rather than writing `!value.IsNull() && !value.IsUnknown()`.

```go
import "github.com/jianyuan/terraform-plugin-framework-utils/fwtypes"
```

## Model fillers

Fill a model from an API response with `func (m *Model) Fill(ctx context.Context, data apiclient.Model) (diags diag.Diagnostics)`, even when nothing in it can produce a diagnostic. Callers append the result and stop on error:

```go
resp.Diagnostics.Append(data.Fill(ctx, *httpResp.JSON200)...)
if resp.Diagnostics.HasError() {
	return
}
```

## Acceptance tests

Name acceptance tests `TestAcc<Resource>Resource_<test>` for resources and `TestAcc<DataSource>DataSource_<test>` for data sources, for example `TestAccProjectResource_basic`. `<test>` names the scenario, such as `basic`, `update` or `import`.

## API clients: prefer `apiclient`

Call Sentry through the generated `apiclient.ClientWithResponses` (`r.apiClient` in resources and data sources, `acctest.SharedApiClient` in tests) rather than the go-sentry `sentry.Client` (`r.client`, `acctest.SharedClient`). If an endpoint is missing, add it to `internal/apiclient/api.yaml` and run `go generate ./internal/apiclient`. Check `httpResp.StatusCode()` yourself, because the generated client does not turn non-2xx responses into errors.

Some older resources and tests still use `sentry.Client`. Don't copy that pattern; switch a call to `apiclient` when you're already changing it.

## Generated code

Many resources and data sources are generated by `internal/providergen` (TypeScript, run with bun). The supertypes mapping lives in `internal/providergen/go-types.ts`. Edit the generator and run `make generate` rather than hand-editing generated files.
