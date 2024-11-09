resource "coolify_project" "example" {
  # name        = "Example Terraformed Project 2"
  # description = "Managed by Terraform 3"
}

output "project_uuid" {
  value = coolify_project.example
}
