terraform {
  required_providers {
    grafana = {
      source = "grafana/grafana"
      version = ">= 1.28.2"
    }
  }
}

provider "grafana" {
  url = var.grafana_url
  auth = var.grafana_api_key
}
