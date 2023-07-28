package core

import (
	"testing"
	"time"

	errs "github.com/a3bd2lra7man/sign/pkg/err"
	"github.com/a3bd2lra7man/sign/pkg/otp/internal/entities"
)

type mockDao struct {
	err        error
	code       string
	expireTime time.Time
}

func (d mockDao) GetCode(identifier, role string) (entities.OtpCode, error) {
	return entities.OtpCode{Code: d.code, ExpireTime: d.expireTime}, d.err

}

func (d mockDao) Save(identifier, code, role string) error {

	return nil
}

func createManager(dao *mockDao) *OtpManager {
	return &OtpManager{
		dao: dao,
	}
}

func TestCreateOtp(t *testing.T) {
	otp := createOtp()

	if len(otp) != 5 {
		t.Fatalf("expected otp to be four character length found : %v", otp)
	}
}

func TestCheckCode(t *testing.T) {

	tests := []struct {
		name          string
		dao           mockDao
		code          string
		expectedError error
	}{
		{
			name:          "code not found in db",
			dao:           mockDao{err: errs.UnExpectedError{}},
			expectedError: errs.UnExpectedError{},
		},
		{
			name:          "code found in db but not match",
			dao:           mockDao{err: nil, code: "code"},
			code:          "mis_match_code",
			expectedError: UnMatch,
		},
		{
			name:          "code found in db but expired",
			dao:           mockDao{err: nil, code: "code", expireTime: time.Now().Add(-time.Hour)},
			code:          "code",
			expectedError: Expired,
		},
		{
			name:          "success",
			dao:           mockDao{err: nil, code: "code", expireTime: time.Now().Add(time.Minute)},
			code:          "code",
			expectedError: nil,
		},
	}
	for _, test := range tests {
		manager := createManager(&test.dao)
		err := manager.CheckCode("", "", test.code)
		if err != test.expectedError {
			t.Fatalf("test %s expected error : %v found : %v", test.name, test.expectedError, err)
		}

	}
}
