.DEFAULT: nalpaca
.PHONY: fmt test gen clean run help sql docker

# command aliases
test := CONFIG_ENV=test go test ./...

targets := nalpaca

VERSION ?= v0.0.0
IMAGE := ahewins/nalpaca
build_flag_path := github.com/AnthonyHewins/nalpaca
BUILD_FLAGS := 
ifneq (,$(wildcard ./vendor))
	$(info Found vendor directory; setting "-mod vendor" to any "go build" commands)
	BUILD_FLAGS += -mod vendor
endif

#======================================
# Local builds
#======================================
$(targets): ## Build a target server binary
	go build $(BUILD_FLAGS) -o bin/$@ ./cmd/$@

#======================================
# Docker
#======================================
docker: ## build docker image w/ $IMAGE
	go mod tidy
	podman build -t $(IMAGE) -f docker/Dockerfile .
	podman push $(IMAGE) docker.io/$(IMAGE)

compose: ## build docker compose
	docker-compose -f ./docker/compose.yaml build

#======================================
# Running
#======================================
run-%: ## Run the server using .env variables
	export $$(cat .env | xargs) && ./bin/$(patsubst run-%,%,$@)

run-compose: ## Run a binary with docker compose
	docker-compose -f ./docker/compose.yaml up

#======================================
# Tooling
#======================================
proto: ## buf generate
	rm -rf gen
	buf generate

#======================================
# App hygiene
#======================================
clean: ## gofmt, go generate, then go mod tidy, and finally rm -rf bin/
	find . -iname *.go -type f -exec gofmt -w -s {} \;
	go generate ./...
	go mod tidy
	rm -rf ./bin

test: ## Run go vet, then test all files
	go vet ./...
	$(test)

help: ## Print help
	@printf "\033[36m%-30s\033[0m %s\n" "(target)" "Build a target binary in current arch for running locally: $(targets)"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
