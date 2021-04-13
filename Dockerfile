FROM golang:1.16-alpine AS builder
RUN apk update && apk add --no-cache git
WORKDIR $GOPATH/src/github.com/pcallewaert/hcloud-gh-actions-provisioner
COPY . .
RUN GOOS=linux GOARCH=amd64 go build -mod vendor -ldflags="-w -s" -o /go/bin/hcloud-gh-actions-provisioner

FROM alpine
COPY --from=builder /go/bin/hcloud-gh-actions-provisioner /go/bin/hcloud-gh-actions-provisioner
ENTRYPOINT ["/go/bin/hcloud-gh-actions-provisioner"]