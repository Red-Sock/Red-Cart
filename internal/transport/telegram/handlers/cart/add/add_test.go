package add

import (
	"context"
	"testing"

	"github.com/Red-Sock/go_tg/model"
	"github.com/Red-Sock/go_tg/model/response"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/stretchr/testify/require"

	"github.com/Red-Sock/Red-Cart/tests"
	"github.com/Red-Sock/Red-Cart/tests/mocks"
)

const (
	successCreatedMessage  = `Предметы были успешно добавлены в корзину!`
	errNotEnoughArgMessage = `Чтобы добавить товар в корзину воспользуйтесь командой /add_item {id} {товар_1} {товар_2}
Пример: /add_item 2 беляши кола сникерс`
	errNotIntegerMessage = `Идентификатор корзины должен быть целочисленным и положительным`
	userId               = int64(1)
)

func Test_Create(t *testing.T) {
	type arguments struct {
		h   *Handler
		In  *model.MessageIn
		Out *mocks.ChatMock
	}

	testCases := map[string]struct {
		create func() arguments
	}{

		"OK": {
			create: func() (a arguments) {
				app := tests.CreateTestApp(tests.UseInMemoryDb, tests.UseServiceV1)
				a.h = New(app.Srv.Cart())

				a.In = &model.MessageIn{
					Ctx: context.Background(),
					Message: &tgbotapi.Message{
						From: &tgbotapi.User{
							ID: userId,
						},
					},
					Args: []string{"1", "сникерс", "баунти"},
				}

				//Создание пользователя
				_, err := a.h.cartService.Create(a.In.Ctx, a.In.Message.From.ID)
				require.NoError(t, err, "error creating test cart")

				a.Out = mocks.NewChatMock(t)
				a.Out.SendMessageMock.Expect(&response.MessageOut{
					Text: successCreatedMessage,
				})
				return
			},
		},
		"NOT_ENOUGH": {
			create: func() (a arguments) {
				app := tests.CreateTestApp(tests.UseInMemoryDb, tests.UseServiceV1)
				a.h = New(app.Srv.Cart())

				a.In = &model.MessageIn{
					Ctx: context.Background(),
					Message: &tgbotapi.Message{
						From: &tgbotapi.User{
							ID: userId,
						},
					},
					Args: []string{"1"},
				}

				//Создание пользователя
				_, err := a.h.cartService.Create(a.In.Ctx, a.In.Message.From.ID)
				require.NoError(t, err, "error creating test cart")

				a.Out = mocks.NewChatMock(t)
				a.Out.SendMessageMock.Expect(&response.MessageOut{
					Text: errNotEnoughArgMessage,
				})
				return
			},
		},
		"NOT_INTEGER": {
			create: func() (a arguments) {
				app := tests.CreateTestApp(tests.UseInMemoryDb, tests.UseServiceV1)
				a.h = New(app.Srv.Cart())

				a.In = &model.MessageIn{
					Ctx: context.Background(),
					Message: &tgbotapi.Message{
						From: &tgbotapi.User{
							ID: userId,
						},
					},
					Args: []string{"Точно не число", "сникерс", "баунти"},
				}

				//Создание пользователя
				_, err := a.h.cartService.Create(a.In.Ctx, a.In.Message.From.ID)
				require.NoError(t, err, "error creating test cart")

				a.Out = mocks.NewChatMock(t)
				a.Out.SendMessageMock.Expect(&response.MessageOut{
					Text: errNotIntegerMessage,
				})
				return
			},
		},
	}

	for name, testCase := range testCases {
		name, testCase := name, testCase
		t.Run(name, func(t *testing.T) {
			args := testCase.create()

			args.h.Handle(args.In, args.Out)

			args.Out.MinimockFinish()
		})

	}
}
