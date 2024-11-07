FROM golang:1.23.2

WORKDIR /app

COPY ./src ./

RUN go mod download

RUN go build cmd/api/main.go

EXPOSE 8080

CMD ["./main"]