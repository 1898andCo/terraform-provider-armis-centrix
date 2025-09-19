# Read in boundary information
data "armis_boundary" "lab" {}

output "armis_lab_boundaries" {
  description = "Armis lab boundaries."
  value       = data.armis_boundary.lab
}
