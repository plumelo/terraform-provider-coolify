# Retrieve all private keys
data "coolify_private_keys" "all" {}

# Retrieve private keys with a specific description and team_id
data "coolify_private_keys" "filtered" {
  filter {
    name   = "description"
    values = ["Created by Coolify"]
  }
  # (AND)
  filter {
    name   = "team_id"
    values = ["0", "1"] # (OR)
  }
}

output "all" {
  value = data.coolify_private_keys.all.private_keys
  # note: private_keys.*.private_key is a sensitive value
  sensitive = true
}

output "filtered" {
  value = data.coolify_private_keys.filtered.private_keys
  # note: private_keys.*.private_key is a sensitive value
  sensitive = true
}
