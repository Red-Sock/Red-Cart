package user

import (
	"context"

	"github.com/Red-Sock/Red-Cart/internal/domain"
)

type UsersService struct {
	userData domain.UserRepo
}

func New(uD domain.UserRepo) *UsersService {
	return &UsersService{
		userData: uD,
	}
}

func (u *UsersService) Start(ctx context.Context, newUser domain.User) (message string, err error) {
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
