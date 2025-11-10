terraform {
  required_version = ">= 1.5.0"
  required_providers {
    armis = {
      source = "1898andCo/armis-centrix"
    }
  }
}

# Configure provider elsewhere or via environment variables
# provider "armis" {}

data "armis_lists" "all" {}

output "armis_lists_names" {
  description = "All Armis list names"
  value       = data.armis_lists.all.lists[*].name
}
