# Two-staged build for the production container
FROM golang:1.20 AS builder

WORKDIR /src

COPY go.mod .
COPY go.sum .

RUN go mod download
COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /bin/main .

# Runner container base image
FROM alpine:edge

WORKDIR /app

COPY --from=builder /bin/main .

CMD ["/app/main"]
