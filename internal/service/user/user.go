package user

import (
	"context"

	errors "github.com/Red-Sock/trace-errors"

	"github.com/Red-Sock/Red-Cart/internal/domain"
)

type UsersService struct {
	userData domain.UserRepo
	cartData domain.CartRepo
}

func New(uD domain.UserRepo, cartData domain.CartRepo) *UsersService {
	return &UsersService{
		userData: uD,
		cartData: cartData,
	}
}

func (u *UsersService) Start(ctx context.Context, newUser domain.User) (message domain.StartMessagePayload, err error) {
	user, err := u.userData.Get(ctx, newUser.Id)
	if err != nil {
		return domain.StartMessagePayload{Msg: domain.DbErrorMsg}, errors.Wrap(err, "error creating new user")
	}

	if user != nil {
		message.Msg = "С возвращением!"
		message.User = *user
	} else {
		err = u.userData.Upsert(ctx, newUser)
		if err != nil {
			return domain.StartMessagePayload{
				Msg: domain.DbErrorMsg,
			}, errors.Wrap(err, "error updating user's profile")
		}
		message.User = newUser
		message.Msg = "Добро пожаловать!"
	}

	cart, err := u.cartData.GetByOwnerId(ctx, newUser.Id)
	if err != nil {
		return domain.StartMessagePayload{
			Msg: domain.DbErrorMsg,
		}, errors.Wrap(err, "error getting cart by owner")
	}

	if cart != nil {
		message.Cart = *cart
	} else {
		message.Cart.Id, err = u.cartData.Create(ctx, newUser.Id)
		if err != nil {
			return domain.StartMessagePayload{
				Msg: domain.DbErrorMsg,
			}, errors.Wrap(err, "error creating cart")
		}

		message.Cart.OwnerId = message.User.Id
	}

	return message, nil
}

func (u *UsersService) getUser() {

}
