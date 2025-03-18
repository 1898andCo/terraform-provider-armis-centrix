# Read in collector information
data "armis_collectors" "lab" {}

output "armis_lab_collectors" {
  description = "Armis virtual appliances in the lab"
  value       = data.armis_collectors.lab.collectors[*].name
}
