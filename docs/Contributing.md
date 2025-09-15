<!--
Copyright (c) 1898 & Co.
SPDX-License-Identifier: Apache-2.0
-->

# Contributing

## Armis Centrix Go SDK

The Armis Centrix Go SDK for authenticating and provisioning resources is located in the `armis` directory. This location may be temporary as we plan to migrate the SDK to a separate, versioned repository. Within this directory, `client.go` and `auth.go` handle authentication and provider requests, requiring an API key and the Armis instance URL to communicate and manage resources. The `models_` files define the Go structs used for API calls and responses. When creating a new data source or resource, a model and a corresponding CRUD operations file (e.g., `roles.go`, `collectors.go`, or `users.go`) must be implemented to interact with the API. Once these components are created and tested, they can be utilized by the provider in the `internal/provider` directory.

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

In addition, you can run task install to set up a developer overrides in your `~/.terraformrc.` This will then allow you to use your locally built provider binary.

```sh
task install
```

When you are finished using a local version of the provider, running `task uninstall` will remove all developer overrides.

```sh
task uninstall
```

To use a released provider in your Terraform environment, run [`terraform init`](https://www.terraform.io/docs/commands/init.html) and Terraform will automatically install the provider. To specify a particular provider version when installing released providers, see the [Terraform documentation on provider versioning](https://www.terraform.io/docs/configuration/providers.html#version-provider-versions).

To instead use a custom-built provider in your Terraform environment (e.g. the provider binary from the build instructions above), follow the instructions to [install it as a plugin.](https://www.terraform.io/docs/plugins/basics.html#installing-plugins) After placing the custom-built provider into your plugins directory, run `terraform init` to initialize it.

## Testing the Provider

> [!NOTE]
> Armis recommends using a tool such as Postman or Paw to quickly develop and test the Armis API.
This will enable the developer to quickly debug requests to and responses from the API. These calls
can then be implemented in your platform of choice.
>
> For more information on the Armis Centrix™ platform, refer to the Armis user guide.

You'll need to set the API key environment variable:

```sh
export ARMIS_API_KEY=<API_KEY>
export ARMIS_API_URL=<API_URL>

# Runs provider tests
task test
```

