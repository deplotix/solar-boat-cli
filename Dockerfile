FROM golang:1.21-alpine AS builder

WORKDIR /app
COPY . .
RUN go build -o /solar-boat-cli

FROM hashicorp/terraform:latest
COPY --from=builder /solar-boat-cli /solar-boat-cli
ENTRYPOINT ["/solar-boat-cli"] 
