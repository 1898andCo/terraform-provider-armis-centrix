# Retrieve policies whose names begin with a shared prefix while omitting a subset
# For example, match policies with the "lab" prefix but exclude archived lab policies.
data "armis_policies" "lab" {
  match_prefix   = "lab"
  exclude_prefix = "lab-archived"
}

output "armis_lab_policy_ids" {
  description = "Identifiers for active Armis lab policies"
  value       = [for policy in data.armis_policies.lab.policies : policy.id]
}

output "armis_lab_policy_names" {
  description = "Names for active Armis lab policies"
  value       = [for policy in data.armis_policies.lab.policies : policy.name]
}
