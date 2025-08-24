APP=komodo-auth-api

.PHONY: run build test docker dev
run: ; PORT=7001 go run ./cmd/$(APP)
build: ; go build -o bin/$(APP) ./cmd/$(APP)
test: ; go test ./...
docker:
\tdocker build -t $(APP):dev .
dev:
\tdocker run --rm -p 7001:7001 -e PORT=7001 -e $(APP):dev