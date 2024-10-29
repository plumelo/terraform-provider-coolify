resource "tls_private_key" "example" {
  algorithm = "ED25519"
}

resource "coolify_private_key" "example" {
  name        = "Example Terraformed Key"
  description = "Managed by Terraform"
  private_key = tls_private_key.test.private_key_pem
}

output "public_key" {
  value     = tls_private_key.example.public_key_pem
  sensitive = false
}
