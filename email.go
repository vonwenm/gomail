// gomail is a simple wrapper around Go's standard net/mail and net/smtp
// based on a gist by andelf: https://gist.github.com/andelf/5004821
package gomail

import (
	"encoding/base64"
	"fmt"
	"log"
	"net/mail"
	"net/smtp"
	"strings"
)

type email struct {
	from    mail.Address
	to      mail.Address
	headers map[string]string
	body    string
}

func NewEmail() *email {
	e := new(email)
	e.headers = make(map[string]string)
	e.Header("From", "")
	e.Header("To", "")
	e.Header("Subject", "")
	e.Header("MIME-Version", "1.0")
	e.Header("Content-Type", "text/plain; charset=\"utf-8\"")
	e.Header("Content-Transfer-Encoding", "base64")
	return e
}

func (e *email) Header(header string, value string) *email {
	e.headers[header] = value
	return e
}

func (e *email) From(address string, name string) *email {
	from := mail.Address{name, address}
	e.from = from
	e.Header("From", from.String())
	return e
}

func (e *email) To(address string, name string) *email {
	to := mail.Address{name, address}
	e.to = to
	e.Header("To", to.String())
	return e
}

func (e *email) Subject(subject string) *email {
	e.Header("Subject", encodeRFC2047(subject))
	return e
}

func (e *email) ReplyTo(address string) *email {
	e.Header("Reply-To", address)
	return e
}

func (e *email) Body(body string) *email {
	e.body = body
	return e
}

func (e *email) Send(server string, auth smtp.Auth) error {
	msg := ""
	for k, v := range e.headers {
		msg += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	msg += "\r\n" + base64.StdEncoding.EncodeToString([]byte(e.body))

	log.Print(msg)
	return smtp.SendMail(
		server,
		auth,
		e.from.Address,
		[]string{e.to.Address},
		[]byte(msg))
}

func encodeRFC2047(String string) string {
	// use mail's rfc2047 to encode any string
	addr := mail.Address{String, ""}
	return strings.Trim(addr.String(), " <>")
}
