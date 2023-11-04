include rscli.mk

dep:
	go install github.com/gojuno/minimock/v3/cmd/minimock@latest


mock:
	minimock -i github.com/Red-Sock/go_tg/interfaces.Chat -o tests/mocks -g -s "_mock.go"