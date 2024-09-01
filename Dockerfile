FROM golang:1.23-alpine as builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o ./build/main ./cmd/

ENV ENV=dev
ENV HTTP_PORT=8080
ENV DATABASE_URL="postgres://postgres:postgres@localhost:5432/test-kami"

EXPOSE 8080

ENTRYPOINT ["./build/main"]