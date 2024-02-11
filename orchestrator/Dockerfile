# syntax=docker/dockerfile:1
FROM golang:1.21

WORKDIR /app

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o app ./cmd/app

CMD ["./app"]