package user

import (
	"context"

	"github.com/Red-Sock/Red-Cart/internal/domain/user"
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

func (u *UsersService) Start(ctx context.Context, newUser user.User) (message string, err error) {
	user, err := u.userData.Get(ctx, newUser.Id)
	if err != nil {
		return "", err
	}

	if user.Id != 0 {
		return "Welcome Back!", nil
	}
	user = newUser
	err = u.userData.Upsert(ctx, user)

	if err != nil {
		return "", err
	}
	return "Hello New User!", nil
}
