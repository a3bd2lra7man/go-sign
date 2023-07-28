//go:build integration
// +build integration

//
package core

import "testing"

func TestSendMail(t *testing.T) {
	err := sendMail("a3bd2lra7man@gmail.com", "1234")

	if err != nil {
		t.Fatalf("Failed to send mail because %v", err)
	}
}
