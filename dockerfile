FROM golang:1.23.4 AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/bin/app ./main.go
FROM alpine:latest
RUN apk --no-cache add ca-certificates
COPY --from=builder /app/bin/app /app/bin/app
WORKDIR /app
EXPOSE 3000
ENV PORT=:3000
CMD ["/app/bin/app"]
# Use the following command to build the Docker image: