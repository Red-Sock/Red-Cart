package users

import (
	"context"
	"sync"

	"github.com/Red-Sock/Red-Cart/internal/domain/user"
)

type Users struct {
	rw sync.RWMutex
	m  map[int64]user.User
}

func NewUsers() *Users {
	return &Users{
		m: make(map[int64]user.User),
	}
}

func (u *Users) Upsert(ctx context.Context, user user.User) error {
	u.rw.Lock()
	u.m[user.Id] = user
	defer u.rw.Unlock()
	return nil
}

func (u *Users) Get(ctx context.Context, id int64) (user.User, error) {

	return u.m[id], nil
}
