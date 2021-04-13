.PHONY: build lint docker-build

build:
	go build -o bin/hcloud-gh-actions-provisioner

lint:
	golangci-lint run

docker-build:
	docker build -t ghcr.io/pcallewaert/hcloud-gh-actions-provisioner .