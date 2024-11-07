# Retrieve a specific server
data "coolify_server_resources" "example" {
  uuid = "abc123"
}

output "server_resource_names" {
  value = sort(data.coolify_server_resources.example.server_resources[*].name)
}
