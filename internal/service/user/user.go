package user

import (
	"context"
	"fmt"

	errors "github.com/Red-Sock/trace-errors"

	"github.com/Red-Sock/Red-Cart/internal/domain"
	"github.com/Red-Sock/Red-Cart/scripts"
)

var (
	ErrNoDefaultCart = errors.New("Отсутствует корзина по умолчанию")
)

type Service struct {
	userData domain.UserRepo
	cartData domain.CartRepo
}

func (u *Service) GetCartByChat(ctx context.Context, userID int64) (domain.UserCart, error) {
	userCart, err := u.cartData.GetCartByChatId(ctx, userID)
	if err != nil {
		return domain.UserCart{}, errors.Wrap(err)
	}

	if userCart == nil {
		return domain.UserCart{}, errors.New("no such cart")
	}

	userCart.Cart.Items, err = u.cartData.ListCartItems(ctx, userCart.Cart.ID)
	if err != nil {
		return domain.UserCart{}, errors.Wrap(err)
	}

	return *userCart, nil
}

func New(uD domain.UserRepo, cartData domain.CartRepo) *Service {
	return &Service{
		userData: uD,
		cartData: cartData,
	}
}

func (u *Service) Start(ctx context.Context, userIn domain.User, chatId int64) (domain.StartMessagePayload, error) {
	message, err := u.welcomeMessage(ctx, userIn)
	if err != nil {
		return message, errors.Wrap(err)
	}

	userCart, err := u.cartData.GetByOwnerId(ctx, userIn.Id)
	if err != nil {
		return domain.StartMessagePayload{Msg: domain.DbErrorMsg}, errors.Wrap(err, "error getting cart by owner")
	}

	if userCart != nil {
		message.Cart = userCart.Cart
		message.Cart.Items, err = u.cartData.ListCartItems(ctx, userCart.Cart.ID)
	} else {
		message.Cart, err = u.createCartForUser(ctx, userIn, chatId)
	}
	if err != nil {
		return domain.StartMessagePayload{
			Msg: domain.DbErrorMsg,
		}, errors.Wrap(err)
	}

	message.Msg += fmt.Sprintf(scripts.Get(ctx, scripts.WelcomeMessagePattern), message.Cart.ID)

	return message, nil
}

func (u *Service) AddToDefaultCart(ctx context.Context, items []domain.Item, userID int64,
) (cart domain.UserCart, err error) {
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
		return domain.UserCart{}, errors.Wrap(err)
	}

	if cart.Cart.ID == 0 {
		return domain.UserCart{}, errors.Wrap(ErrNoDefaultCart)
	}

	err = u.cartData.AddCartItems(ctx, items, cart.Cart.ID, userID)
	if err != nil {
		return domain.UserCart{}, errors.Wrap(err)
	}

	cart.Cart.Items, err = u.cartData.ListCartItems(ctx, cart.Cart.ID)
	if err != nil {
		return domain.UserCart{}, errors.Wrap(err, "error getting cartItems")
	}

	return cart, nil
}

func (u *Service) createCartForUser(ctx context.Context, user domain.User, chatId int64) (domain.Cart, error) {
	cartId, err := u.cartData.Create(ctx, user.Id)
	if err != nil {
		return domain.Cart{}, errors.Wrap(err, "error creating cart")
	}

	err = u.cartData.LinkUserToCart(ctx, user.Id, cartId, chatId)
	if err != nil {
		return domain.Cart{}, errors.Wrap(err, "error linking cart")
	}

	err = u.cartData.SetDefaultCart(ctx, user.Id, cartId)
	if err != nil {
		return domain.Cart{}, errors.Wrap(err, "error setting default cart")
	}

	cart, err := u.cartData.GetCartByID(ctx, cartId)
	if err != nil {
		return domain.Cart{}, errors.Wrap(err)
	}

	return cart.Cart, nil
}

func (u *Service) welcomeMessage(ctx context.Context, userIn domain.User) (domain.StartMessagePayload, error) {
	message := domain.StartMessagePayload{}

	dbUser, err := u.userData.Get(ctx, userIn.Id)
	if err != nil {
		return domain.StartMessagePayload{Msg: domain.DbErrorMsg}, errors.Wrap(err, "error creating new dbUser")
	}

	if dbUser != nil {
		message.Msg = scripts.Get(ctx, scripts.WelcomeBack)
		message.User = *dbUser
		return message, nil
	}

	err = u.userData.Upsert(ctx, userIn)
	if err != nil {
		return domain.StartMessagePayload{
			Msg: domain.DbErrorMsg,
		}, errors.Wrap(err, "error updating user's profile")
	}
	message.User = userIn
	message.Msg = scripts.Get(ctx, scripts.Welcome)

	return message, nil
}
