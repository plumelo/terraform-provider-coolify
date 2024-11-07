# Retrieve all servers
data "coolify_servers" "all" {}

# Retrieve servers with a specific description and team_id
data "coolify_servers" "filtered" {
  filter {
    name   = "user"
    values = ["root"]
  }
  # (AND)
  filter {
    name   = "ip"
    values = ["127.0.0.1", "localhost", "host.docker.internal"] # (OR)
  }
}

output "all" {
  value = data.coolify_servers.all.servers
}

output "filtered" {
  value = data.coolify_servers.filtered.servers
}
