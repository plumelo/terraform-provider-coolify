terraform {
  required_providers {
    coolify = {
      source  = "sierrajc/coolify"
      version = "~> 0.0.1"
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

# ...
