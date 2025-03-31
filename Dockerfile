FROM golang:1.23-alpine

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . /app

RUN go build -o auth ./cmd/main.go

EXPOSE 8020

CMD ["./auth"]