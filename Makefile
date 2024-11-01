default: generate

fetch-schema:
	# @curl -s https://raw.githubusercontent.com/coollabsio/coolify/main/openapi.yaml > tools/openapi.yml
	@cp ../coolify/openapi.yaml tools/openapi.yml

install:
	@go install .

generate:
	cd tools; go generate ./...

test:
	go test -v -cover -timeout=2m -parallel=10 ./...

testacc:
	TF_ACC=1 go test -v -cover -timeout 10m ./internal/provider/...

.PHONY: fetch-schema install generate test testacc