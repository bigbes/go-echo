build:
	go build -o .artifacts/echo $(realpath ./cmd/go-echo)/*.go

# Example:
# export SOVEREN_PUSH_DEBUG_REGISTRY_URI="ghcr.io/nbasov/debug"
.PHONY: push-registry
push-registry: build
	@[ "${SOVEREN_PUSH_DEBUG_REGISTRY_URI}" ] || ( echo "SOVEREN_PUSH_DEBUG_REGISTRY_URI is not set"; exit 1 )
	docker buildx create --use
	docker buildx build										\
	--platform linux/amd64,linux/arm64                      \
	--tag $(SOVEREN_PUSH_DEBUG_REGISTRY_URI):go-echo 		\
	--push 													\
	.

start:
	go run ./cmd/go-echo/main.go
