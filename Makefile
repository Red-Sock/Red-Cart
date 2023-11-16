include rscli.mk

dep:
	go install github.com/gojuno/minimock/v3/cmd/minimock@latest
	go install github.com/pressly/goose/v3/cmd/goose@latest


mock:
	minimock -i github.com/Red-Sock/go_tg/interfaces.Chat -o tests/mocks -g -s "_mock.go"

migrate-up:
	goose -dir migration postgres "user=red_cart dbname=red_cart host=localhost port=5432 sslmode=disable" up

