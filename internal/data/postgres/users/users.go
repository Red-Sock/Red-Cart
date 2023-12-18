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
	    (tg_id,
	     user_name,
	     first_name,
	     last_name)
VALUES	($1,
        $2,
        $3,
        $4)`,
		user.Id,
		user.UserName,
		user.FirstName,
		user.LastName,
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
    tg_id,
    user_name,
    first_name,
    last_name
    FROM tg_users
WHERE tg_id = $1`,
		userId,
	).Scan(
		&dbUser.Id,
		&dbUser.UserName,
		&dbUser.FirstName,
		&dbUser.LastName,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return dbUser, nil
	}

	if err != nil {
		return user.User{}, errors.Wrap(err, "error getting user from database")
	}

	return dbUser, nil
}
