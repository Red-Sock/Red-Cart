package cart

import (
	"github.com/Red-Sock/Red-Cart/internal/interfaces/data"
)

type Cart struct {
	cartsData data.Carts
}

func New(userData data.Carts) *Cart {
	return &Cart{
		cartsData: userData,
	}
}
