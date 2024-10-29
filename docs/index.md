---
page_title: "Provider: Coolify"
---

# Coolify Provider

The "coolify" provider facilitates interaction with resources supported by [Coolify](https://coolify.io/). Before using this provider, you must configure it with your credentials, typically by setting the environment variable `COOLIFY_TOKEN`. For instructions on obtaining an API token, refer to Coolify's [API documentation](https://coolify.io/docs/api-reference/authorization).

For detailed information on the available resources, please refer to the links in the navigation bar.

## Example Usage

```terraform
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

# Create a server
resource "coolify_server" "vps" {
  # ...
}
```

## Argument Reference

- `endpoint` (String) Coolify endpoint. If not set checks env for `COOLIFY_ENDPOINT`. Default: `https://app.coolify.io/api/v1`
- `token` (String, Sensitive) Coolify token. If not set checks env for `COOLIFY_TOKEN`.
