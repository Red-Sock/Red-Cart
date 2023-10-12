package data

import (
	"github.com/Red-Sock/Red-Cart/internal/interfaces/service"
)

type Service interface {
	User() service.UserService
	CartService() service.CartService
}
