# syntax=docker/dockerfile:1
FROM golang:1.23-bookworm AS build
WORKDIR /app
COPY backend/ ./
RUN --mount=type=secret,id=proxy_ca \
    if [ -f /run/secrets/proxy_ca ]; then export SSL_CERT_FILE=/run/secrets/proxy_ca; fi; \
    go mod tidy && go test ./... && CGO_ENABLED=0 go build -trimpath -o /app/fluxo .
FROM alpine:3.21
RUN --mount=type=secret,id=proxy_ca \
    if [ -f /run/secrets/proxy_ca ]; then export SSL_CERT_FILE=/run/secrets/proxy_ca; fi; \
    apk add --no-cache ca-certificates && addgroup -S fluxo && adduser -S fluxo -G fluxo
WORKDIR /app
COPY --from=build /app/fluxo /app/fluxo
RUN mkdir /app/uploads && chown fluxo:fluxo /app/uploads
USER fluxo
EXPOSE 8080
CMD ["/app/fluxo"]
