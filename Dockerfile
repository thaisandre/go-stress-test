FROM golang:1.22-alpine AS builder

WORKDIR /app
COPY go.mod ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o /stress-test .

FROM alpine:3.20

RUN apk --no-cache add ca-certificates
COPY --from=builder /stress-test /usr/local/bin/stress-test

ENTRYPOINT ["stress-test"]
