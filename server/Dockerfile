# Build stage
FROM golang:1.20-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o main main.go

# Run stage
FROM alpine:3.15
WORKDIR /app
RUN apk add tzdata && ln -snf /usr/share/zoneinfo/Asia/Seoul /etc/localtime
COPY --from=builder /app/main .
COPY app.env .
COPY wait-for.sh .
COPY db/migration ./db/migration

EXPOSE 8000
CMD ["/app/main"]