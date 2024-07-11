package domain

import (
	"context"
)

type UserRepo interface {
	Upsert(ctx context.Context, u User) error
	Get(ctx context.Context, id int64) (*User, error)
}

type User struct {
	Id        int64
	UserName  string
	FirstName string
	LastName  string
}
