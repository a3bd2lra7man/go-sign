package core

import "errors"

type RolesErrors uint8

const (
	UnActive RolesErrors = 0
)

func (err RolesErrors) Error() string {
	switch err {
	case UnActive:
		return "UnActive"
	}
	return errors.New("error").Error()
}
