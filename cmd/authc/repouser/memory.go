package repouser

import (
	"context"
	"errors"
)

type Memory struct {
	data []UserInfo
}

func (m Memory) GetUserByEmail(ctx context.Context, email string) (user UserInfo, err error) {
	for _, u := range m.data {
		if u.Email == email {
			return u, nil
		}
	}

	return UserInfo{}, errors.New("user not found")
}

func NewMemory(data []UserInfo) Memory {
	return Memory{data: data}
}

var _ Repository = (*Memory)(nil)
