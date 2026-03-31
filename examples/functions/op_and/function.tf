output "operand" {
  value = provider::sentry::op_and(
    provider::sentry::op_status_code_check("greater_than", 199),
    provider::sentry::op_status_code_check("less_than", 300),
  )
}
