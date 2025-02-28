# Create a role for auditing
resource "armis_role" "auditor" {
  name = "Auditor"

  permissions = {
    advanced_permissions = {
      all = false
      behavioral = {
        all              = false
        application_name = false
        host_name        = true
        service_name     = false
      }
      device = {
        all           = false
        device_names  = true
        ip_addresses  = true
        mac_addresses = true
        phone_numbers = false
      }
    }
    alert = {
      all = false
      manage = {
        all               = false
        resolve           = true
        suppress          = true
        whitelist_devices = false
      }
    }
  }
}
