IMAGE_NAME=trusch/sha3
IMAGE_BUILDER=sudo podman

install: $(GOPATH)/bin/sha3

image: go.mod go.sum main.go Dockerfile
	 $(IMAGE_BUILDER) build -t $(IMAGE_NAME) .

$(GOPATH)/bin/sha3: go.mod go.sum main.go
	@go install -v
