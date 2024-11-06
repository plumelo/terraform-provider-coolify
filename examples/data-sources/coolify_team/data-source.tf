# Retrieve a specific team
data "coolify_team" "example" {
  id = 123
}

output "team_members" {
  value = data.coolify_team.example.members
}

# Retrieve the team for the current authenticated API Key
data "coolify_team" "current" {}

output "current_team_id" {
  value = data.coolify_team.current.id
}
