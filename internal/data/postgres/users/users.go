package users

import (
	"context"

	"github.com/Red-Sock/Red-Cart/internal/clients/postgres"
	"github.com/Red-Sock/Red-Cart/internal/domain/user"
)

type Users struct {
	conn postgres.Conn
}

func New(conn postgres.Conn) *Users {
	return &Users{conn: conn}
}

func (u *Users) Upsert(ctx context.Context, user user.User) error {
	// TODO
	return nil
}
