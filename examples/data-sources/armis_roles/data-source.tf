# Read in role permission information
data "armis_roles" "stakeholder" {
  name = "Stakeholder"
}

output "armis_stakeholder_permissions" {
  description = "Armis stakeholder role permissions"
  value       = data.armis_role.stakeholder.permissions
}

# Read in roles with a matching prefix
data "armis_roles" "this" {
  match_prefix = "custom-role-"
}

output "role_names" {
  value = [for p in data.armis_roles.this.roles : p.name]
}

