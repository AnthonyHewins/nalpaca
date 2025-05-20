.DEFAULT: all
.PHONY: fmt test gen clean run help sql docker

# command aliases
test := CONFIG_ENV=test go test ./...

targets := nalpaca

VERSION ?= v?.?.?
COMMIT ?= $(shell git rev-list -1 HEAD)
IMAGE := docker.io/ahewins/nalpaca
BUILD_FLAGS := 
docker_bin ?= podman
ifneq (,$(wildcard ./vendor))
	$(info Found vendor directory; setting "-mod vendor" to any "go build" commands)
	BUILD_FLAGS += -mod vendor
endif

#======================================
# Builds
#======================================
$(targets): ## Build a target server binary
	go build $(BUILD_FLAGS) -ldflags "-X main.version=$(COMMIT)" -o bin/$@ ./cmd/$@

all: $(targets) ## Build all targets

#======================================
# Docker
#======================================
docker: ## build docker image 
	go mod tidy
	$(docker_bin) login docker.io
	$(docker_bin) build -t $(IMAGE) --build-arg target=nalpaca -f docker/Dockerfile .
	$(docker_bin) push $(IMAGE)

composer := docker-compose -f ./docker/compose.yaml 
compose: ## build docker compose
	$(composer) build

up: ## run docker compose
	$(composer) up

down: ## teardown docker compose
	$(composer) down

#======================================
# Running
#======================================
run-%: ## Run the server using .env variables
	export $$(cat .env | xargs) && ./bin/$(patsubst run-%,%,$@)

#======================================
# Protobuf
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
