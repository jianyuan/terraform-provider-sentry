package sentrydata

// https://github.com/getsentry/sentry/blob/master/src/sentry/constants.py#L223-L230
var LogLevels = []string{
	"sample",
	"debug",
	"info",
	"warning",
	"error",
	"fatal",
}

// https://github.com/getsentry/sentry/blob/master/src/sentry/constants.py#L223-L230
var LogLevelNameToId = map[string]string{
	"sample":  "0",
	"debug":   "10",
	"info":    "20",
	"warning": "30",
	"error":   "40",
	"fatal":   "50",
}

// https://github.com/getsentry/sentry/blob/master/src/sentry/constants.py#L223-L230
var LogLevelIdToName = map[string]string{
	"0":  "sample",
	"10": "debug",
	"20": "info",
	"30": "warning",
	"40": "error",
	"50": "fatal",
}

// https://github.com/getsentry/sentry/blob/master/src/sentry/issues/grouptype.py#L31-L39
var IssueGroupCategories = []string{
	"Error",
	"Performance",
	"Profile",
	"Cron",
	"Replay",
	"Feedback",
	"Uptime",
	"Metric_Alert",
}

// https://github.com/getsentry/sentry/blob/master/src/sentry/issues/grouptype.py#L31-L39
var IssueGroupCategoryNameToId = map[string]string{
	"Error":        "1",
	"Performance":  "2",
	"Profile":      "3",
	"Cron":         "4",
	"Replay":       "5",
	"Feedback":     "6",
	"Uptime":       "7",
	"Metric_Alert": "8",
}

// https://github.com/getsentry/sentry/blob/master/src/sentry/issues/grouptype.py#L31-L39
var IssueGroupCategoryIdToName = map[string]string{
	"1": "Error",
	"2": "Performance",
	"3": "Profile",
	"4": "Cron",
	"5": "Replay",
	"6": "Feedback",
	"7": "Uptime",
	"8": "Metric_Alert",
}

// https://github.com/getsentry/sentry/blob/master/src/sentry/rules/conditions/event_attribute.py#L41-L68
var EventAttributes = []string{
	"message",
	"platform",
	"environment",
	"type",
	"error.handled",
	"error.unhandled",
	"error.main_thread",
	"exception.type",
	"exception.value",
	"user.id",
	"user.email",
	"user.username",
	"user.ip_address",
	"http.method",
	"http.url",
	"http.status_code",
	"sdk.name",
	"stacktrace.code",
	"stacktrace.module",
	"stacktrace.filename",
	"stacktrace.abs_path",
	"stacktrace.package",
	"unreal.crash_type",
	"app.in_foreground",
	"os.distribution_name",
	"os.distribution_version",
}

// https://github.com/getsentry/sentry/blob/master/src/sentry/rules/match.py#L6-L22
var MatchTypes = []string{
	"CONTAINS",
	"ENDS_WITH",
	"EQUAL",
	"GREATER_OR_EQUAL",
	"GREATER",
	"IS_SET",
	"IS_IN",
	"LESS_OR_EQUAL",
	"LESS",
	"NOT_CONTAINS",
	"NOT_ENDS_WITH",
	"NOT_EQUAL",
	"NOT_SET",
	"NOT_STARTS_WITH",
	"NOT_IN",
	"STARTS_WITH",
}

// https://github.com/getsentry/sentry/blob/master/src/sentry/rules/match.py#L6-L22
var MatchTypeNameToId = map[string]string{
	"CONTAINS":         "co",
	"ENDS_WITH":        "ew",
	"EQUAL":            "eq",
	"GREATER_OR_EQUAL": "gte",
	"GREATER":          "gt",
	"IS_SET":           "is",
	"IS_IN":            "in",
	"LESS_OR_EQUAL":    "lte",
	"LESS":             "lt",
	"NOT_CONTAINS":     "nc",
	"NOT_ENDS_WITH":    "new",
	"NOT_EQUAL":        "ne",
	"NOT_SET":          "ns",
	"NOT_STARTS_WITH":  "nsw",
	"NOT_IN":           "nin",
	"STARTS_WITH":      "sw",
}

// https://github.com/getsentry/sentry/blob/master/src/sentry/rules/match.py#L6-L22
var MatchTypeIdToName = map[string]string{
	"co":  "CONTAINS",
	"ew":  "ENDS_WITH",
	"eq":  "EQUAL",
	"gte": "GREATER_OR_EQUAL",
	"gt":  "GREATER",
	"is":  "IS_SET",
	"in":  "IS_IN",
	"lte": "LESS_OR_EQUAL",
	"lt":  "LESS",
	"nc":  "NOT_CONTAINS",
	"new": "NOT_ENDS_WITH",
	"ne":  "NOT_EQUAL",
	"ns":  "NOT_SET",
	"nsw": "NOT_STARTS_WITH",
	"nin": "NOT_IN",
	"sw":  "STARTS_WITH",
}

// https://github.com/getsentry/sentry/blob/master/src/sentry/rules/match.py#L25-L29
var LevelMatchTypes = []string{
	"EQUAL",
	"GREATER_OR_EQUAL",
	"LESS_OR_EQUAL",
}

// https://github.com/getsentry/sentry/blob/master/src/sentry/models/dashboard_widget.py#L49-L78
var DashboardWidgetTypes = []string{
	"discover",
	"issue",
	"metrics",
	"error-events",
	"transaction-like",
	"spans",
}

// https://github.com/getsentry/sentry/blob/master/src/sentry/models/dashboard_widget.py#L128-L145
var DashboardWidgetDisplayTypes = []string{
	"line",
	"area",
	"stacked_area",
	"bar",
	"table",
	"big_number",
	"top_n",
}
