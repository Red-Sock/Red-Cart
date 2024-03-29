package postgres

import (
	"github.com/Red-Sock/Red-Cart/internal/clients/postgres"
	"github.com/Red-Sock/Red-Cart/internal/data/postgres/carts"
	"github.com/Red-Sock/Red-Cart/internal/data/postgres/items"
	"github.com/Red-Sock/Red-Cart/internal/data/postgres/users"
	"github.com/Red-Sock/Red-Cart/internal/domain"
)

type Storage struct {
	Users domain.UserRepo
	Carts domain.CartRepo
	Items domain.ItemRepo
}

func New(conn postgres.Conn) *Storage {
	return &Storage{
		Users: users.New(conn),
		Carts: carts.New(conn),
		Items: items.New(conn),
	}
}

func (s *Storage) User() domain.UserRepo {
	return s.Users
}

func (s *Storage) Cart() domain.CartRepo {
	return s.Carts
}

func (s *Storage) Item() domain.ItemRepo {
	return s.Items
}
