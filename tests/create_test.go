package tests

import (
	"context"
	"testing"

	tgapi "github.com/Red-Sock/go_tg/interfaces"
	"github.com/Red-Sock/go_tg/model"
	"github.com/Red-Sock/go_tg/model/response"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/stretchr/testify/require"

	"github.com/Red-Sock/Red-Cart/internal/domain/user"
	"github.com/Red-Sock/Red-Cart/internal/transport/telegram/handlers/cart/create"
	"github.com/Red-Sock/Red-Cart/tests/mocks"
)

func Test_Create(t *testing.T) {
	const (
		successCreatedMessageRegExp = `Корзина c id = \d+? была успешно создана.
Друзья могут добавить корзину через
/add_item 1 имя_товара_1 имя_товара_2`
		errCreateMessageRegExp = `У вас уже есть корзина с идентификатором = \d+?`
	)

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
				a.Out.SendMessageMock.Set(func(out tgapi.MessageOut) {
					message, ok := out.(*response.MessageOut)
					require.Truef(t, ok, "output message must be of type *response.MessageOut but %T is passed", message)
					require.Regexpf(t, successCreatedMessageRegExp, message.Text, "unexpected message response text")
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
				a.Out.SendMessageMock.Set(func(out tgapi.MessageOut) {
					message, ok := out.(*response.MessageOut)
					require.Truef(t, ok, "output message must be of type *response.MessageOut but %T is passed", message)
					require.Regexpf(t, errCreateMessageRegExp, message.Text, "unexpected message response text")

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
