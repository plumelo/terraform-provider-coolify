
# Retrieve a specific project
data "coolify_project" "example" {
  uuid = "abc123"
}

output "project_name" {
  value = data.coolify_project.example.name
}

output "project_description" {
  value = data.coolify_project.example.description
}

output "project_environments" {
  value = data.coolify_project.example.environments[*].name
}
