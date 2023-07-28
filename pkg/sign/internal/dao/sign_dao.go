package dao

import "github.com/a3bd2lra7man/sign/pkg/sign/internal/entities"

type ISingDao interface {
	GetUser(identifier string) (entities.User, error)
	IsNotExist(identifier string) bool
	CreateUser(identifier string, password []byte) error
}
