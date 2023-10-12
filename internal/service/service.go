package service

import (
	"github.com/Red-Sock/Red-Cart/internal/interfaces/data"
	"github.com/Red-Sock/Red-Cart/internal/interfaces/service"
	"github.com/Red-Sock/Red-Cart/internal/service/cart"
	"github.com/Red-Sock/Red-Cart/internal/service/user"
)

type Service struct {
	User service.UserService
	Cart service.CartService
}

func (s Service) Start(id int64) (message string, err error) {
	msg, err := s.User.Start(id)

	return msg, err
}

func New(sD data.Service) *Service {
	return &Service{
		User: user.New(),
		Cart: cart.New(),
	}
}
