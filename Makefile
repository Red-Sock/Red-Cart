include rscli.mk

makeDep:
	go install github.com/gojuno/minimock/v3/cmd/minimock@latest
	go install github.com/pressly/goose/v3/cmd/goose@latest


mock:
	minimock -i github.com/Red-Sock/go_tg/interfaces.Chat -o tests/mocks -g -s "_mock.go"

# Не знаю как именно должна выглядит миграция
upMigration:
	goose postgres "user=red_cart dbname=red_cart port=5444 sslmode=disable" up
