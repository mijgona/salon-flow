FROM golang:1.21-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /salon-crm ./cmd/app

FROM alpine:3.19
RUN apk --no-cache add ca-certificates
WORKDIR /app
COPY --from=builder /salon-crm .
COPY configs/ ./configs/
COPY migrations/ ./migrations/

EXPOSE 8080
CMD ["./salon-crm"]
