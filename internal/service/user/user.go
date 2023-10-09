package user

import (
	"context"
	"github.com/Red-Sock/Red-Cart/internal/interfaces/data"
	"log"
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

func (u *UsersService) Start(id int64) (message string) {
	user, err := u.userData.Get(u.ctx, id)
	if err != nil {
		//TODO обработка ошибки
	}

	if user.Id != 0 {
		return "Welcome Back!"
	}
	user.Id = id
	err = u.userData.Add(u.ctx, user)

	if err != nil {
		log.Fatal("ошибка добавления пользователя: ", err)
	}
	return "Hello New User!"
}
