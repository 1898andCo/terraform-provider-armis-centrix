# Read in role permission information
data "armis_role" "stakeholder" {
  name = "Stakeholder"
}

output "armis_stakeholder_permissions" {
  description = "Armis stakeholder role permissions"
  value       = data.armis_role.stakeholder.permissions
}
