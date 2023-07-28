package core

import "errors"

type OtpError uint8

const (
	UnMatch OtpError = 1
	Expired OtpError = 2
)

func (err OtpError) Error() string {
	switch err {
	case UnMatch:
		return "UnMatch"
	case Expired:
		return "Expired"
	}
	return errors.New("error").Error()
}
