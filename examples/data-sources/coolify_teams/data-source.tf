# Retrieve all teams
data "coolify_teams" "all" {}

# Retrieve teams with specific field values
data "coolify_teams" "filtered" {
  filter {
    name   = "discord_enabled"
    values = ["true"]
  }
  # (AND)
  filter {
    name   = "id"
    values = ["0", "1"] # (OR)
  }
}

output "all" {
  value = nonsensitive(data.coolify_teams.all.teams[*].name)
}

output "filtered" {
  value = nonsensitive(data.coolify_teams.filtered.teams[*].name)
}
