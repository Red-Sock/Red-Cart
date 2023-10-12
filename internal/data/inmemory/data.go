package inmemory

import (
	"context"

	"github.com/Red-Sock/Red-Cart/internal/domain/user"
	"github.com/Red-Sock/Red-Cart/internal/interfaces/service"
)

type Users struct {
	m map[int64]user.User
}

type Carts struct {
	m map[int64]user.User
}

type Service struct {
	Users Users
	Cart  Carts
}

func (s *Service) User() service.UserService {
	//TODO implement me
	panic("implement me")
}

func (s *Service) CartService() service.CartService {
	//TODO implement me
	panic("implement me")
}

func New() *Service {
	return &Service{
		Users: Users{
			m: make(map[int64]user.User),
		},
		Cart: Carts{
			m: make(map[int64]user.User),
		},
	}
}

func (s *Service) Upsert(ctx context.Context, user user.User) error {

	s.Users.m[user.Id] = user
	return nil
}

func (s *Service) Get(ctx context.Context, id int64) (user.User, error) {

	return s.Users.m[id], nil
}
