resource "coolify_application_envs" "example" {
  # Application UUID
  uuid = "mc8gw00wscww4gskgk0gwgw0"

  env {
    key   = "key1"
    value = "value1"
  }

  # Set value for key1 only on preview environments
  env {
    key        = "key1"
    value      = "value1-on-preview"
    is_preview = true
  }

  env {
    key   = "key2"
    value = "value2"
  }
}
