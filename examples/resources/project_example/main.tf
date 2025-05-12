terraform {
  required_providers {
    gitlab = {
      source = "custom/gitlab"
    }
  }
}

resource "gitlab_project" "this" {
  name        = "abc1234"
  description = "hehe"
}
