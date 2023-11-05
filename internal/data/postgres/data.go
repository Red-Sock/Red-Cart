package inmemory

import (
	"github.com/Red-Sock/Red-Cart/internal/data/inmemory/carts"
	"github.com/Red-Sock/Red-Cart/internal/data/inmemory/users"
	"github.com/Red-Sock/Red-Cart/internal/interfaces/data"
)

type Storage struct {
	Users data.Users
	Carts data.Carts
}

func New() *Storage {
	return &Storage{
		Users: users.NewUsers(),
		Carts: carts.New()}
}

func (s *Storage) User() data.Users {
	return s.Users
}

func (s *Storage) Cart() data.Carts {
	return s.Carts
}
