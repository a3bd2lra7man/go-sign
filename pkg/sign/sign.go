package sign

import (
	"github.com/a3bd2lra7man/go-sign/pkg/otp"
	"github.com/a3bd2lra7man/go-sign/pkg/roles"
	"github.com/a3bd2lra7man/go-sign/pkg/sign/internal/core"
	"github.com/a3bd2lra7man/go-sign/pkg/sign/internal/dao"
	"go.mongodb.org/mongo-driver/mongo"
)

var manager core.SignManager

func SetUp(db *mongo.Database) {
	manager = core.New(dao.New(db))
	roles.SetUp(db)
}

func SignIn(identifier string, password []byte, role string) error {
	err := manager.Sign(identifier, password)
	if err != nil {
		return err
	}
	err = roles.CheckStatus(identifier, role)

	if err != nil {
		return err
	}

	return nil
}

func CreateUser(identifier string, password []byte) error {
	err := manager.SignUp(identifier, password)

	if err != nil && err != core.UserIsExist {
		return err
	}

	return nil
}

func SignUpForRole(identifier string, password []byte, role string) error {
	err := manager.SignUp(identifier, password)
	if err != nil && err != core.UserIsExist {
		return err
	}

	err = roles.CheckStatus(role, identifier)
	if err != nil {
		return err
	}

	err = otp.SendOtpForRole(identifier, role)
	if err != nil {
		return err
	}

	return nil
}

func ConfirmCode(identifier, role, code string) error {
	err := otp.CheckOtp(identifier, role, code)
	if err != nil {
		return err
	}

	err = roles.Activate(role, identifier)

	if err != nil {
		return err
	}

	return nil
}
