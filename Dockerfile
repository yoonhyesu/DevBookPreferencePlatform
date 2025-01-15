FROM golang:1.23.4

WORKDIR /app

COPY . .

RUN go mod download

WORKDIR /app/cmd

RUN CGO_ENABLED=0 GOOS=linux go build -o main .

EXPOSE 7777

CMD ["/app/cmd/main"]