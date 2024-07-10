package service

import (
	"github.com/Red-Sock/Red-Cart/internal/interfaces/data"
	"github.com/Red-Sock/Red-Cart/internal/interfaces/service"
	"github.com/Red-Sock/Red-Cart/internal/service/cart"
	"github.com/Red-Sock/Red-Cart/internal/service/item"
	"github.com/Red-Sock/Red-Cart/internal/service/user"
)

type Storage struct {
	UserService *user.Service
	CartService *cart.Service
	ItemService *item.Service
}

func New(sD data.Storage) *Storage {
	return &Storage{
		UserService: user.New(sD.User(), sD.Cart()),
		CartService: cart.New(sD.Cart()),
		ItemService: item.New(sD.Item()),
	}
}

func (s *Storage) User() service.UserService {
	return s.UserService
}

func (s *Storage) Cart() service.CartService {
	return s.CartService
}

func (s *Storage) Item() service.ItemService {
	return s.ItemService
}
