FROM golang:1.21-alpine AS builder

WORKDIR /usr/local/src

COPY ["go.mod", "go.sum", "./"]
RUN go mod download

COPY ./ ./
RUN go build -o ./bin/app cmd/api/main.go

FROM alpine AS runner

COPY --from=builder /usr/local/src/bin/app /app
COPY ./config/config.yaml ./config/config.yaml

CMD ["/app"]