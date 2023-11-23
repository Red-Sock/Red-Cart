package tests

import (
	"context"
	"testing"

	"github.com/Red-Sock/go_tg/model"
	"github.com/Red-Sock/go_tg/model/response"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/stretchr/testify/require"

	"github.com/Red-Sock/Red-Cart/internal/domain/user"
	"github.com/Red-Sock/Red-Cart/internal/transport/telegram/handlers/cart/create"
	"github.com/Red-Sock/Red-Cart/tests/mocks"
)

const (
	successCreatedMessage = `Корзина c id = 1 была успешно создана.
Друзья могут добавить корзину через
/add_item 1 имя_товара_1 имя_товара_2`
	errCreateMessage = `У вас уже есть корзина с идентификатором = 2`
)

func Test_Create(t *testing.T) {
	type arguments struct {
		h   *create.Handler
		In  *model.MessageIn
		Out *mocks.ChatMock
	}

	testCases := map[string]struct {
		create func() arguments
	}{

		"OK": {
			create: func() (a arguments) {
				app := CreateTestApp(UsePgDb, UseServiceV1)
				a.h = create.New(app.Srv.Cart())

				userId := GetUserID()

				a.In = &model.MessageIn{
					Ctx: context.Background(),
					Message: &tgbotapi.Message{
						From: &tgbotapi.User{
							ID: userId,
						},
					},
				}
				newUser := user.User{
					Id: userId,
				}
				err := app.Db.User().Upsert(context.Background(), newUser)
				require.NoError(t, err, "error creating test cart")

				a.Out = mocks.NewChatMock(t)
				a.Out.SendMessageMock.Expect(&response.MessageOut{
					Text: successCreatedMessage,
				})
				return
			},
		},
		"ERR_CART_EXISTS": {
			create: func() (a arguments) {
				app := CreateTestApp(UsePgDb, UseServiceV1)
				a.h = create.New(app.Srv.Cart())

				userId := GetUserID()

				newUser := user.User{
					Id: userId,
				}
				err := app.Db.User().Upsert(context.Background(), newUser)
				require.NoError(t, err, "error creating test cart")

				_, err = app.Db.Cart().Create(context.Background(), userId)
				require.NoError(t, err, "error creating test cart")
				a.In = &model.MessageIn{
					Ctx: context.Background(),
					Message: &tgbotapi.Message{
						From: &tgbotapi.User{
							ID: userId,
						},
					},
				}

				a.Out = mocks.NewChatMock(t)
				a.Out.SendMessageMock.Expect(&response.MessageOut{
					Text: errCreateMessage,
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
