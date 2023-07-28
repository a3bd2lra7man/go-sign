package core

import (
	"bytes"
	"math/rand"
	"time"

	"github.com/a3bd2lra7man/sign/pkg/otp/internal/dao"
)

type OtpManager struct {
	dao dao.IOtpDao
}

func New(dao dao.IOtpDao) OtpManager {
	return OtpManager{dao: dao}
}

func createOtp() string {
	allowedKeys := "1234567890"

	var buffer bytes.Buffer

	for i := 0; i < 5; i++ {
		index := rand.Intn(10)
		buffer.WriteByte(allowedKeys[index])
	}

	return buffer.String()
}

func (m OtpManager) Create(identifier string, role string) error {
	code := createOtp()
	err := sendMail(identifier, code)
	if err != nil {
		return err
	}
	err = m.dao.Save(identifier, code, role)
	if err != nil {
		return err
	}
	return nil
}

func (m OtpManager) CheckCode(identifier string, role string, code string) error {
	codeInDb, err := m.dao.GetCode(identifier, role)
	if err != nil {
		return err
	}
	if codeInDb.Code != code {
		return UnMatch
	}
	if time.Now().Compare(codeInDb.ExpireTime) >= 0 {
		return Expired
	}
	return nil
}
