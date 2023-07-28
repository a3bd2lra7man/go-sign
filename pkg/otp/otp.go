package otp

import (
	"github.com/a3bd2lra7man/go-sign/pkg/otp/internal/core"
	"github.com/a3bd2lra7man/go-sign/pkg/otp/internal/dao"
	"go.mongodb.org/mongo-driver/mongo"
)

var manager core.OtpManager

func SetUp(db *mongo.Database) {
	manager = core.New(dao.New(db))
}

func SendOtpForRole(identifier string, role string) error {
	return manager.Create(identifier, role)
}

func CheckOtp(identifier string, role string, code string) error {
	return manager.CheckCode(identifier, role, code)
}
