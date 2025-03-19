FROM golang:1.23-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN apk add --no-cache gcc musl-dev && \
    CGO_ENABLED=1 GOOS=linux go build -o todo-api ./cmd/main.go

FROM alpine:latest

WORKDIR /root/

COPY --from=builder /app/todo-api .

EXPOSE 8080

CMD ["./todo-api"]