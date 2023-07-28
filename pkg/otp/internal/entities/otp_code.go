package entities

import "time"

type OtpCode struct {
	Identifier string
	Role       string
	Code       string
	ExpireTime time.Time
}
