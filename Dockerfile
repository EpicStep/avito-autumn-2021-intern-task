FROM golang:alpine AS builder

WORKDIR $GOPATH/src/github.com/EpicStep/avito-autumn-2021-intern-task/

COPY . .
RUN go build -o /go/bin/balance-service ./cmd/balance/main.go

FROM alpine:latest
COPY --from=builder /go/bin/balance-service /go/bin/
CMD ["/go/bin/balance-service"]