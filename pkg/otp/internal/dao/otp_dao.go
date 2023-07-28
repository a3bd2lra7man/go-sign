package dao

import "github.com/a3bd2lra7man/go-sign/pkg/otp/internal/entities"

type IOtpDao interface {
	Save(identifier, code, role string) error
	GetCode(identifier, role string) (entities.OtpCode, error)
}
