resource "grafana_notification_policy" "my_policy" {
  group_by      = ["alertname"]
  contact_point = grafana_contact_point.webhook.name

  policy {
    matcher {
      label = "a"
      match = "="
      value = "b"
    }
    group_by      = ["alertname"]
    contact_point = grafana_contact_point.webhook.name
  }
}

resource "grafana_contact_point" "webhook" {
  name = "Send to webhook"

  webhook {
    url     = var.grafana_alerts_webhook_url
    title   = "Alert FIRING"
    message = <<EOT
{{ len .Alerts.Firing }} alerts are firing!

Alert summaries:
{{ range .Alerts.Firing }}
{{ template "Alert Instance Template" . }}
{{ end }}
EOT
  }
}

resource "grafana_rule_group" "rule_group" {
  name             = "My Alert Rules"
  folder_uid       = grafana_folder.rule_folder.uid
  interval_seconds = 60

  rule {
    name      = "My Random Walk Alert"
    condition = "C"
    for       = "0s"

    // Query the datasource.
    data {
      ref_id = "A"
      relative_time_range {
        from = 600
        to   = 0
      }
      datasource_uid = grafana_data_source.testdata_datasource.uid
      // `model` is a JSON blob that sends datasource-specific data.
      // It's different for every datasource. The alert's query is defined here.
      model = jsonencode({
        intervalMs    = 1000
        maxDataPoints = 43200
        refId         = "A"
      })
    }
  }
}

resource "grafana_data_source" "testdata_datasource" {
  name = "TestData"
  type = "testdata"
}

resource "grafana_folder" "rule_folder" {
  title = "My Rule Folder"
}
