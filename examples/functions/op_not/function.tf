output "operand" {
  value = provider::sentry::op_not(
    provider::sentry::op_header_check(
      # Key
      "equals",
      provider::sentry::op_header_operand_literal("X-Header-Key"),
      # Value
      "equals",
      provider::sentry::op_header_operand_glob("pattern-*"),
    ),
  )
}
