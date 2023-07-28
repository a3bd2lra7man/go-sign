package core

import (
	errs "github.com/a3bd2lra7man/sign/pkg/err"
	"github.com/a3bd2lra7man/sign/pkg/sign/internal/dao"
	"golang.org/x/crypto/bcrypt"
)

type SignManager struct {
	dao dao.ISingDao
}

func New(dao dao.ISingDao) SignManager {
	return SignManager{dao: dao}
}

func (m SignManager) Sign(identifier string, password []byte) error {
	user, err := m.dao.GetUser(identifier)
	if err != nil {
		return UnFound
	}

	check := bcrypt.CompareHashAndPassword([]byte(user.Password), password)

	if check != nil {
		return PasswordUnMatch
	}

	if !user.Allowed {
		return UnAllowed
	}

	return nil
}

func (m SignManager) SignUp(identifier string, password []byte) error {
	IsNotExist := m.dao.IsNotExist(identifier)
	if !IsNotExist {
		return UserIsExist
	}

	password, err := bcrypt.GenerateFromPassword([]byte(password), 10)

	if err != nil {
		return errs.UnExpectedError{}
	}

	err = m.dao.CreateUser(identifier, password)

	if err != nil {
		return errs.UnExpectedError{}
	}

	return nil
}
