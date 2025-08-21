// Copyright (c) 1898 & Co.
// SPDX-License-Identifier: Apache-2.0

package provider_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAcc_RoleResource(t *testing.T) {
	resourceName := "armis_role.test"
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccRoleResourceConfig(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", "test-role"),

					resource.TestCheckResourceAttr(resourceName, "permissions.advanced_permissions.all", "false"),
					resource.TestCheckResourceAttr(resourceName, "permissions.advanced_permissions.behavioral.all", "false"),
					resource.TestCheckResourceAttr(resourceName, "permissions.advanced_permissions.behavioral.application_name", "false"),
					resource.TestCheckResourceAttr(resourceName, "permissions.advanced_permissions.behavioral.host_name", "false"),
					resource.TestCheckResourceAttr(resourceName, "permissions.advanced_permissions.behavioral.service_name", "false"),

					resource.TestCheckResourceAttr(resourceName, "permissions.advanced_permissions.device.all", "false"),
					resource.TestCheckResourceAttr(resourceName, "permissions.advanced_permissions.device.device_names", "false"),
					resource.TestCheckResourceAttr(resourceName, "permissions.advanced_permissions.device.ip_addresses", "false"),
					resource.TestCheckResourceAttr(resourceName, "permissions.advanced_permissions.device.mac_addresses", "false"),
					resource.TestCheckResourceAttr(resourceName, "permissions.advanced_permissions.device.phone_numbers", "false"),

					resource.TestCheckResourceAttr(resourceName, "permissions.alert.all", "false"),
					resource.TestCheckResourceAttr(resourceName, "permissions.alert.read", "true"),
					resource.TestCheckResourceAttr(resourceName, "permissions.alert.manage.all", "false"),
					resource.TestCheckResourceAttr(resourceName, "permissions.alert.manage.resolve", "true"),
					resource.TestCheckResourceAttr(resourceName, "permissions.alert.manage.suppress", "false"),
					resource.TestCheckResourceAttr(resourceName, "permissions.alert.manage.whitelist_devices", "true"),

					resource.TestCheckResourceAttr(resourceName, "permissions.device.all", "false"),
					resource.TestCheckResourceAttr(resourceName, "permissions.device.read", "true"),
					resource.TestCheckResourceAttr(resourceName, "permissions.device.manage.all", "false"),
					resource.TestCheckResourceAttr(resourceName, "permissions.device.manage.create", "true"),
					resource.TestCheckResourceAttr(resourceName, "permissions.device.manage.delete", "true"),
					resource.TestCheckResourceAttr(resourceName, "permissions.device.manage.edit", "true"),
					resource.TestCheckResourceAttr(resourceName, "permissions.device.manage.merge", "true"),
					resource.TestCheckResourceAttr(resourceName, "permissions.device.manage.request_deleted_data", "false"),
					resource.TestCheckResourceAttr(resourceName, "permissions.device.manage.tags", "true"),
					resource.TestCheckResourceAttr(resourceName, "permissions.device.manage.enforce.all", "true"),
					resource.TestCheckResourceAttr(resourceName, "permissions.device.manage.enforce.create", "true"),
					resource.TestCheckResourceAttr(resourceName, "permissions.device.manage.enforce.delete", "true"),

					resource.TestCheckResourceAttr(resourceName, "permissions.policy.all", "false"),
					resource.TestCheckResourceAttr(resourceName, "permissions.policy.manage", "false"),
					resource.TestCheckResourceAttr(resourceName, "permissions.policy.read", "true"),

					resource.TestCheckResourceAttr(resourceName, "permissions.report.all", "true"),
					resource.TestCheckResourceAttr(resourceName, "permissions.report.export", "true"),
					resource.TestCheckResourceAttr(resourceName, "permissions.report.read", "true"),
					resource.TestCheckResourceAttr(resourceName, "permissions.report.manage.all", "true"),
					resource.TestCheckResourceAttr(resourceName, "permissions.report.manage.create", "true"),
					resource.TestCheckResourceAttr(resourceName, "permissions.report.manage.delete", "true"),
					resource.TestCheckResourceAttr(resourceName, "permissions.report.manage.edit", "true"),
				),
			},
		},
	})
}

func testAccRoleResourceConfig() string {
	return `
resource "armis_role" "test" {
  name = "test-role"

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
        resolve           = true
        suppress          = false
        whitelist_devices = true
      }
    }

    device = {
      all  = false
      read = true
      manage = {
        all                  = false
        create               = true
        delete               = true
        edit                 = true
        merge                = true
        request_deleted_data = false
        tags                 = true
        enforce = {
          all    = true
          create = true
          delete = true
        }
      }
    }

    policy = {
      all    = false
      manage = false
      read   = true
    }

    report = {
      all    = true
      export = true
      read   = true
      manage = {
        all    = true
        create = true
        delete = true
        edit   = true
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
      audit_log         = false
      secret_key        = false
      security_settings = false

      boundary = {
        all  = true
        read = true
        manage = {
          all    = true
          create = true
          delete = true
          edit   = true
        }
      }

      business_impact = {
        all    = true
        manage = true
        read   = true
      }

      collector = {
        all    = false
        manage = false
        read   = false
      }

      custom_properties = {
        all    = true
        manage = true
        read   = true
      }

      integration = {
        all    = false
        manage = false
        read   = false
      }

      internal_ips = {
        all    = false
        manage = false
        read   = false
      }

      notifications = {
        all    = false
        manage = false
        read   = false
      }

      oidc = {
        all    = false
        manage = false
        read   = false
      }

      saml = {
        all    = false
        manage = false
        read   = false
      }

      sites_and_sensors = {
        all  = true
        read = true
        manage = {
          all     = true
          sensors = true
          sites   = true
        }
      }

      users_and_roles = {
        all  = true
        read = true
        manage = {
          all = true
          roles = {
            all    = true
            create = true
            delete = true
            edit   = true
          }
          users = {
            all    = true
            create = true
            delete = true
            edit   = true
          }
        }
      }
    }

    user = {
      all  = true
      read = true
      manage = {
        all    = true
        upsert = true
      }
    }

    vulnerability = {
      all  = true
      read = true
      manage = {
        all     = true
        ignore  = true
        resolve = true
        write   = true
      }
    }
  }
}
`
}
