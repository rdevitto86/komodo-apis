FROM golang:1.26 AS build

COPY komodo-forge-sdk-go /komodo-forge-sdk-go

WORKDIR /app

COPY komodo-auth-api/go.mod komodo-auth-api/go.sum ./
RUN go mod download

COPY komodo-auth-api ./
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /bin/komodo ./cmd/komodo-auth-api

FROM gcr.io/distroless/base-debian12
COPY --from=build /bin/komodo /komodo
COPY --from=build /app/internal/config/validation_rules.yml /app/config/validation_rules.yml
EXPOSE 7011
ENTRYPOINT ["/komodo"]