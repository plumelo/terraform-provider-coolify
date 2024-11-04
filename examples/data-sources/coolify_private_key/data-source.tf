# Retrieve a specific private key
data "coolify_private_key" "example" {
  uuid = "abc123"
}

# Example outputs
output "private_key_name" {
  value = data.coolify_private_key.example.name
}

output "private_key_description" {
  value = data.coolify_private_key.example.description
}

output "private_key_value" {
  value     = data.coolify_private_key.example.private_key
  sensitive = true
}

