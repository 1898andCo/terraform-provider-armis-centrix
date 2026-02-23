# Monthly security assessment report
resource "armis_report" "comprehensive_report" {
  report_name   = "Monthly Security Assessment"
  asq           = "in:devices riskLevel:\"High\" OR in:devices riskLevel:\"Critical\""
  email_subject = "Monthly High Risk Device Report"

  schedule = {
    email              = ["ciso@example.com"]
    repeat_amount      = "1"
    repeat_unit        = "Days"
    report_file_format = "csv"
    time_of_day        = "06:00"
    timezone           = "UTC"
  }

  export_configuration = {
    columns = {
      devices         = ["name", "ipAddress", "riskLevel", "lastSeen"]
      vulnerabilities = ["cveId", "severity", "description"]
    }
  }
}
