define setup_env
    $(eval ENV_FILE := .env)
    $(eval include .env)
    $(eval export)
endef

default: fmt lint install generate

fetch-schema:
	@curl -s https://raw.githubusercontent.com/coollabsio/coolify/main/openapi.yaml > tools/openapi.yml

loadEnv:
	$(call setup_env)

build:
	go build -v ./...

install: build
	go install -v ./...

lint:
	# todo: golangci-lint run

generate:
	cd tools; go generate ./...

fmt:
	gofmt -s -w -e .

test:
	go test -v -cover -timeout=2m -parallel=10 ./...

testacc: loadEnv
	TF_ACC=1 go test -v -cover -timeout 10m -run '^(TestAcc|TestProtocol6ProviderServerConfigure)' ./...

.PHONY: fetch-schema loadEnv fmt lint test testacc build install generate