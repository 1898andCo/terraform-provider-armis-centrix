# Configure the Armis Centrix provider using the required_providers stanza.
# You may optionally use a version directive to prevent breaking
# changes occurring unannounced.
terraform {
  required_version = ">= 1.0.0"

  required_providers {
    armis = {
      source  = "1898andCo/armis-centrix"
      version = "~> 1.0"
    }
  }
}

# Provider values can also be set using the ARMIS_API_KEY and ARMIS_API_URL environment variables.
provider "armis" {
  api_key = var.api_key
  api_url = "https://example-lab.armis.com"
}
