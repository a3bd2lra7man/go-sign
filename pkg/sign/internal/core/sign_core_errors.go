package core

import "errors"

type SignError uint8

const (
	UnFound         SignError = 0
	PasswordUnMatch SignError = 1
	UnAllowed       SignError = 2
)

func (err SignError) Error() string {
	switch err {
	case UnFound:
		return "UnFound"
	case PasswordUnMatch:
		return "PasswordUnMatch"
	case UnAllowed:
		return "UnAllowed"
	}
	return errors.New("error").Error()
}

type SignUpError uint8

const (
	UserIsExist SignUpError = 0
)

func (err SignUpError) Error() string {
	switch err {
	case UserIsExist:
		return "Exist"
	}
	return errors.New("error").Error()
}

