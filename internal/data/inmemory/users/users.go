package users

import (
	"context"
	"github.com/Red-Sock/Red-Cart/internal/domain/user"
)

type Users struct {
	m map[int64]user.User
}

func (u2 Users) Upsert(ctx context.Context, u user.User) error {
	return nil
}

func (u2 Users) Get(ctx context.Context, id int64) (user.User, error) {
	return user.User{}, nil
}

func NewUsers() *Users {
	return &Users{
		m: make(map[int64]user.User),
	}
}
