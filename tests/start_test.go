package tests

import (
	"context"
	"testing"

	"github.com/Red-Sock/go_tg/model"
	"github.com/Red-Sock/go_tg/model/keyboard"
	"github.com/Red-Sock/go_tg/model/response"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/stretchr/testify/require"

	"github.com/Red-Sock/Red-Cart/internal/domain"
	"github.com/Red-Sock/Red-Cart/internal/transport/telegram/handlers/start"
	"github.com/Red-Sock/Red-Cart/tests/mocks"
)

func Test_Start(t *testing.T) {

	const (
		successCreatedMessage  = `Hello New User!`
		successReturnedMessage = `Welcome Back!`
	)

	ourContext := context.Background()

	type arguments struct {
		h   *start.Handler
		In  *model.MessageIn
		Out *mocks.ChatMock
	}

	testCases := map[string]struct {
		create func() arguments
	}{

		"OK_FIRST_TIME": {
			create: func() (a arguments) {
				app := CreateTestApp(UsePgDb, UseServiceV1)
				a.h = start.New(app.Srv.User(), app.Srv.Cart())

				userId := GetUserID()

				a.In = &model.MessageIn{
					Ctx: ourContext,
					Message: &tgbotapi.Message{
						From: &tgbotapi.User{
							ID: userId,
						},
					},
				}

				a.Out = mocks.NewChatMock(t)
				kb := keyboard.InlineKeyboard{}
				kb.AddButton("Создать корзину", "/create_cart")
				kb.AddButton("Добавить товар", "/add_item")
				a.Out.SendMessageMock.Expect(&response.MessageOut{
					Text: successCreatedMessage,
					Keys: &kb,
				})
				return
			},
		},
		"OK_USER_EXISTS": {
			create: func() (a arguments) {
				app := CreateTestApp(UsePgDb, UseServiceV1)
				a.h = start.New(app.Srv.User(), app.Srv.Cart())

				userId := GetUserID()

				err := app.Db.User().Upsert(ourContext, domain.User{ID: userId})
				require.NoError(t, err)

				a.In = &model.MessageIn{
					Ctx: ourContext,
					Message: &tgbotapi.Message{
						From: &tgbotapi.User{
							ID: userId,
						},
					},
				}

				a.Out = mocks.NewChatMock(t)
				kb := keyboard.InlineKeyboard{}
				kb.AddButton("Создать корзину", "/create_cart")
				kb.AddButton("Добавить товар", "/add_item")
				a.Out.SendMessageMock.Expect(&response.MessageOut{
					Text: successReturnedMessage,
					Keys: &kb,
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
