package inmemory

import (
	"github.com/Red-Sock/Red-Cart/internal/data/inmemory/carts"
	"github.com/Red-Sock/Red-Cart/internal/data/inmemory/users"
	"github.com/Red-Sock/Red-Cart/internal/domain/user"
	"github.com/Red-Sock/Red-Cart/internal/interfaces/data"
)

type Carts struct {
	m map[int64]user.User
}

type Storage struct {
	Users *users.Users
	Carts *carts.Carts
}

func New() *Storage {
	return &Storage{Users: users.NewUsers(),
		Carts: carts.New()}
}

func (s *Storage) User() data.Users {
	return s.Users
}

func (s *Storage) Cart() data.Carts {
	return s.Carts
}
