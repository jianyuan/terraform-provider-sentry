provider "sentry" {
    token = "ba4c11e226cf4c0c914e1057dea1649bc7f9a5e993ec4e9aaa4de6cbd0944e34"
    base_url = "http://localhost:9000/api/"
}

resource "sentry_organization" "my_organization" {
    name = "My Organization"
    agree_terms = true
}

resource "sentry_team" "engineering" {
    organization = "${sentry_organization.my_organization.id}"
    name = "The Engineering Team"
}

resource "sentry_project" "web_app" {
    organization = "${sentry_team.engineering.organization}"
    team = "${sentry_team.engineering.id}"
    name = "Web App"
}

// Using the first parameter
data "sentry_key" "via_first" {
    organization = "${sentry_project.web_app.organization}"
    project = "${sentry_project.web_app.id}"
    first = true
}

// Using the name parameter
data "sentry_key" "via_name" {
    organization = "${sentry_project.web_app.organization}"
    project = "${sentry_project.web_app.id}"
    name = "Default"
}

output "sentry_key_dsn_secret" {
    value = "${data.sentry_key.via_name.dsn_secret}"
}
