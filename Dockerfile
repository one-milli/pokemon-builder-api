# Goのビルド環境
FROM golang:1.16 as builder
WORKDIR /app
COPY . .
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o pokemon-builder-api .

# 実行環境
FROM alpine:latest  
WORKDIR /root/
COPY --from=builder /app/pokemon-builder-api .
CMD ["./pokemon-builder-api"]
