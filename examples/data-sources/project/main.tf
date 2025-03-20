terraform {
  required_providers {
    gitlab = {
      source = "custom/gitlab"
    }
  }
}

provider "gitlab" {
  token = "glpat-mylittletoken" // BYTT MEG ELLER BRUK GITLAB_TOKEN miljø-variabel
}

data "gitlab_project" "this" {
}

output "project" {
  value = data.gitlab_project.this
}
