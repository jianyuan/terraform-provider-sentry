# import using the member ID and team slug from the URL:
# https://[org-slug].sentry.io/settings/teams/[team-slug]/members/
# https://[org-slug].sentry.io/settings/members/[member-id]/
terraform import sentry_team_member.default org-slug/team-slug/member-id
