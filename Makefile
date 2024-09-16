
.PHONY: keys dev-tools deps deps-cleancache check

dev-tools:
	@echo "Installing dev tools..."
	@go install go.uber.org/mock/mockgen@latest
	@go install github.com/air-verse/air@latest
	@go install github.com/a-h/templ/cmd/templ@latest
	@echo "Done."

deps:
	@echo "Installing dependencies..."
	@go mod tidy
	@echo "Done."

deps-cleancache:
	@go clean -modcache

check:
	golangci-lint run

keys:
	@echo "Generating secret keys..."
	@mkdir -p keys/
	@openssl genpkey -algorithm ed25519 -out keys/jwt-private.pem
	@openssl pkey -in keys/jwt-private.pem -pubout -out keys/jwt-public.pem
	@echo "Done."
