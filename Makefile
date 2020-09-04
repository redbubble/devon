SHELL = /bin/bash
VERSION ?= 0.0.1-pre5
GIT_HASH = $(shell git rev-parse --short HEAD)
PROJECT_NAME = devon

BUILDKITE_BRANCH ?= $(shell git rev-parse --abbrev-ref HEAD)
BUILDKITE_BUILD_NUMBER ?= 0

define GOLANG_CONTAINER
	docker container run --rm                          \
		--volume $(shell pwd):/go/src/github.com/redbubble/${PROJECT_NAME} \
		--workdir /go/src/github.com/redbubble/${PROJECT_NAME}       \
		--env GO111MODULE=on                                \
		--env CGO_ENABLED=0                                 \
		--interactive golang:1.15
endef

.PHONY: test
test:
	${GOLANG_CONTAINER} go test ./...

.PHONY: release
release:
	git tag -f -a "v${VERSION}" -m "Releasing version ${VERSION}"
	git push --tags
	GO111MODULE=on goreleaser --rm-dist

.PHONY: fmt
fmt:
	${GOLANG_CONTAINER} go fmt ./...
	git diff --exit-code

target/${PROJECT_NAME}:
	mkdir -p target
	${GOLANG_CONTAINER} go build -o $@ -ldflags "-s -X main.version=${VERSION}" -a -installsuffix cgo

.PHONY: install
install:
	GO111MODULE=on go install -ldflags "-s -X main.version=${VERSION}" -a -installsuffix cgo

# .PHONY: completions
# completions: target/${PROJECT_NAME}
# 	target/${PROJECT_NAME} completions bash > static/completions/${PROJECT_NAME}.bash
# 	target/${PROJECT_NAME} completions zsh > static/completions/${PROJECT_NAME}.zsh

.PHONY: clean
clean: ## Cleanup artifacts generated by build.
	@echo "--- :shower: Cleaning up :shower:"
	rm -rf target
