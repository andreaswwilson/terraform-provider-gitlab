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

data "gitlab_current_user" "this" {
}

output "current_user" {
  value = data.gitlab_current_user.this
}
