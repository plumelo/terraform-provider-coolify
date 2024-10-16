fetch-schema:
	# @curl -s https://raw.githubusercontent.com/coollabsio/coolify/main/openapi.yaml > config/coolify_openapi.yml
	@cp ../coolify/openapi.yaml config/coolify_openapi.yml

install:
	@go install .

generate: fetch-schema
	@go generate ./...

testacc:
	TF_ACC=1 go test ./internal/provider/ -count=1 -v -cover -timeout 10m