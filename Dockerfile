FROM golang:1.18-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o main .

FROM alpine:3.14
WORKDIR /app
COPY --from=builder /app/main .
COPY config ./config
EXPOSE 8080
CMD ["./main"]
