# Read in all Armis tags
data "armis_tags" "all" {}

output "armis_tags" {
  description = "All Armis tags"
  value       = data.armis_tags.all.tags
}

# Read Armis tags matching a prefix
data "armis_tags" "ot_tags" {
  match_prefix = "OT"
}

output "armis_ot_tags" {
  description = "Armis tags starting with OT"
  value       = data.armis_tags.ot_tags.tags
}

# Read Armis tags excluding a prefix
data "armis_tags" "non_ot_tags" {
  exclude_prefix = "OT"
}

output "armis_non_ot_tags" {
  description = "Armis tags not starting with OT"
  value       = data.armis_tags.non_ot_tags.tags
}
