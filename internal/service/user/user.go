package user

import (
	"context"
	"fmt"

	errors "github.com/Red-Sock/trace-errors"

	"github.com/Red-Sock/Red-Cart/internal/domain"
)

type UsersService struct {
	userData domain.UserRepo
	cartData domain.CartRepo
}

func (u *UsersService) GetDefaultCart(ctx context.Context, userID int64) (resp domain.UserCart, err error) {
	user, err := u.userData.Get(ctx, userID)
	if err != nil {
		return domain.UserCart{}, err
	}

	if user == nil {
		return domain.UserCart{}, errors.New("no such user")
	}

	resp.User = *user

	resp.Cart, err = u.cartData.GetUserDefaultCart(ctx, userID)
	if err != nil {
		return domain.UserCart{}, err
	}

	resp.Cart.Items, err = u.cartData.ListCartItems(ctx, resp.Cart.ID)
	if err != nil {
		return domain.UserCart{}, err
	}

	return resp, nil
}

func New(uD domain.UserRepo, cartData domain.CartRepo) *UsersService {
	return &UsersService{
		userData: uD,
		cartData: cartData,
	}
}

func (u *UsersService) Start(ctx context.Context, newUser domain.User) (message domain.StartMessagePayload, err error) {
	user, err := u.userData.Get(ctx, newUser.ID)
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

	cart, err := u.cartData.GetByOwnerId(ctx, newUser.ID)
	if err != nil {
		return domain.StartMessagePayload{
			Msg: domain.DbErrorMsg,
		}, errors.Wrap(err, "error getting cart by owner")
	}

	if cart != nil {
		message.Cart = *cart
		message.Cart.Items, err = u.cartData.ListCartItems(ctx, cart.ID)
	} else {
		message.Cart.ID, err = u.createCartForUser(ctx, newUser)
	}
	if err != nil {
		return domain.StartMessagePayload{
			Msg: domain.DbErrorMsg,
		}, err
	}

	message.Msg += fmt.Sprintf(` 🛒

Корзина по умолчанию: %d

Для добавления продуктов просто введите их название
`, message.Cart.ID)

	return message, nil
}

func (u *UsersService) AddToDefaultCart(ctx context.Context, items []domain.Item, userID int64) (cart domain.UserCart, err error) {
	user, err := u.userData.Get(ctx, userID)
	if err != nil {
		return domain.UserCart{}, errors.New("error getting user data")
	}

	if user == nil {
		return domain.UserCart{}, errors.New("no such user")
	}

	cart.User = *user

	cart.Cart, err = u.cartData.GetUserDefaultCart(ctx, userID)
	if err != nil {
		return domain.UserCart{}, err
	}

	if cart.Cart.ID == 0 {
		return domain.UserCart{}, errors.New(fmt.Sprintf("Для пользователя id = %d не задано корзины по умолчанию ", userID))
	}

	err = u.cartData.AddCartItems(ctx, items, cart.Cart.ID, userID)
	if err != nil {
		return domain.UserCart{}, err
	}

	cart.Cart.Items, err = u.cartData.ListCartItems(ctx, cart.Cart.ID)
	if err != nil {
		return domain.UserCart{}, errors.Wrap(err, "error getting cartItems")
	}

	return cart, nil
}

func (u *UsersService) createCartForUser(ctx context.Context, user domain.User) (int64, error) {
	cartID, err := u.cartData.Create(ctx, user.ID)
	if err != nil {
		return 0, errors.Wrap(err, "error creating cart")
	}
	err = u.cartData.LinkUserToCart(ctx, user.ID, cartID)
	if err != nil {
		return 0, errors.Wrap(err, "error linking cart")
	}

	err = u.cartData.SetDefaultCart(ctx, user.ID, cartID)
	if err != nil {
		return 0, errors.Wrap(err, "error setting default cart")
	}

	return cartID, nil
}
