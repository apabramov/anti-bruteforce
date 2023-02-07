package storage

import (
	"errors"
)

var (
	ErrExists    = errors.New("subnet already exists")
	ErrNotExists = errors.New("subnet not exists")
)

type Authorize struct {
	Login    string
	Password string
	IP       string
}

func NewAuthorize(login string, password string, ip string) Authorize {
	return Authorize{
		Login:    login,
		Password: password,
		IP:       ip,
	}
}
