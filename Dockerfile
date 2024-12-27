FROM golang:1.23-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o upassed-statistic-service ./cmd/main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates
RUN mkdir -p /upassed-statistic-service/config
COPY --from=builder /app/upassed-statistic-service /upassed-statistic-service/upassed-statistic-service
COPY --from=builder /app/config/* /upassed-statistic-service/config
RUN chmod +x /upassed-statistic-service/upassed-statistic-service
ENV APP_CONFIG_PATH="/upassed-statistic-service/config/local.yml"
EXPOSE 44044
CMD ["/upassed-statistic-service/upassed-statistic-service"]
