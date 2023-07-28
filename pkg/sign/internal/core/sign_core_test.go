package core

import (
	"errors"
	"testing"

	"github.com/a3bd2lra7man/go-sign/pkg/sign/internal/entities"
	"golang.org/x/crypto/bcrypt"
)

type mockDao struct {
	user        entities.User
	err         error
	_isNotExist bool
}

func (m *mockDao) GetUser(id string) (entities.User, error) {
	return m.user, m.err
}

func (m *mockDao) IsNotExist(id string) bool {
	return m._isNotExist
}

func (m *mockDao) CreateUser(id string, Password []byte) error {
	return m.err
}

func createManager(dao *mockDao) *SignManager {
	if dao.user.Password != "" {
		pass, _ := bcrypt.GenerateFromPassword([]byte(dao.user.Password), 10)
		dao.user.Password = string(pass)
	}
	return &SignManager{
		dao: dao,
	}
}

func TestSignIn(t *testing.T) {

	tests := []struct {
		name       string
		dao        mockDao
		Identifier string
		Password   string
		err        error
		id         string
	}{
		{
			name: "user not found in db",
			dao:  mockDao{user: entities.User{}, err: errors.New("any error")},
			err:  UnFound,
		},
		{
			name:       "Password doesn't match",
			Identifier: "test@test.com",
			Password:   "87654321",
			dao:        mockDao{user: entities.User{Identifier: "test@test.com", Password: "12345678"}},
			err:        PasswordUnMatch,
		},
		{
			name:       "user is blocked",
			Identifier: "test@test.com",
			Password:   "12345678",
			dao:        mockDao{user: entities.User{Identifier: "test@test.com", Password: "12345678", Allowed: false}},
			err:        UnAllowed,
		},
		{
			name:       "success",
			Identifier: "test@test.com",
			Password:   "12345678",
			id:         "idFromDB",
			dao:        mockDao{user: entities.User{Identifier: "test@test.com", Password: "12345678", Allowed: true}},
			err:        nil,
		},
	}

	for _, test := range tests {
		signManager := createManager(&test.dao)
		err := signManager.Sign(test.Identifier, []byte(test.Password))
		if err != test.err {
			t.Fatalf("test %s expected error : %v found : %v", test.name, test.err, err)
			continue
		}
	}
}

func TestSignUp(t *testing.T) {

	tests := []struct {
		name       string
		Identifier string
		Password   string
		expectedId string
		err        error
		dao        mockDao
	}{
		{
			name:       "user is exist in db",
			Identifier: "test@test.com",
			Password:   "87654321",
			dao:        mockDao{_isNotExist: false},
			err:        UserIsExist,
		},
		{
			name:       "success",
			Identifier: "test@test.com",
			Password:   "12345678",
			dao:        mockDao{_isNotExist: true},
			expectedId: "expected_id",
			err:        nil,
		},
	}

	for _, test := range tests {
		signManager := createManager(&test.dao)
		err := signManager.SignUp(test.Identifier, []byte(test.Password))
		if err != test.err {
			t.Fatalf("test %s expected error : %v found : %v", test.name, test.err, err)
		}
	}
}
