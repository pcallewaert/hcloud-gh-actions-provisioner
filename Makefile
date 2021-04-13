.PHONY: build lint docker-build

build:
	go build -o bin/hcloud-gh-actions-provisioner

lint:
	golangci-lint run

docker-build:
	docker build -t docker.pkg.github.com/pcallewaert/hcloud-gh-actions-provisioner/hcloud-gh-actions-provisioner .