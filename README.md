[contributors-shield]: https://img.shields.io/github/contributors/edgefarm/provider-natssecrets.svg?style=for-the-badge
[contributors-url]: https://github.com/edgefarm/provider-natssecrets/graphs/contributors
[forks-shield]: https://img.shields.io/github/forks/edgefarm/provider-natssecrets.svg?style=for-the-badge
[forks-url]: https://github.com/edgefarm/provider-natssecrets/network/members
[stars-shield]: https://img.shields.io/github/stars/edgefarm/provider-natssecrets.svg?style=for-the-badge
[stars-url]: https://github.com/edgefarm/provider-natssecrets/stargazers
[issues-shield]: https://img.shields.io/github/issues/edgefarm/provider-natssecrets.svg?style=for-the-badge
[issues-url]: https://github.com/edgefarm/provider-natssecrets/issues
[license-shield]: https://img.shields.io/github/license/edgefarm/provider-natssecrets?logo=apache2&style=for-the-badge
[license-url]: https://opensource.org/license/apache-2-0
[release-shield]:  https://img.shields.io/github/release/edgefarm/provider-natssecrets.svg?style=for-the-badge&sort=semver
[release-url]: https://github.com/edgefarm/provider-natssecrets/releases
[tag-shield]:  https://img.shields.io/github/tag/edgefarm/provider-natssecrets.svg?include_prereleases&sort=semver&style=for-the-badge
[tag-url]: https://github.com/edgefarm/provider-natssecrets/tags
[ci-shield]:  https://img.shields.io/github/actions/workflow/status/edgefarm/provider-natssecrets/ci.yml?branch=main&style=for-the-badge
[ci-url]: https://github.com/edgefarm/provider-natssecrets/actions/workflows/ci.yml

[![Contributors][contributors-shield]][contributors-url]
[![Forks][forks-shield]][forks-url]
[![Stargazers][stars-shield]][stars-url]
[![Issues][issues-shield]][issues-url]
[![Apache 2.0 License][license-shield]][license-url]
[![Release][release-shield]][release-url]
[![Latest Tag][tag-shield]][tag-url]
[![CI][ci-shield]][ci-url]

# provider-natssecretssecrets

