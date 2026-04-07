package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
	"github.com/jianyuan/terraform-provider-sentry/internal/acctest"
)

func TestAccMetricMonitorDataSource_threshold(t *testing.T) {
	teamName := acctest.RandomWithPrefix("tf-team")
	projectName := acctest.RandomWithPrefix("tf-project")
	monitorName := acctest.RandomWithPrefix("tf-metric-monitor")
	rn := "data.sentry_metric_monitor.test"

	checks := []statecheck.StateCheck{
		statecheck.ExpectKnownValue(rn, tfjsonpath.New("id"), knownvalue.NotNull()),
		statecheck.ExpectKnownValue(rn, tfjsonpath.New("organization"), knownvalue.StringExact(acctest.TestOrganization)),
		statecheck.ExpectKnownValue(rn, tfjsonpath.New("project_id"), knownvalue.NotNull()),
		statecheck.ExpectKnownValue(rn, tfjsonpath.New("owner"), knownvalue.ObjectExact(map[string]knownvalue.Check{
			"user_id": knownvalue.Null(),
			"team_id": knownvalue.NotNull(),
		})),
	}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccMetricMonitorDataSourceConfig(teamName, projectName, monitorName, `
					aggregate = "count()"
					dataset = "events"
					event_types = ["default", "error"]
					query = "is:unresolved"
					query_type = "error"
					time_window_seconds = 3600

					condition_group = {
						conditions = [
							{
								type = "gt"
								comparison = 100
								condition_result = 75
							},
							{
								type = "lte"
								comparison = 50
								condition_result = 0
							},
						]
					}

					issue_detection = {
						type = "static"
					}
				`),
				ConfigStateChecks: append(
					checks,
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("enabled"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("name"), knownvalue.StringExact(monitorName)),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("aggregate"), knownvalue.StringExact("count()")),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("dataset"), knownvalue.StringExact("events")),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("event_types"), knownvalue.SetExact([]knownvalue.Check{
						knownvalue.StringExact("default"),
						knownvalue.StringExact("error"),
					})),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("query"), knownvalue.StringExact("is:unresolved")),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("query_type"), knownvalue.StringExact("error")),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("time_window_seconds"), knownvalue.Int64Exact(3600)),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("condition_group"), knownvalue.ObjectExact(map[string]knownvalue.Check{
						"logic_type": knownvalue.StringExact("any"),
						"conditions": knownvalue.ListExact([]knownvalue.Check{
							knownvalue.ObjectPartial(map[string]knownvalue.Check{
								"type":             knownvalue.StringExact("gt"),
								"comparison":       knownvalue.Int64Exact(100),
								"condition_result": knownvalue.Int64Exact(75),
							}),
							knownvalue.ObjectPartial(map[string]knownvalue.Check{
								"type":             knownvalue.StringExact("lte"),
								"comparison":       knownvalue.Int64Exact(50),
								"condition_result": knownvalue.Int64Exact(0),
							}),
						}),
					})),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("issue_detection"), knownvalue.ObjectExact(map[string]knownvalue.Check{
						"type":             knownvalue.StringExact("static"),
						"comparison_delta": knownvalue.Null(),
					})),
				),
			},
			{
				Config: testAccMetricMonitorDataSourceConfig(teamName, projectName, monitorName, `
					aggregate = "count()"
					dataset = "events"
					event_types = ["default", "error"]
					query = "is:unresolved"
					time_window_seconds = 3600

					condition_group = {
						conditions = [
							{
								type = "gt"
								comparison = 100
								condition_result = 75
							},
							{
								type = "lte"
								comparison = 50
								condition_result = 0
							},
						]
					}

					issue_detection = {
						type = "static"
					}
				`),
				ConfigStateChecks: append(
					checks,
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("enabled"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("name"), knownvalue.StringExact(monitorName)),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("aggregate"), knownvalue.StringExact("count()")),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("dataset"), knownvalue.StringExact("events")),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("event_types"), knownvalue.SetExact([]knownvalue.Check{
						knownvalue.StringExact("default"),
						knownvalue.StringExact("error"),
					})),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("query"), knownvalue.StringExact("is:unresolved")),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("query_type"), knownvalue.StringExact("error")),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("time_window_seconds"), knownvalue.Int64Exact(3600)),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("condition_group"), knownvalue.ObjectExact(map[string]knownvalue.Check{
						"logic_type": knownvalue.StringExact("any"),
						"conditions": knownvalue.ListExact([]knownvalue.Check{
							knownvalue.ObjectPartial(map[string]knownvalue.Check{
								"type":             knownvalue.StringExact("gt"),
								"comparison":       knownvalue.Int64Exact(100),
								"condition_result": knownvalue.Int64Exact(75),
							}),
							knownvalue.ObjectPartial(map[string]knownvalue.Check{
								"type":             knownvalue.StringExact("lte"),
								"comparison":       knownvalue.Int64Exact(50),
								"condition_result": knownvalue.Int64Exact(0),
							}),
						}),
					})),
				),
			},
			{
				Config: testAccMetricMonitorDataSourceConfig(teamName, projectName, monitorName+"-updated", `
					aggregate = "count()"
					dataset = "events"
					event_types = ["default", "error"]

					condition_group = {
						conditions = [
							{
								type = "gt"
								comparison = 100
								condition_result = 75
							},
							{
								type = "lte"
								comparison = 50
								condition_result = 0
							},
						]
					}

					issue_detection = {
						type = "static"
					}
				`),
				ConfigStateChecks: append(
					checks,
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("enabled"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("name"), knownvalue.StringExact(monitorName+"-updated")),
				),
			},
			{
				Config: testAccMetricMonitorDataSourceConfig(teamName, projectName, monitorName+"-updated", `
					enabled = false

					aggregate = "count()"
					dataset = "events"
					event_types = ["default", "error"]

					condition_group = {
						conditions = [
							{
								type = "gt"
								comparison = 100
								condition_result = 75
							},
							{
								type = "lte"
								comparison = 50
								condition_result = 0
							},
						]
					}

					issue_detection = {
						type = "static"
					}
				`),
				ConfigStateChecks: append(
					checks,
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("enabled"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("name"), knownvalue.StringExact(monitorName+"-updated")),
				),
			},
		},
	})
}

