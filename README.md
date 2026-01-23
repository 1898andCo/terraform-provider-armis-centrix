<a href="https://opentofu.org">
  <picture>
    <img src=".github/opentofu.svg" alt="OpenTofu logo" title="OpenTofu" align="right" height="50">
  </picture>
</a>

# Terraform/OpenTofu Armis Centrix Provider

## Requirements

- [Terraform](https://www.terraform.io/downloads.html) >= 0.12 or [OpenTofu](https://opentofu.org/) >= 1.5
- [Go](https://golang.org/doc/install) >= 1.25 (to build the provider plugin)

## Authenticating to Armis Centrix

The authorization token is used for authentication of the Armis API. To obtain your secret key from the Armis console, go to **Settings > API Management**.

If the secret key has not already been created, do the following:

1. Click **Create** to create the secret key.
2. Click **Show** to access the secret key. The following dialog is displayed, from which you can copy
your secret key.
3. Set the `ARMIS_API_KEY` environment variable or declare in the Armis Centrix provider configuration with the `api_key` parameter.

## Examples

```terraform
terraform {
  required_version = ">= 1.5.0"

  required_providers {
    armis = {
      source = "1898andCo/armis-centrix"
    }
  }
}

provider "armis" {
  api_key = var.api_key
  api_url = var.api_url
}
```

## Building the Provider

> [!NOTE]
> The following installation uses [Taskfile](https://taskfile.dev/), which can be downloaded by running the following command:
>
> `brew install go-task`

Clone repository:

```sh
git clone https://github.com/1898andCo/terraform-provider-armis-centrix.git
```

Enter the provider directory and build the provider:

```sh
cd terraform-provider-armis-centrix
task build
```

To use a released provider in your Terraform environment, run [`terraform init`](https://www.terraform.io/docs/commands/init.html) and Terraform will automatically install the provider. To specify a particular provider version when installing released providers, see the [Terraform documentation on provider versioning](https://www.terraform.io/docs/configuration/providers.html#version-provider-versions).

To instead use a custom-built provider in your Terraform environment (e.g. the provider binary from the build instructions above), follow the instructions to [install it as a plugin.](https://www.terraform.io/docs/plugins/basics.html#installing-plugins) After placing the custom-built provider into your plugins directory, run `terraform init` to initialize it.

## LLM Support

This repository implements the [llms.txt specification](https://llmstxt.org/) to help Large Language Models better understand and interact with the codebase. See [`llms.txt`](./llms.txt) for structured documentation.

## Contributing

Check out our [Contributing Docs](./docs/Contributing.md) for more information on how to support new resources and data sources, test, and contribute to the provider!

For bug reports & feature requests, please use the [issue tracker](https://github.com/1898andCo/terraform-provider-armis-centrix/issues).

PRs are welcome! We follow the typical "fork-and-pull" Git workflow.
 1. **Fork** the repo on GitHub
 2. **Clone** the project to your own machine
 3. **Commit** changes to your own branch
 4. **Push** your work back up to your fork
 5. Submit a **Pull Request** so that we can review your changes

> [!TIP]
> Be sure to merge the latest changes from "upstream" before making a pull request!

