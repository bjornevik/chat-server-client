FROM golang:1.22 AS builder
WORKDIR /app
COPY go.mod ./
RUN go mod download
COPY . .
RUN go build -o client main.go

FROM golang:1.22
WORKDIR /app
COPY --from=builder /app/client .
ENV SERVER_ADDRESS=0.0.0.0
ENV SERVER_PORT=8080
CMD ["./client"]
