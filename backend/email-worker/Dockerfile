FROM golang:1.25 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o email-worker .

FROM alpine:3.22

WORKDIR /app

COPY --from=builder /app/email-worker .

CMD ["./backend/email-worker"]
