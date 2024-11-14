terraform {
  required_providers {
    coolify = {
      source  = "sierrajc/coolify"
      version = "~> 0"
    }
  }
}

provider "coolify" {
  endpoint = "https://coolify.domain.com/api/v1"
  # Instead of setting token here, define a COOLIFY_TOKEN
  # environment variable, e.g. by adding the following line to .bashrc:
  # export COOLIFY_TOKEN="Your API token"
  token = "Your API token"
}

# Generate a new private key, and create a server with that key.

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
