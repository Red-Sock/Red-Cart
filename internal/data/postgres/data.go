package postgres

import (
	"github.com/Red-Sock/Red-Cart/internal/clients/postgres"
	"github.com/Red-Sock/Red-Cart/internal/data/postgres/carts"
	"github.com/Red-Sock/Red-Cart/internal/data/postgres/users"
	"github.com/Red-Sock/Red-Cart/internal/interfaces/data"
)

type Storage struct {
	Users data.Users
	Carts data.Carts
}

func New(conn postgres.Conn) *Storage {
	return &Storage{
		Users: users.New(conn),
		Carts: carts.New(conn),
	}
}

func (s *Storage) User() data.Users {
	return s.Users
}

func (s *Storage) Cart() data.Carts {
	return s.Carts
}
