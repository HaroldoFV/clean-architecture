FROM golang:1.22

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN mv /app/cmd/ordersystem/.env /app/.env

RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/ordersystem

EXPOSE 8000 50051 8080

CMD ["./main"]