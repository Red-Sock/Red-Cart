package users

import (
	"context"

	"github.com/Red-Sock/Red-Cart/internal/domain/user"
)

type Users struct {
	m map[int64]user.User
}

func NewUsers() *Users {
	return &Users{
		m: make(map[int64]user.User),
	}
}

func (u *Users) Upsert(ctx context.Context, user user.User) error {
	u.m[user.Id] = user
	return nil
}

func (u *Users) Get(ctx context.Context, id int64) (user.User, error) {

	return u.m[id], nil
}
