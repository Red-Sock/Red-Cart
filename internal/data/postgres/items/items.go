package items

import (
	"context"

	errors "github.com/Red-Sock/trace-errors"

	"github.com/Red-Sock/Red-Cart/internal/clients/postgres"
)

type Repository struct {
	conn postgres.Conn
}

func New(conn postgres.Conn) *Repository {
	return &Repository{
		conn: conn,
	}
}

func (s *Repository) ChangeName(ctx context.Context, cartID int64, oldName, newName string) error {
	_, err := s.conn.Exec(ctx, `
		UPDATE cart_items 
		SET item_name = $1
		WHERE cart_id = $2 
		AND item_name = $3
`, newName, cartID, oldName)
	if err != nil {
		return errors.Wrap(err, "error changing name in db")
	}

	return nil
}
