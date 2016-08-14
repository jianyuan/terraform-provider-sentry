provider "sentry" {
    token = "ba4c11e226cf4c0c914e1057dea1649bc7f9a5e993ec4e9aaa4de6cbd0944e34"
    base_url = "http://localhost:9000/api/"
}

resource "sentry_organization" "my_organization" {
    name = "My Organization"
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

resource "sentry_project" "worker_app" {
    organization = "${sentry_team.engineering.organization}"
    team = "${sentry_team.engineering.id}"
    name = "Worker App"
}
