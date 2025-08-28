# ---------- build ----------
FROM golang:1.25 AS build
WORKDIR /app

# Go deps cache
COPY go.mod go.sum ./
RUN go mod download

# Source
COPY . .

# Build a static binary for distroless
ENV CGO_ENABLED=0
RUN GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o /out/server ./cmd/komodo-address-api

# ---------- runtime ----------
# Distroless has no shell; run the binary directly.
FROM gcr.io/distroless/static:nonroot

# Workdir is optional; keep paths explicit
WORKDIR /app

# Copy binary only
COPY --from=build /out/server /app/server

# Expose service port (mapped by orchestrator)
EXPOSE 7010

# Run the server directly; all configuration comes from runtime envs
ENTRYPOINT ["/app/server"]
