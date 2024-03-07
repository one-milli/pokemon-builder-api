FROM golang:1.22 as builder
WORKDIR /app
COPY . .
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o pokemon-builder-api .

FROM alpine:latest  
WORKDIR /root/
COPY --from=builder /app/pokemon-builder-api .
CMD ["./pokemon-builder-api"]
