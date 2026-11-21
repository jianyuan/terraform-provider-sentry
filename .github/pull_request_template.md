<!--
PR title: Conventional Commits, scoped to the affected resource(s) or data source(s).
Always use the full Terraform name as the scope (`sentry_project`, not `project`).

  fix(sentry_project): keep project in state when a later step fails
  feat(sentry_project_codeowners): new resource
  fix(sentry_alert, sentry_issue_alert): handle null action intervals
  feat(sentry_project)!: remove deprecated `foo` attribute   <- `!` marks a breaking change

Types: feat, fix, docs, refactor, test, ci, chore. Use `deps` as the scope for dependency bumps.
-->

## What

<!-- What does this change? Link the issue it fixes, e.g. "Fixes #123". -->

## Why

<!-- The problem or use case. For bugs: what users saw (error, plan diff, panic). -->

## How

<!-- Approach, and anything reviewers should look at closely. Link Sentry API docs if the change depends on API behaviour. -->

## Testing

<!-- How you tested it, e.g. `make testacc TESTARGS='-run=TestAccProjectResource'`, plus the result. -->

## Checklist

- [ ] PR title follows Conventional Commits with the full resource or data source name(s) as scope, e.g. `sentry_project`
- [ ] Acceptance tests added or updated, and they pass against a real Sentry organization
- [ ] Examples in `examples/` added or updated for new resources, data sources or attributes
- [ ] `make generate` run, and the generated code and docs are committed
- [ ] `make lint` passes
- [ ] New code uses `supertypes` (see `CLAUDE.md`); generated files edited through `internal/providergen`, not by hand
- [ ] No change to existing schema or state, or the PR is marked breaking (`!`) and the migration is described above
