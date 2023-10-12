package users

import (
	"context"
	"github.com/Red-Sock/Red-Cart/internal/domain/user"
)

type Users struct {
	m map[int64]user.User
}

func (u2 Users) Upsert(ctx context.Context, u user.User) error {
	//TODO implement me
	panic("implement me")
}

func (u2 Users) Get(ctx context.Context, id int64) (user.User, error) {
	//TODO implement me
	panic("implement me")
}

func NewUsers() *Users {
	return &Users{
		m: make(map[int64]user.User),
	}
}
