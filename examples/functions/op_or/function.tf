output "operand" {
  value = provider::sentry::op_or(
    provider::sentry::op_status_code_check("equals", 404),
    provider::sentry::op_status_code_check("equals", 500),
  )
}
