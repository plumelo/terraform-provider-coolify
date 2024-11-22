resource "coolify_postgresql_database" "example" {
  name        = "Example Terraformed Database 3"
  description = "Managed by Terraform"

  server_uuid      = "rg8ks8c"
  project_uuid     = "uoswco88w8swo40k48o8kcwk"
  environment_name = "production"

  image             = "postgres:16-alpine"
  postgres_db       = "my_database"
  postgres_user     = "postgres"
  postgres_password = "hunter12"

  instant_deploy = false
}
