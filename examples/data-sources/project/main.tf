terraform {
  required_providers {
    gitlab = {
      source = "custom/gitlab"
    }
  }
}

data "gitlab_project" "this" {
  path_with_namespace = "gitlab-org/api/client-go"
}

output "project" {
  value = data.gitlab_project.this
}
