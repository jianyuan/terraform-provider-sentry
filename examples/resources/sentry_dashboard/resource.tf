resource "sentry_dashboard" "main" {
  organization = data.sentry_organization.main.id
  title        = "Test dashboard"

  widget {
    title        = "Number of Errors"
    display_type = "big_number"
    interval     = "5m"
    widget_type  = "discover"

    query {
      fields     = ["count()"]
      aggregates = ["count()"]
      conditions = "!event.type:transaction"
      order_by   = "count()"
    }

    layout {
      x     = 0
      y     = 0
      w     = 1
      h     = 1
      min_h = 1
    }
  }

  widget {
    title        = "Number of Issues"
    display_type = "big_number"
    interval     = "5m"
    widget_type  = "discover"

    query {
      fields     = ["count_unique(issue)"]
      aggregates = ["count_unique(issue)"]
      conditions = "!event.type:transaction"
      order_by   = "count_unique(issue)"
    }

    layout {
      x     = 1
      y     = 0
      w     = 1
      h     = 1
      min_h = 1
    }
  }

  widget {
    title        = "Events"
    display_type = "line"
    interval     = "5m"
    widget_type  = "discover"

    query {
      name       = "Events"
      fields     = ["count()"]
      aggregates = ["count()"]
      conditions = "!event.type:transaction"
      order_by   = "count()"
    }

    layout {
      x     = 2
      y     = 0
      w     = 4
      h     = 2
      min_h = 2
    }
  }

  widget {
    title        = "Affected Users"
    display_type = "line"
    interval     = "5m"
    widget_type  = "discover"

    query {
      name       = "Known Users"
      fields     = ["count_unique(user)"]
      aggregates = ["count_unique(user)"]
      conditions = "has:user.email !event.type:transaction"
      order_by   = "count_unique(user)"
    }

    query {
      name       = "Anonymous Users"
      fields     = ["count_unique(user)"]
      aggregates = ["count_unique(user)"]
      conditions = "!has:user.email !event.type:transaction"
      order_by   = "count_unique(user)"
    }

    layout {
      x     = 1
      y     = 2
      w     = 1
      h     = 2
      min_h = 2
    }
  }

  widget {
    title        = "Handled vs. Unhandled"
    display_type = "line"
    interval     = "5m"
    widget_type  = "discover"

    query {
      name       = "Handled"
      fields     = ["count()"]
      aggregates = ["count()"]
      conditions = "error.handled:true"
      order_by   = "count()"
    }

    query {
      name       = "Unhandled"
      fields     = ["count()"]
      aggregates = ["count()"]
      conditions = "error.handled:false"
      order_by   = "count()"
    }

    layout {
      x     = 0
      y     = 2
      w     = 1
      h     = 2
      min_h = 2
    }
  }

  widget {
    title        = "Errors by Country"
    display_type = "world_map"
    interval     = "5m"
    widget_type  = "discover"

    query {
      fields     = ["count()"]
      aggregates = ["count()"]
      conditions = "!event.type:transaction has:geo.country_code"
      order_by   = "count()"
    }

    layout {
      x     = 4
      y     = 6
      w     = 2
      h     = 4
      min_h = 2
    }
  }

  widget {
    title        = "High Throughput Transactions"
    display_type = "table"
    interval     = "5m"
    widget_type  = "discover"

    query {
      fields     = ["count()", "transaction"]
      aggregates = ["count()"]
      columns    = ["transaction"]
      conditions = "!event.type:error"
      order_by   = "-count()"
    }

    layout {
      x     = 0
      y     = 6
      w     = 2
      h     = 4
      min_h = 2
    }
  }

  widget {
    title        = "Errors by Browser"
    display_type = "table"
    interval     = "5m"
    widget_type  = "discover"

    query {
      fields     = ["browser.name", "count()"]
      aggregates = ["count()"]
      columns    = ["browser.name"]
      conditions = "!event.type:transaction has:browser.name"
      order_by   = "-count()"
    }

    layout {
      x     = 5
      y     = 2
      w     = 1
      h     = 4
      min_h = 2
    }
  }

  widget {
    title        = "Overall User Misery"
    display_type = "big_number"
    interval     = "5m"
    widget_type  = "discover"

    query {
      fields     = ["user_misery(300)"]
      aggregates = ["user_misery(300)"]
    }

    layout {
      x     = 0
      y     = 1
      w     = 1
      h     = 1
      min_h = 1
    }
  }

  widget {
    title        = "Overall Apdex"
    display_type = "big_number"
    interval     = "5m"
    widget_type  = "discover"

    query {
      fields     = ["apdex(300)"]
      aggregates = ["apdex(300)"]
    }

    layout {
      x     = 1
      y     = 1
      w     = 1
      h     = 1
      min_h = 1
    }
  }

  widget {
    title        = "High Throughput Transactions"
    display_type = "top_n"
    interval     = "5m"
    widget_type  = "discover"

    query {
      fields     = ["transaction", "count()"]
      aggregates = ["count()"]
      columns    = ["transaction"]
      conditions = "!event.type:error"
      order_by   = "-count()"
    }

    layout {
      x     = 0
      y     = 4
      w     = 2
      h     = 2
      min_h = 2
    }
  }

  widget {
    title        = "Issues Assigned to Me or My Teams"
    display_type = "table"
    interval     = "5m"
    widget_type  = "issue"

    query {
      fields     = ["assignee", "issue", "title"]
      columns    = ["assignee", "issue", "title"]
      conditions = "assigned_or_suggested:me is:unresolved"
      order_by   = "priority"
    }

    layout {
      x     = 2
      y     = 2
      w     = 2
      h     = 4
      min_h = 2
    }
  }

  widget {
    title        = "Transactions Ordered by Misery"
    display_type = "table"
    interval     = "5m"
    widget_type  = "discover"

    query {
      fields     = ["transaction", "user_misery(300)"]
      aggregates = ["user_misery(300)"]
      columns    = ["transaction"]
      order_by   = "-user_misery(300)"
    }

    layout {
      x     = 2
      y     = 6
      w     = 2
      h     = 4
      min_h = 2
    }
  }

  widget {
    title        = "Errors by Browser Over Time"
    display_type = "top_n"
    interval     = "5m"
    widget_type  = "discover"

    query {
      fields     = ["browser.name", "count()"]
      aggregates = ["count()"]
      columns    = ["browser.name"]
      conditions = "event.type:error has:browser.name"
      order_by   = "-count()"
    }

    layout {
      x     = 4
      y     = 2
      w     = 1
      h     = 4
      min_h = 2
    }
  }
}