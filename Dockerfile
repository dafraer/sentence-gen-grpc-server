# Stage 1: Build
FROM golang:1.26.1-alpine AS builder
WORKDIR /app

#Copy and download dependencies first
COPY go.mod go.sum ./
RUN go mod download

#Copy the rest
COPY . .

RUN go build -o sengen ./cmd

# Stage 2: Run
FROM alpine:3.21
WORKDIR /app
COPY --from=builder /app/sengen .
EXPOSE 50051
CMD ["./sengen"]