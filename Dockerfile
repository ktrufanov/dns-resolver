FROM golang:1.15 AS builder

WORKDIR $GOPATH/src/github.com/ktrufanov/dns-resolver
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /dns_resolver dns_resolver.go
FROM alpine:latest
COPY --from=builder /dns_resolver /dns_resolver

CMD ["/dns_resolver"]
