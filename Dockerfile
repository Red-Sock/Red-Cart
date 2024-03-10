FROM golang as builder

WORKDIR /app
COPY . .

RUN go build -o /deploy/server/redcart ./cmd/Red-Cart/main.go

FROM alpine

LABEL MATRESHKA_CONFIG_ENABLED=true

WORKDIR /app
COPY --from=builder ./deploy/server/ .
COPY --from=builder /app/config/config.yaml ./config/config.yaml
EXPOSE 8080

ENTRYPOINT ["./redcart"]