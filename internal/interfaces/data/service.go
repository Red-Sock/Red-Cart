package data

import (
	"github.com/Red-Sock/Red-Cart/internal/domain"
)

type Storage interface {
	User() domain.UserRepo
	Cart() domain.CartRepo
	Item() domain.ItemRepo
}
