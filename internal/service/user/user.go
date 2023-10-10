package user

import (
	"context"

	"github.com/Red-Sock/Red-Cart/internal/interfaces/data"
)

type UsersService struct {
	ctx      context.Context
	userData data.Users
}

func New(uD data.Users) *UsersService {
	return &UsersService{
		userData: uD,
	}
}

func (u *UsersService) Start(id int64) (message string, err error) {
	user, err := u.userData.Get(u.ctx, id)
	if err != nil {
		return "", err
	}

	if user.Id != 0 {
		return "Welcome Back!", nil
	}
	user.Id = id
	err = u.userData.Upsert(u.ctx, user)

	if err != nil {
		return "", err
	}
	return "Hello New User!", nil
}
