
FROM golang:1.22-alpine AS builder

WORKDIR /comment-system

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o server .

FROM alpine:latest

WORKDIR /comment-system

COPY --from=builder /comment-system/server .

CMD ["./server"]