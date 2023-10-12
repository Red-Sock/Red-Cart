package service

import "context"

type UserService interface {
	Start(ctx context.Context, id int64) (message string, err error)
}
