# `provider-natssecrets` Helm Chart

This Helm chart deploys the `provider-natssecrets` container, which provides a NATS-based secrets provider for Crossplane.

## Prerequisites

- A Kubernetes cluster with Crossplane installed
- A Vault instance for storing secrets

You can customize the installation by creating a values.yaml file and specifying the values you want to change. For more information on the available values, please refer to the values.yaml file in this chart.

## Configuration

| Field     | Description                                                | Default Value     |
| --------- | ---------------------------------------------------------- | ----------------- |
| namespace | Specifies the namespace where the resource will be created | crossplane-system |

The namespace field in the values.yaml file specifies the namespace where the resource will be created. By default, this is set to crossplane-system.

### Provider

The provider section of the values.yaml file specifies the configuration for the provider container. You can specify the image tag, pull policy, number of replicas, and command-line arguments for the container.

| Field             | Description                                                                               | Default Value |
| ----------------- | ----------------------------------------------------------------------------------------- | ------------- |
| package           | The Docker image for the provider container.                                              | N/A           |
| tag               | The tag for the Docker image. If not specified, the chart's appVersion will be used.      | Chart version |
| packagePullPolicy | The pull policy for the Docker image. By default, this is set to IfNotPresent.            | IfNotPresent  |
| replicas          | The number of replicas for the provider container. By default, this is set to 1.          | 1             |
| args              | The command-line arguments for the provider container. By default, this is an empty list. | []            |

### Vault

The vault section of the values.yaml file specifies the configuration for accessing a Vault instance. You can specify the credentials for accessing the Vault instance, including the address, TLS settings, token, and path to the secrets.

| Field                                       | Description                                                                                    | Default Value                |
| ------------------------------------------- | ---------------------------------------------------------------------------------------------- | ---------------------------- |
| credentials.secretRef.name                  | The name of the Kubernetes secret containing the credentials for accessing the Vault instance. | provider-natssecrets         |
| credentials.secretRef.key                   | The key in the Kubernetes secret containing the credentials for accessing the Vault instance.  | credentials                  |
| credentials.data.address                    | The address of the Vault instance.                                                             | https://vault.vault.svc:8200 |
| credentials.data.tls                        | Whether to use TLS when connecting to the Vault instance.                                      | true                         |
| credentials.data.insecure                   | Whether to skip TLS verification when connecting to the Vault instance.                        | true                         |
| credentials.data.token.value                | The token for accessing the Vault instance. If not defined, then fromSecret is used.           | N/A (must be set by user)    |
| credentials.data.token.fromSecret.name      | The name of the Kubernetes secret containing the token for accessing the Vault instance.       | bank-vaults                  |
| credentials.data.token.fromSecret.namespace | The namespace of the Kubernetes secret containing the token for accessing the Vault instance.  | vault                        |
| credentials.data.token.fromSecret.key       | The key in the Kubernetes secret containing the token for accessing the Vault instance.        | vault-root                   |
