package user

import (
	"context"

	"github.com/Red-Sock/Red-Cart/internal/interfaces/data"
)

type UsersService struct {
	userData data.Users
}

func New(uD data.Users) *UsersService {
	return &UsersService{
		userData: uD,
	}
}

func (u *UsersService) Start(ctx context.Context, id int64) (message string, err error) {
	user, err := u.userData.Get(ctx, id)
	if err != nil {
		return "", err
	}

	if user.Id != 0 {
		return "Welcome Back!", nil
	}
	user.Id = id
	err = u.userData.Upsert(ctx, user)

	if err != nil {
		return "", err
	}
	return "Hello New User!", nil
}
