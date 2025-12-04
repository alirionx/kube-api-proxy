# multi-stage Dockerfile: build a static Go binary, then copy into scratch
# Save as /home/ubuntu/kube-api-proxy/Dockerfile

FROM golang:1.25-alpine AS builder
WORKDIR /src

# add git for modules if needed
RUN apk add --no-cache git

COPY src/ .

ENV CGO_ENABLED=0 GOOS=linux GOARCH=amd64
# build a small, stripped binary
RUN go build -ldflags="-s -w" -o /kube-api-proxy .

FROM scratch
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /kube-api-proxy /kube-api-proxy

EXPOSE 8080
ENTRYPOINT ["/kube-api-proxy"]