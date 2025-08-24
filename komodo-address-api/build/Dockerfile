# ---------- build ----------
FROM golang:1.25 AS build
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /bin/komodo ./cmd/komodo-address-api

# ---------- runtime ----------
FROM gcr.io/distroless/base-debian12
COPY --from=build /bin/komodo /komodo
EXPOSE 7010
ENTRYPOINT ["/komodo"]