resource "armis_role" "example" {}

# Import an existing role by ID
import {
  to = armis_role.example
  id = "92012"
}
