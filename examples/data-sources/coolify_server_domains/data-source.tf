# Retrieve a specific server
data "coolify_server_domains" "example" {
  uuid = "abc123"
}

output "server_domains_ips" {
  value = data.coolify_server_domains.example.server_domains[*].ip
}

output "server_all_domains" {
  value = flatten(data.coolify_server_domains.example.server_domains[*].domains)
}
