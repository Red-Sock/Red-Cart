package cart

import (
	"github.com/Red-Sock/Red-Cart/internal/interfaces/data"
)

type CartItemService struct {
	cartsData data.Carts
}

func NewCartItem() *CartItemService {
	return &CartItemService{}
}
