data "armis_lists" "all" {}

output "armis_lists_names" {
  description = "All Armis list names"
  value       = data.armis_lists.all.lists[*].name
}
