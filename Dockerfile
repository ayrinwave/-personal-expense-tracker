FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o server ./

FROM alpine:3.18

WORKDIR /root/
COPY --from=builder /app/server .

CMD ["./server"]

