terraform {
  backend "s3" {
    bucket         = "terraform-state-storage-586877430255"
    dynamodb_table = "terraform-state-lock-586877430255"
    region         = "us-west-2"

    // THIS MUST BE UNIQUE
    key = "smee-mayday.tfstate"
  }
}

provider "aws" {
  region = "us-west-2"
}

data "aws_ssm_parameter" "eks_cluster_endpoint" {
  name = "/eks/av-cluster-endpoint"
}

provider "kubernetes" {
  host = data.aws_ssm_parameter.eks_cluster_endpoint.value
  config_path = "~/.kube/config"
}

data "aws_ssm_parameter" "slack_webhook" {
  name = "/env/slack-webhook"
}

module "shipyard_prd" {
  source = "github.com/byuoitav/terraform//modules/kubernetes-deployment"

  // required
  name           = "mayday-prd"
  image          = "docker.pkg.github.com/byuoitav/mayday/mayday-dev"
  image_version  = "e1d51f9"
  container_port = 80 // doesn't actually exist in container
  repo_url       = "https://github.com/byuoitav/mayday"

  // optional
  image_pull_secret = "github-docker-registry"
  container_args = [
    "--webhook", data.aws_ssm_parameter.slack_webhook,
  ]
  health_check = false
}