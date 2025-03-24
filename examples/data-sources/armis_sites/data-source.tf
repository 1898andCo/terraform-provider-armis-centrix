# Read in site information
data "armis_sites" "lab" {}

output "armis_lab_sites" {
  description = "Armis sites across the lab environment"
  value       = data.armis_sites.lab.sites[*].name
}
