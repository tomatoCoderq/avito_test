FROM golang:1.25.4-alpine3.21

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build src/api/main.go

EXPOSE 8080

CMD ["./main"]