package users

import (
	"context"

	errors "github.com/Red-Sock/trace-errors"
	"github.com/jackc/pgx/v5"

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
	_, err := u.conn.Exec(ctx, `
INSERT INTO tg_users
	    (tg_id)
VALUES	(   $1)`,
		user.Id,
	)
	if err != nil {
		return errors.Wrap(err, "error inserting user to database")
	}

	return nil
}

func (u *Users) Get(ctx context.Context, userId int64) (user.User, error) {
	var dbUser user.User
	err := u.conn.QueryRow(ctx, `
SELECT 
    tg_id
    FROM tg_users
WHERE tg_id = $1`,
		userId,
	).Scan(
		&dbUser.Id,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return dbUser, nil
	}

	if err != nil {
		return user.User{}, errors.Wrap(err, "error getting user from database")
	}

	return dbUser, nil
}
