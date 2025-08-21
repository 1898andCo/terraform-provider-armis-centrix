# Create a role for auditing
resource "armis_role" "auditor" {
  name = "Auditor"

  permissions = {
    advanced_permissions = {
      all = false
      behavioral = {
        all              = false
        application_name = false
        host_name        = false
        service_name     = false
      }
      device = {
        all           = false
        device_names  = false
        ip_addresses  = false
        mac_addresses = false
        phone_numbers = false
      }
    }

    alert = {
      all  = false
      read = true
      manage = {
        all               = false
        resolve           = false
        suppress          = false
        whitelist_devices = false
      }
    }

    device = {
      all  = false
      read = true
      manage = {
        all                  = false
        create               = false
        delete               = false
        edit                 = false
        merge                = false
        request_deleted_data = false
        tags                 = false
        enforce = {
          all    = false
          create = false
          delete = false
        }
      }
    }

    policy = {
      all    = false
      manage = false
      read   = true
    }

    report = {
      all    = false
      export = true
      read   = true
      manage = {
        all    = false
        create = false
        delete = false
        edit   = false
      }
    }

    risk_factor = {
      all  = false
      read = true
      manage = {
        all = false
        customization = {
          all     = false
          create  = false
          disable = false
          edit    = false
        }
        status = {
          all     = false
          ignore  = false
          resolve = false
        }
      }
    }

    settings = {
      all               = false
      audit_log         = true
      secret_key        = false
      security_settings = false

      boundary = {
        all  = false
        read = true
        manage = {
          all    = false
          create = false
          delete = false
          edit   = false
        }
      }

      business_impact = {
        all    = false
        manage = false
        read   = true
      }

      collector = {
        all    = false
        manage = false
        read   = true
      }

      custom_properties = {
        all    = false
        manage = false
        read   = true
      }

      integration = {
        all    = false
        manage = false
        read   = true
      }

      internal_ips = {
        all    = false
        manage = false
        read   = true
      }

      notifications = {
        all    = false
        manage = false
        read   = true
      }

      oidc = {
        all    = false
        manage = false
        read   = true
      }

      saml = {
        all    = false
        manage = false
        read   = true
      }

      sites_and_sensors = {
        all  = false
        read = true
        manage = {
          all     = false
          sensors = false
          sites   = false
        }
      }

      users_and_roles = {
        all  = false
        read = true
        manage = {
          all = false
          roles = {
            all    = false
            create = false
            delete = false
            edit   = false
          }
          users = {
            all    = false
            create = false
            delete = false
            edit   = false
          }
        }
      }
    }

    user = {
      all  = false
      read = true
      manage = {
        all    = false
        upsert = false
      }
    }

    vulnerability = {
      all  = false
      read = true
      manage = {
        all     = false
        ignore  = false
        resolve = false
        write   = false
      }
    }
  }
}

