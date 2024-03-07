FROM golang as builder

WORKDIR /app
COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -o /deploy/server/redcart ./cmd/Red-Cart/main.go

FROM alpine

WORKDIR /app
COPY --from=builder ./deploy/server/ .

EXPOSE 8080

ENTRYPOINT ["./redcart"]