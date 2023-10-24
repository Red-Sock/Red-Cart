package service

import (
	"github.com/Red-Sock/Red-Cart/internal/interfaces/data"
	"github.com/Red-Sock/Red-Cart/internal/interfaces/service"
	"github.com/Red-Sock/Red-Cart/internal/service/cart"
	"github.com/Red-Sock/Red-Cart/internal/service/user"
)

type Storage struct {
	UserService     *user.UsersService
	CartService     *cart.CartsService
	CartItemService *cart.CartItemService
}

func New(sD data.Storage) *Storage {
	return &Storage{
		UserService: user.New(sD.User()),
		CartService: cart.New(sD.Cart()),
	}
}

func (s *Storage) CartItem() service.CartItem {
	return s.CartItem()
}

func (s *Storage) User() service.UserService {
	return s.UserService
}

func (s *Storage) Cart() service.CartService {
	return s.CartService
}
