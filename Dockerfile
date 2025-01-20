FROM golang:1.23.4

WORKDIR /app

COPY . .

RUN go mod download

# 이미지 저장 디렉토리
RUN mkdir -p /app/storage/image/dev && \
    mkdir -p /app/storage/image/profile && \
    chmod -R 755 /app/storage

WORKDIR /app/cmd

RUN CGO_ENABLED=0 GOOS=linux go build -o main .

EXPOSE 7777

CMD ["/app/cmd/main"]