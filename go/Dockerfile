FROM golang:1.24.4

WORKDIR /app

COPY go.mod .

RUN go mod tidy

CMD ["go", "run", "main.go"]
