# ---------- build ----------
FROM golang:1.25 AS build

# make the sibling module available inside the image (context root contains komodo-internal-lib-apis-go)
COPY komodo-internal-lib-apis-go /komodo-internal-lib-apis-go

WORKDIR /app

# copy module files for komodo-auth-api from the context
COPY komodo-auth-api/go.mod komodo-auth-api/go.sum ./
RUN go mod download

# copy the komodo-auth-api source
COPY komodo-auth-api ./
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /bin/komodo ./cmd/komodo-auth-api

# ---------- runtime ----------
FROM gcr.io/distroless/base-debian12
COPY --from=build /bin/komodo /komodo
COPY --from=build /app/internal/config/validation_rules.yml /app/config/validation_rules.yml
EXPOSE 7001
ENTRYPOINT ["/komodo"]