FROM golang:1.17 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o main .

FROM golang:1.17

WORKDIR /app

COPY --from=builder /app/main .

EXPOSE 3000

CMD ["./main"]