`provider-natssecretssecrets` is a [Crossplane](https://crossplane.io/) Provider
that implements [EdgeFarm's Vault Nats Secrets Plugin](https://github.com/edgefarm/vault-plugin-secrets-nats) as managed resources.

## Features

The provider supports the following resources:
- Operators
- Accounts
- Users
- Operator signing keys
- Account signing keys

## 🎯 Installation

Make sure you have Crossplane installed. See the [Crossplane installation guide](https://docs.crossplane.io/latest/software/install/)

Create a `Provider` resource:

```yaml
apiVersion: pkg.crossplane.io/v1
kind: Provider
metadata:
  name: provider-natssecrets
spec:
  package: ghcr.io/edgefarm/provider-natssecrets/provider-natssecrets:master
  packagePullPolicy: IfNotPresent
  revisionActivationPolicy: Automatic
  revisionHistoryLimit: 1
```

**NOTE: Instead of using package version `master` [have a look at the available versions](https://github.com/edgefarm/provider-natssecrets/pkgs/container/provider-natssecrets%2Fprovider-natssecrets)**

## 📖 Examples

You might find the [examples](examples) directory helpful. Every example in this directory is deployable in a `make dev` environment.

For a full spec of possible options 
* go to docs.crds.dev for 
[Operator](https://doc.crds.dev/github.com/edgefarm/provider-natssecrets/nats.crossplane.io/Operator/v1alpha1#spec-forProvider), [Account](https://doc.crds.dev/github.com/edgefarm/provider-natssecrets/nats.crossplane.io/Account/v1alpha1#spec-forProvider), [User](https://doc.crds.dev/github.com/edgefarm/provider-natssecrets/nats.crossplane.io/User/v1alpha1#spec-forProvider), [AccountSigningKey](https://doc.crds.dev/github.com/edgefarm/provider-natssecrets/nats.crossplane.io/AccountSigningKey/v1alpha1#spec-forProvider) and [OperatorSigningKey](https://doc.crds.dev/github.com/edgefarm/provider-natssecrets/nats.crossplane.io/OperatorSigningKey/v1alpha1#spec-forProvider)
* use the `kubectl explain` command

```bash
# How to use kubectl explain
$ kubectl explain operator.spec.forProvider
$ kubectl explain account.spec.forProvider
$ kubectl explain user.spec.forProvider
$ kubectl explain accountsigningkey.spec.forProvider
$ kubectl explain operatorsigningkey.spec.forProvider
```

### Examples of operator, account and user resources

<details>
  <summary>Example operator resource</summary>

```yaml
apiVersion: issue.natssecrets.crossplane.io/v1alpha1
kind: Operator
metadata:
  name: myoperator
spec:
  forProvider:
    syncAccountServer: true
    createSystemAccount: true
    claims:
      operator:
        accountServerUrl: "nats://nats.nats:4222"
        signingKeys:
          - opsk1
        strictSigningKeyUsage: false
  providerConfigRef:
    name: vault-creds
  writeConnectionSecretToRef:
    namespace: crossplane-system
    name: myoperator
```
</details>  


<details>
  <summary>Example sys account resource</summary>

```yaml
apiVersion: issue.natssecrets.crossplane.io/v1alpha1
kind: Account
metadata:
  name: sys
spec:
  forProvider:
    operator: myoperator
    useSigningKey: opsk1
    claims:
      account:
        signingKeys:
          - sask1
        limits:
          subs: -1
          conn: -1
          leafNodeConn: -1
          data: -1
          payload: -1
          wildcardExports: true
          imports: -1
          exports: -1
        exports:
          - name: account-monitoring-streams
            subject: "$SYS.ACCOUNT.*.>"
            type: Stream
            accountTokenPosition: 3
            description: Account specific monitoring stream
            infoURL: https://docs.nats.io/nats-server/configuration/sys_accounts
          - name: account-monitoring-services
            subject: "$SYS.ACCOUNT.*.*"
            type: Service
            responseType: Stream
            accountTokenPosition: 4
            description:
              "Request account specific monitoring services for: SUBSZ, CONNZ,
              LEAFZ, JSZ and INFO"
            infoURL: https://docs.nats.io/nats-server/configuration/sys_accounts
  providerConfigRef:
    name: vault-creds
  writeConnectionSecretToRef:
    namespace: crossplane-system
    name: sys

```
</details>  


<details>
  <summary>Example default sys account user resource</summary>

```yaml
apiVersion: issue.natssecrets.crossplane.io/v1alpha1
kind: User
metadata:
  name: default-push
spec:
  forProvider:
    operator: myoperator
    account: sys
    useSigningKey: sask1
    claims:
      user:
        data: -1
        payload: -1
        subs: -1
        pub:
          allow:
            - "$SYS.REQ.CLAIMS.LIST"
            - "$SYS.REQ.CLAIMS.UPDATE"
            - "$SYS.REQ.CLAIMS.DELETE"
        resp:
        sub:
          allow:
            - _INBOX.>
  providerConfigRef:
    name: vault-creds

```
</details>  


<details>
  <summary>Example standard account resource</summary>

```yaml
apiVersion: issue.natssecrets.crossplane.io/v1alpha1
kind: Account
metadata:
  name: myaccount
spec:
  forProvider:
    operator: myoperator
    claims:
      account:
        defaultPermissions:
          pub:
            allow:
              - foo
              - bar
        limits:
          subs: -1
          conn: -1
          leafNodeConn: -1
          data: -1
          payload: -1
          wildcardExports: true
          imports: -1
          exports: -1
  providerConfigRef:
    name: vault-creds
  writeConnectionSecretToRef:
    namespace: crossplane-system
    name: myaccount
```
</details>  


<details>
  <summary>Example user resource</summary>

```yaml
apiVersion: issue.natssecrets.crossplane.io/v1alpha1
kind: User
metadata:
  name: myuser
spec:
  forProvider:
    operator: myoperator
    account: myaccount
    claims:
      user:
        data: 100
        payload: 200
        subs: 300
        pub:
          allow:
            - foo
  providerConfigRef:
    name: vault-creds
  writeConnectionSecretToRef:
    namespace: crossplane-system
    name: myuser
```
</details>  

### Signing keys

Signing keys can be either generated by the provider or imported from an existing secret. The secret must contain the base64 encoded nkey seed.

<details>
  <summary>Example of an `OperatorSigningKey` resource importing a secret</summary>

```yaml
apiVersion: nkey.natssecrets.crossplane.io/v1alpha1
kind: OperatorSigningKey
metadata:
  name: opsk1
spec:
  forProvider:
    operator: myoperator
    config:
      import:
        secretRef:
          name: opsk1
          namespace: default
          key: seed
  providerConfigRef:
    name: vault-creds
---
apiVersion: v1
kind: Secret
metadata:
  name: opsk1
  namespace: default
data:
  # base64 encoded operator nkey seed
  seed: U09BT0dMWFpDUzVUU1ZTTVBMM01QUjYzM0JaQUI2VkNJS1FJM1RMVTRaNUxFRlZEM0syRVQ1TUtQVQo=
```
</details>  

<details>
  <summary>Example of an `AccountSigningKey` without importing a secret</summary>

```yaml
apiVersion: nkey.natssecrets.crossplane.io/v1alpha1
kind: AccountSigningKey
metadata:
  name: mykey1
spec:
  forProvider:
    operator: myoperator
    account: myaccount
  providerConfigRef:
    name: vault-creds
```
</details>  


## 🐞 Debugging

Just start the debugger of your choice to debug `cmd/provider/main.go`.
The only thing that is important is, that your KUBECONFIG points to a `dev` cluster with the CRDs deployed (see `Developing locally`).

## 🧪 Test environment

To test the provider locally, you can use `devspace` to spin up a local `kind` cluster with the following components installed:
- Hashicorp Vault (with custom TLS certificate)
- NATS Server
- Crossplane
- provider-natssecrets (this project)

To start the test environment, run the following command:

```console
$ devspace run create-kind-cluster
$ devspace run-pipeline init
$ devspace run-pipeline deploy-vault
$ devspace run-pipeline deploy-crossplane
$ devspace run-pipeline deploy-nats
```

Once the environment is up and running you can use the `nats` cli to connect to the NATS server and publish messages.

```console
# Create the account and user and get the creds for the user
$ devspace run-pipeline create-custom-nats-account
$ kubectl port-forward -n nats svc/nats 4222:4222 &
$ PID=$!

# Publish and subscribe using the creds previously fetched
$ docker run -it -d --rm --name nats-subscribe --network host -v $(pwd)/.devspace/creds/creds:/creds natsio/nats-box:0.13.4 nats sub -s nats://localhost:4222 --creds /creds foo 
$ docker run --rm -d -it --name nats-publish --network host -v $(pwd)/.devspace/creds/creds:/creds natsio/nats-box:0.13.4 nats pub -s nats://localhost:4222 --creds /creds foo --count 3 "Message {{Count}} @ {{Time}}"

# Log output shows that authenticating with the creds file works for pub and sub
$ docker logs nats-subscribe
14:49:35 Subscribing on foo 
[#1] Received on "foo"
Message 1 @ 2:49PM

[#2] Received on "foo"
Message 2 @ 2:49PM

[#3] Received on "foo"
Message 3 @ 2:49PM

# Cleanup
$ docker kill nats-subscribe
$ pkill $PID

```

# 🤝🏽 Contributing

Code contributions are very much **welcome**.

1. Fork the Project
2. Create your Branch (`git checkout -b AmazingFeature`)
3. Commit your Changes (`git commit -m 'Add some AmazingFeature")
4. Push to the Branch (`git push origin AmazingFeature`)
5. Open a Pull Request targetting the `beta` branch.
