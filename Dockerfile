# Builder
FROM golang:1.22 AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o weather-service

# Runtime
FROM gcr.io/distroless/base-debian10
COPY --from=builder /app/weather-service /weather-service
ENV PORT=8080
EXPOSE 8080
ENTRYPOINT ["/weather-service"]