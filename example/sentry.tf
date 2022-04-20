terraform {
  required_providers {
    sentry = {
      source  = "hashicorp.com/edu/sentry"
      version = "0.7.51"
    }
  }
}

provider "sentry" {
  token    = "8583cbb12ed34be39dfc979ab1a3fed097f4c73688f742da859f9f17fde7bea4"
  base_url = "https://sentry.io/api"
}

# data "sentry_apm_rules" "all_rules" {
#   organization = "snappcar"
#   project      = "terraform-test"
# }

resource "sentry_apm_rule" "default" {
  organization     = "snappcar"
  project          = "terraform-test"
  name             = "test-from-terraform"
  # environment      = "testing"
  dataset          = "transactions"
  query            = "http.url:http://testservicer.com/stats"
  time_window       = 50.0
  aggregate        = "p50(transaction.duration)"
  threshold_type    = 0
  resolve_threshold = 100.0
  # triggers = [
  #   {
  #     actions          = [{}]
  #     alert_threshold   = 10000
  #     label            = "critical"
  #     resolve_threshold = 100.0
  #     threshold_type    = 0
  #   }
  # ]
  #When using Type.Set
  triggers {
      actions          = []
      alert_threshold   = 1000
      label            = "critical"
      resolve_threshold = 100.0
      threshold_type    = 0
  }

  triggers {
      actions          = []
      alert_threshold   = 500
      label            = "warning"
      resolve_threshold = 100.0
      threshold_type    = 0
  }

  #when using Type.List
  # triggers = [
  #   {
  #     # actions          = []
  #     alertThreshold   = 10000
  #     label            = "critical"
  #     resolveThreshold = 100.0
  #     thresholdType    = 0
  #   }
  # ]

  # triggers {
  #   # trigger {
  #     actions          = []
  #     alert_threshold   = 1000
  #     label            = "warning"
  #     resolve_threshold = 100.0
  #     threshold_type    = 0
  #   # }
  # }
    
  
#   projects = ["terraform-test"]
  owner = "team:52051"
}

# Returns all apm rules
# output "all_rules" {
#   value = data.sentry_apm_rules.all_rules.apm_rules
# }

# resource "sentry_organization" "my_organization" {
#     name = "My Organization"
#     agree_terms = true
# }

# resource "sentry_team" "engineering" {
#     organization = "${sentry_organization.my_organization.id}"
#     name = "The Engineering Team"
# }

# resource "sentry_project" "web_app" {
#     organization = "${sentry_team.engineering.organization}"
#     team = "${sentry_team.engineering.id}"
#     name = "Web App"
# }

# // Using the first parameter
# data "sentry_key" "via_first" {
#     organization = "${sentry_project.web_app.organization}"
#     project = "${sentry_project.web_app.id}"
#     first = true
# }

# // Using the name parameter
# data "sentry_key" "via_name" {
#     organization = "${sentry_project.web_app.organization}"
#     project = "${sentry_project.web_app.id}"
#     name = "Default"
# }

# output "sentry_key_dsn_secret" {
#     value = "${data.sentry_key.via_name.dsn_secret}"
# }
