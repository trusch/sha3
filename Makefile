IMAGE_NAME=trusch/sha3
IMAGE_BUILDER=sudo podman

.GIT_COMMIT=$(shell git rev-parse HEAD)
.GIT_VERSION=$(shell git describe --tags 2>/dev/null || echo "$(.GIT_COMMIT)")
.GIT_UNTRACKEDCHANGES := $(shell git status --porcelain --untracked-files=no)
ifneq ($(.GIT_UNTRACKEDCHANGES),)
	.GIT_COMMIT := $(.GIT_COMMIT)-dirty
endif

install: $(GOPATH)/bin/sha3

image: go.mod go.sum main.go Dockerfile
	 $(IMAGE_BUILDER) build -t $(IMAGE_NAME):latest .
	 $(IMAGE_BUILDER) tag $(IMAGE_NAME):latest $(IMAGE_NAME):$(.GIT_COMMIT)
	 $(IMAGE_BUILDER) tag $(IMAGE_NAME):latest $(IMAGE_NAME):$(.GIT_VERSION)

$(GOPATH)/bin/sha3: go.mod go.sum main.go
	@go install -v
