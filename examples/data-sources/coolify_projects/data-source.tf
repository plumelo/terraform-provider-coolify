# Retrieve all projects
data "coolify_projects" "all" {}

# Retrieve projects with a specific name and description
data "coolify_projects" "filtered" {
  filter {
    name   = "name"
    values = ["my project"]
  }
  # (AND)
  filter {
    name   = "description"
    values = ["description 1", "description 2"] # (OR)
  }
}

output "all" {
  value = data.coolify_projects.all.projects
}

output "filtered" {
  value = data.coolify_projects.filtered.projects
}
