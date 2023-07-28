
package core

import (
	"crypto/tls"
	"fmt"
	"net/smtp"
)

func sendMail(to string, otp string) error {

	// Sender data.
	from := "a3bd2llah@gmail.com"
	password := "oeskmjbjsocrwtzb"

	// smtp server configuration.
	host := "smtp.gmail.com"
	port := "465"

	// Authentication.
	auth := smtp.PlainAuth("", from, password, host)

	tlsConfig := &tls.Config{
		InsecureSkipVerify: true,
		ServerName:         host,
	}

	conn, err := tls.Dial("tcp", fmt.Sprintf("%s:%s", host, port), tlsConfig)
	if err != nil {
		return err
	}

	client, err := smtp.NewClient(conn, host)
	if err != nil {
		return err
	}

	if err = client.Auth(auth); err != nil {
		return err
	}

	if err = client.Mail(from); err != nil {
		return err
	}

	// Message.
	if err = client.Rcpt(to); err != nil {
		return err
	}

	writer, err := client.Data()
	if err != nil {
		return err
	}

	html := fmt.Sprintf(`<h1>Email Confirmation</h1>
		<h2>Hello %s</h2>
		<p>Thank you for subscribing. Please copy and paste this code</p>
		<a> %s</a>


		</div>`, "name", otp)

	message := fmt.Sprintf("To: %s\r\n"+
		"Subject: M7 ECommerce\r\n"+
		"Content-Type: text/html; charset=UTF-8\r\n"+
		"\r\n"+
		"%s\r\n", to, html)

	_, err = writer.Write([]byte(message))
	if err != nil {
		return err
	}

	if err = writer.Close(); err != nil {
		return err
	}

	client.Quit()

	return nil

}
