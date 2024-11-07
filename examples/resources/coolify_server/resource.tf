# Create a server using a Terraform generated private key

resource "tls_private_key" "example" {
  algorithm = "ED25519"
}

resource "coolify_private_key" "example" {
  name        = "Example Terraformed Key"
  description = "Managed by Terraform"
  private_key = tls_private_key.example.private_key_pem
}

resource "coolify_server" "example" {
  name             = "Example Terraformed Server"
  description      = "Managed by Terraform"
  ip               = "localhost"
  port             = 22
  user             = "root"
  private_key_uuid = coolify_private_key.example.uuid
  instant_validate = false
}
