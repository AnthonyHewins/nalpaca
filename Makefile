.DEFAULT: falpaca
.PHONY: fmt test gen clean run help sql

# command aliases
test := CONFIG_ENV=test go test ./...

targets := falpaca

VERSION ?= v0.0.0
IMAGE := ahewins/renaissance
build_flag_path := github.com/AnthonyHewins/falpaca
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
	docker build -t $(IMAGE) -f docker/Dockerfile .
	docker push $(IMAGE)

#======================================
# Running
#======================================
run-%: ## Run the server using .env variables
	export $$(cat .env | xargs) && ./bin/$(patsubst run-%,%,$@)

run-docker-compose: ## Run a binary with docker compose
	docker-compose -f ./docker/%-compose.yaml up

#======================================
# Tooling
#======================================
proto: ## buf generate
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