func TestAccMetricMonitorDataSource_dynamic(t *testing.T) {
	teamName := acctest.RandomWithPrefix("tf-team")
	projectName := acctest.RandomWithPrefix("tf-project")
	monitorName := acctest.RandomWithPrefix("tf-metric-monitor")
	rn := "data.sentry_metric_monitor.test"

	checks := []statecheck.StateCheck{
		statecheck.ExpectKnownValue(rn, tfjsonpath.New("id"), knownvalue.NotNull()),
		statecheck.ExpectKnownValue(rn, tfjsonpath.New("organization"), knownvalue.StringExact(acctest.TestOrganization)),
		statecheck.ExpectKnownValue(rn, tfjsonpath.New("project_id"), knownvalue.NotNull()),
		statecheck.ExpectKnownValue(rn, tfjsonpath.New("owner"), knownvalue.ObjectExact(map[string]knownvalue.Check{
			"user_id": knownvalue.Null(),
			"team_id": knownvalue.NotNull(),
		})),
	}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccMetricMonitorDataSourceConfig(teamName, projectName, monitorName, `
					aggregate = "count()"
					dataset = "events"
					event_types = ["default", "error"]
					query = "is:unresolved"
					query_type = "error"
					time_window_seconds = 3600

					condition_group = {
						conditions = [
							{
								type = "anomaly_detection"
								comparison_sensitivity = "high"
								comparison_threshold_type = "above_and_below"
								condition_result = 75
							},
						]
					}

					issue_detection = {
						type = "dynamic"
					}
				`),
				ConfigStateChecks: append(
					checks,
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("enabled"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("name"), knownvalue.StringExact(monitorName)),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("aggregate"), knownvalue.StringExact("count()")),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("dataset"), knownvalue.StringExact("events")),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("event_types"), knownvalue.SetExact([]knownvalue.Check{
						knownvalue.StringExact("default"),
						knownvalue.StringExact("error"),
					})),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("query"), knownvalue.StringExact("is:unresolved")),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("query_type"), knownvalue.StringExact("error")),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("time_window_seconds"), knownvalue.Int64Exact(3600)),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("condition_group"), knownvalue.ObjectExact(map[string]knownvalue.Check{
						"logic_type": knownvalue.StringExact("any"),
						"conditions": knownvalue.ListExact([]knownvalue.Check{
							knownvalue.ObjectPartial(map[string]knownvalue.Check{
								"type":                      knownvalue.StringExact("anomaly_detection"),
								"comparison_sensitivity":    knownvalue.StringExact("high"),
								"comparison_threshold_type": knownvalue.StringExact("above_and_below"),
								"condition_result":          knownvalue.Int64Exact(75),
							}),
						}),
					})),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("issue_detection"), knownvalue.ObjectExact(map[string]knownvalue.Check{
						"type":             knownvalue.StringExact("dynamic"),
						"comparison_delta": knownvalue.Null(),
					})),
				),
			},
		},
	})
}

func testAccMetricMonitorDataSourceConfig(teamName, projectName, name, extras string) string {
	return testAccMetricMonitorResourceConfig(teamName, projectName, name, extras) + `
		data "sentry_metric_monitor" "test" {
			organization = sentry_metric_monitor.test.organization
			id           = sentry_metric_monitor.test.id
		}
	`
}
