FROM golang:1.24-alpine AS build

WORKDIR /app


RUN apk add --no-cache git ca-certificates upx


COPY go.mod go.sum* ./
RUN if [ -f go.mod ]; then go mod download; fi


COPY . .


RUN mkdir -p static config && \
    if [ -f index.html ] && [ ! -f static/index.html ]; then cp index.html static/index.html; fi && \
    if [ -f config.yml ] && [ ! -f config/config.yml ]; then cp config.yml config/config.yml; fi


RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /bin/warehouse ./cmd/api/main.go || \
    (echo "Tip: ensure go.mod module path matches imports (e.g., 'module warehouse')" && exit 1)


RUN upx -q --lzma /bin/warehouse || true


FROM alpine:3.20

WORKDIR /app
ENV GIN_MODE=release

RUN apk add --no-cache ca-certificates tzdata && update-ca-certificates

COPY --from=build /bin/warehouse /app/warehouse
COPY --from=build /app/static /app/static
COPY --from=build /app/config /app/config
COPY --from=build /app/migrations /app/migrations

EXPOSE 8080
ENTRYPOINT ["/app/warehouse"]
