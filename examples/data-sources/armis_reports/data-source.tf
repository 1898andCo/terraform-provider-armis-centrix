# Retrieve all reports
data "armis_reports" "all" {}

# Retrieve a specific report by ID
data "armis_reports" "by_id" {
  report_id = "123"
}

# Retrieve reports by name
data "armis_reports" "weekly" {
  report_name = "Weekly Security Report"
}
