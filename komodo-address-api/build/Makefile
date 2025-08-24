APP=komodo-address-api

.PHONY: run build test docker dev
run: ; PORT=7010 GEOCODER=mock go run ./cmd/$(APP)
build: ; go build -o bin/$(APP) ./cmd/$(APP)
test: ; go test ./...
docker:
\tdocker build -t $(APP):dev .
dev:
\tdocker run --rm -p 7010:7010 -e PORT=7010 -e GEOCODER=mock $(APP):dev