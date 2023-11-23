package tests

import (
	"context"
	"strconv"
	"testing"

	"github.com/Red-Sock/go_tg/model"
	"github.com/Red-Sock/go_tg/model/response"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/stretchr/testify/require"

	"github.com/Red-Sock/Red-Cart/internal/domain/user"
	"github.com/Red-Sock/Red-Cart/internal/transport/telegram/handlers/cart/add"
	"github.com/Red-Sock/Red-Cart/tests/mocks"
)

const (
	cartAddSuccessCreatedMessage = `Предметы были успешно добавлены в корзину!`
	errNotEnoughArgMessage       = `Чтобы добавить товар в корзину воспользуйтесь командой /add_item {id} {товар_1} {товар_2}
Пример: /add_item 2 беляши кола сникерс`
	errNotIntegerMessage = `Идентификатор корзины должен быть целочисленным и положительным`
	errNoIdInDBMessage   = `Корзины с id = 4 не существует`
)

func Test_Add(t *testing.T) {
	type arguments struct {
		h   *add.Handler
		In  *model.MessageIn
		Out *mocks.ChatMock
	}

	testCases := map[string]struct {
		create func() arguments
	}{

		"OK": {
			create: func() (a arguments) {
				app := CreateTestApp(UsePgDb, UseServiceV1)
				a.h = add.New(app.Srv.Cart())

				userId := GetUserID()

				a.In = &model.MessageIn{
					Ctx: context.Background(),
					Message: &tgbotapi.Message{
						From: &tgbotapi.User{
							ID: userId,
						},
					},
					Args: []string{"1", "сникерс", "баунти"},
				}

				newUser := user.User{
					Id: userId,
				}
				err := app.Db.User().Upsert(context.Background(), newUser)
				require.NoError(t, err, "error creating test cart")

				_, err = app.Db.Cart().Create(context.Background(), userId)
				require.NoError(t, err, "error creating test cart")

				a.Out = mocks.NewChatMock(t)
				a.Out.SendMessageMock.Expect(&response.MessageOut{
					Text: cartAddSuccessCreatedMessage,
				})
				return
			},
		},
		"NOT_ENOUGH": {
			create: func() (a arguments) {
				app := CreateTestApp(UsePgDb, UseServiceV1)
				a.h = add.New(app.Srv.Cart())

				userId := GetUserID()

				a.In = &model.MessageIn{
					Ctx: context.Background(),
					Message: &tgbotapi.Message{
						From: &tgbotapi.User{
							ID: userId,
						},
					},
					Args: []string{"1"},
				}

				newUser := user.User{
					Id: userId,
				}
				err := app.Db.User().Upsert(context.Background(), newUser)
				require.NoError(t, err, "error creating test cart")

				_, err = app.Db.Cart().Create(context.Background(), userId)
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
				app := CreateTestApp(UsePgDb, UseServiceV1)
				a.h = add.New(app.Srv.Cart())

				userId := GetUserID()

				a.In = &model.MessageIn{
					Ctx: context.Background(),
					Message: &tgbotapi.Message{
						From: &tgbotapi.User{
							ID: userId,
						},
					},
					Args: []string{"Точно не число", "сникерс", "баунти"},
				}

				newUser := user.User{
					Id: userId,
				}
				err := app.Db.User().Upsert(context.Background(), newUser)
				require.NoError(t, err, "error creating test cart")

				_, err = app.Db.Cart().Create(context.Background(), userId)
				require.NoError(t, err, "error creating test cart")

				a.Out = mocks.NewChatMock(t)
				a.Out.SendMessageMock.Expect(&response.MessageOut{
					Text: errNotIntegerMessage,
				})
				return
			},
		},

		//TODO если запускать все тесты, то тут ломается, потому что выше создается корзина
		"NO_ID_IN_DB": {
			create: func() (a arguments) {
				app := CreateTestApp(UsePgDb, UseServiceV1)
				a.h = add.New(app.Srv.Cart())

				userId := GetUserID()

				a.In = &model.MessageIn{
					Ctx: context.Background(),
					Message: &tgbotapi.Message{
						From: &tgbotapi.User{
							ID: userId,
						},
					},
					Args: []string{strconv.Itoa(int(userId)), "сникерс", "баунти"},
				}

				a.Out = mocks.NewChatMock(t)
				a.Out.SendMessageMock.Expect(&response.MessageOut{
					Text: errNoIdInDBMessage,
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
