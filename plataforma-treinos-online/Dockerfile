FROM golang:1.22-alpine AS builder

WORKDIR /app

COPY app/go.mod app/go.sum* ./
RUN go mod download

COPY app/ ./
RUN CGO_ENABLED=0 GOOS=linux go build -o plataforma-treinos .

FROM alpine:3.20

WORKDIR /app

COPY --from=builder /app/plataforma-treinos .

EXPOSE 8080

CMD ["./plataforma-treinos"]
