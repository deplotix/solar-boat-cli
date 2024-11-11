FROM golang:1.21-alpine AS builder

WORKDIR /app
COPY . .
RUN go build -o /solar-boat-cli

FROM alpine:3.18
COPY --from=builder /solar-boat-cli /solar-boat-cli
ENTRYPOINT ["/solar-boat-cli"] 
