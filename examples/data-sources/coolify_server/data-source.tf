# Retrieve a specific server
data "coolify_server" "example" {
  uuid = "abc123"
}

output "server_address" {
  # user@ip:port
  value = "${data.coolify_server.example.user}@${data.coolify_server.example.ip}:${data.coolify_server.example.port}"
}
