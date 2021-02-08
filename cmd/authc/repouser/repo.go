package repouser

import (
	"context"
)

type Repository interface {
	GetUserByEmail(ctx context.Context, email string) (user UserInfo, err error)
}

type UserInfo struct {
	ID       int
	Email    string
	Password string
}
