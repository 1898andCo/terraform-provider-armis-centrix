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

The provider includes two types of tests:

### Unit Tests

Unit tests validate the internal logic, data transformations, and utility functions without requiring API access. These tests are fast, run in parallel, and don't need authentication credentials.

```sh
# Run all unit tests
go test ./...

# Run unit tests with coverage
go test ./... -coverprofile=coverage.out
go tool cover -html=coverage.out

# Run specific package tests
go test ./internal/utils -v
go test ./armis -v

# Run specific test function
go test ./internal/utils -v -run TestBuildRoleRequest
```

**Current Coverage:**
- `internal/utils` package: 86.1% coverage (126 tests)
- `armis` SDK package: 73.3% coverage (51 tests)
- Overall project: 42.7% coverage

### Acceptance Tests

Acceptance tests validate the provider's integration with the Armis API and require valid API credentials. These tests create, read, update, and delete real resources in the configured Armis instance.

You'll need to set the API key environment variable:

```sh
export ARMIS_API_KEY=<API_KEY>
export ARMIS_API_URL=<API_URL>

# Runs acceptance tests (requires API credentials)
task test
```

**Note:** The `task test` command runs acceptance tests with `TF_ACC=true` and includes a cleanup sweep at the end.

### Test Structure

```
.
├── armis/                      # SDK tests
│   ├── auth_test.go           # Authentication tests
│   ├── collectors_test.go     # Collector CRUD tests
│   └── ...
├── internal/
│   ├── provider/              # Acceptance tests
│   │   ├── collector_resource_test.go
│   │   ├── policy_resource_test.go
│   │   └── ...
│   └── utils/                 # Unit tests
│       ├── roles_utils_test.go      # Role transformation tests
│       └── policy_utils_test.go     # Policy transformation tests
```

### Writing Tests

All tests follow Go best practices:

- **Table-driven tests** for comprehensive coverage
- **Parallel execution** with `t.Parallel()` for performance
- **Clear test names** describing what is being tested
- **Validation functions** for complex assertions

Example:

```go
func TestMyFunction(t *testing.T) {
    t.Parallel()

    tests := []struct {
        name     string
        input    MyInput
        expected MyOutput
    }{
        {
            name: "valid input",
            input: MyInput{Field: "value"},
            expected: MyOutput{Result: "expected"},
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            t.Parallel()
            result := MyFunction(tt.input)
            if result != tt.expected {
                t.Errorf("got %v, want %v", result, tt.expected)
            }
        })
    }
}
```

