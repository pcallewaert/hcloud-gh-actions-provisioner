# Example packer file to create packer image

## Create packer

export your hcloud api token:

```sh
export HCLOUD_TOKEN=<your token>
```

```sh
packer build hcloud-buildserver.json
```