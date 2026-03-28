output "operand" {
  value = provider::sentry::op_jsonpath(
    provider::sentry::op_jsonpath_operand_literal("value"),
    "equals",
    "$.status",
  )
}
