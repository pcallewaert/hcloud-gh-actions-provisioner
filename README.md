# hcloud-gh-actions-provisioner

This is a CLI tool to dynamically add or remove self-hosted GitHub Actions runners to your organization with [hetzner cloud](https://www.hetzner.com/cloud) instances.

This can be needed for several reasons:

- Cost reduction: Hetzner Cloud is cheap to start with, but you can also scale your runners to a minimum in the weekend/night
- Need different VM size (default GitHub Actions runner is [2vCPU/7GB RAM](https://docs.github.com/en/actions/using-github-hosted-runners/about-github-hosted-runners#supported-runners-and-hardware-resources)). By using VMs with a lot of cores, you can reduce your workflow times(for example the parallel tests)

It's possible to configure the location and the server type of the instances.
The servers can be protected by a Hetzner Cloud firewall by defining the firewall name with the `--hcloud-firewall-name` flag.

A new instance is configured by injecting a bash script with a GitHub registration token and an optional label list.

## Different pools of servers

You can define multiple pools of servers by configuring different `--name-prefix`.
It will only manage the servers that start with the `name-prefix`. The `static-labels` defines the labels in Github Actions.

This can be used in your Github Actions workflow (eg. `--static-labels=build`):

```yaml
jobs:
  test:
    name: My cool job
    runs-on: [self-hosted, build]
```

## Requirements

- Github PAT (needs `admin:org` permissions)
- Hetzner Cloud token (needs `write` permissions)
- A snapshot image at Hetzner Cloud with your basic tools needed as runner. An example to create this image with [packer](https://packer.io) can be found in examples/packer. We need the image ID, this can be found with the `hcloud` cli: `hcloud image list`

## Installation

### Kubernetes

See [kubernetes example](examples/kubernetes)

### cronjob

You can also just run it as a cron job by adding a scale up and scale down in crontab.

## Flags

```sh
Help:

```sh
Usage of ./bin/hcloud-gh-actions-provisioner:
  -github-owner="": Github Organisation owner
  -github-pat="": Github Personal Access Token
  -hcloud-firewall-name="": Hetzner Firewall Name
  -hcloud-location="fsn1": Hetzner Location
  -hcloud-server-type="cpx21": Hetzner Server type
  -hcloud-token="": Hetzner Cloud API Token
  -image-snapshot=-1: Image ID of the snapshot to use
  -loglevel="INFO": Log level
  -name-prefix="hcloud-github-actions-": Name prefix of the servers
  -number-of-builders=-1: The number of builders that have to be scaled
  -static-labels="": Labels that are added to Github runner and hetzner server
```

## Development

Code is checked for linting with default [golangci-lint](https://golangci-lint.run). You can use `make lint` to run.
