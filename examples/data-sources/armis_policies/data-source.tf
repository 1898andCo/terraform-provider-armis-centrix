# Retrieve policies whose names begin with a shared prefix
# For example, use this to discover all lab policies with a "lab" prefix.
data "armis_policies" "lab" {
  match_prefix = "lab"
}

output "armis_lab_policy_ids" {
  description = "Identifiers for Armis policies matching the lab prefix"
  value       = [for policy in data.armis_policies.lab.policies : policy.id]
}

output "armis_lab_policy_names" {
  description = "Names for Armis policies matching the lab prefix"
  value       = [for policy in data.armis_policies.lab.policies : policy.name]
}
