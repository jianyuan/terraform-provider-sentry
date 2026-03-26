output "operand" {
  value = provider::sentry::op_status_code_check("greater_than", 199)
}
