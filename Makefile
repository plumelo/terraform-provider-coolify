default: fmt lint install generate

# @curl -s https://raw.githubusercontent.com/coollabsio/coolify/main/openapi.yaml > tools/openapi.yml
fetch-schema:
	@cp ../coolify/openapi.yaml tools/openapi.yml

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

testacc:
	TF_ACC=1 go test -v -cover -timeout 10m ./internal/provider/...

.PHONY: fetch-schema fmt lint test testacc build install generate