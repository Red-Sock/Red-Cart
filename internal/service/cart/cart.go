package cart

import (
	"context"
	"encoding/json"

	tgapi "github.com/Red-Sock/go_tg/interfaces"
	errors "github.com/Red-Sock/trace-errors"

	"github.com/Red-Sock/Red-Cart/internal/domain"
)

type Service struct {
	cartData domain.CartRepo
}

func New(cartData domain.CartRepo) *Service {
	return &Service{
		cartData: cartData,
	}
}

func (c *Service) SyncCartMessage(ctx context.Context, userCart domain.UserCart, msg tgapi.MessageOut) error {
	if userCart.Cart.MessageId != nil && *userCart.Cart.MessageId == msg.GetMessageId() {
		return nil
	}

	userCart.Cart.ChatId = msg.GetChatId()

	msgID := msg.GetMessageId()
	userCart.Cart.MessageId = &msgID

	err := c.cartData.UpdateCartReference(ctx, userCart)
	if err != nil {
		return errors.Wrap(err, "error updating cart chat reference")
	}

	return nil
}

func (c *Service) GetCartByChatId(ctx context.Context, chatID int64) (domain.UserCart, error) {
	userCart, err := c.cartData.GetCartByChatId(ctx, chatID)
	if err != nil {
		return domain.UserCart{}, errors.Wrap(err, "error reading cart for chat")
	}

	userCart.Cart.Items, err = c.cartData.ListCartItems(ctx, userCart.Cart.ID)
	if err != nil {
		return domain.UserCart{}, errors.Wrap(err, "error reading cart items")
	}

	return *userCart, nil
}

func (c *Service) Add(ctx context.Context, items []domain.Item, cartID int64, userID int64) (domain.UserCart, error) {
	cart, err := c.cartData.GetCartByID(ctx, cartID)
	if err != nil {
		return domain.UserCart{}, err
	}

	if cart == nil {
		return domain.UserCart{}, errors.New("no such cart")
	}

	err = c.cartData.AddCartItems(ctx, items, cartID, userID)
	if err != nil {
		return domain.UserCart{}, errors.Wrap(err, "error adding items to cart ")
	}

	cart.Cart.Items, err = c.cartData.ListCartItems(ctx, cartID)
	if err != nil {
		return domain.UserCart{}, errors.Wrap(err, "error getting actual cart items")
	}

	return *cart, err
}

func (c *Service) GetCartById(ctx context.Context, cartID int64) (domain.UserCart, error) {
	uc, err := c.cartData.GetCartByID(ctx, cartID)
	if err != nil {
		return domain.UserCart{}, err
	}

	if uc == nil {
		return domain.UserCart{}, errors.New("Cart doesn't exists")
	}

	uc.Cart.Items, err = c.cartData.ListCartItems(ctx, cartID)
	if err != nil {
		return domain.UserCart{}, errors.Wrap(err, "Can't get items for cart")
	}

	return *uc, nil
}

func (c *Service) AwaitNameChange(ctx context.Context, cartID int64, item domain.Item) (err error) {
	req := domain.Cart{
		ID:    cartID,
		State: domain.CartStateEditingItemName,
	}

	req.StatePayload, err = json.Marshal(domain.ChangeItemNamePayload{ItemName: item.Name})
	err = c.cartData.ChangeState(ctx, req)
	if err != nil {
		return errors.Wrap(err, "error changing state of cart")
	}

	return nil
}

func (c *Service) AwaitItemsAdded(ctx context.Context, cartID int64) (err error) {
	req := domain.Cart{
		ID:    cartID,
		State: domain.CartStateAdding,
	}

	err = c.cartData.ChangeState(ctx, req)
	if err != nil {
		return errors.Wrap(err, "error changing state of cart")
	}

	return nil
}

func (c *Service) PurgeCart(ctx context.Context, cartId int64) error {
	return c.cartData.PurgeCart(ctx, cartId)
}
