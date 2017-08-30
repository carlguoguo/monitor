package email

import (
	"crypto/tls"
	"encoding/base64"
	"errors"
	"fmt"
	"net"
	"net/smtp"
	"strings"
)

var (
	e SMTPInfo
)

// SMTPInfo is the details for the SMTP server
type SMTPInfo struct {
	Username string
	Password string
	Hostname string
	Port     int
	From     string
}

// Configure adds the settings for the SMTP server
func Configure(c SMTPInfo) {
	e = c
}

// ReadConfig returns the SMTP information
func ReadConfig() SMTPInfo {
	return e
}

// LoginAuth is loginauth
type LoginAuth struct {
	username, password string
}

// NewLoginAuth init a loginauth
func NewLoginAuth(username, password string) smtp.Auth {
	return &LoginAuth{username, password}
}

// Start function
func (a *LoginAuth) Start(server *smtp.ServerInfo) (string, []byte, error) {
	return "LOGIN", []byte{}, nil
}

// Next function
func (a *LoginAuth) Next(fromServer []byte, more bool) ([]byte, error) {
	if more {
		switch string(fromServer) {
		case "Username:":
			return []byte(a.username), nil
		case "Password:":
			return []byte(a.password), nil
		default:
			return nil, errors.New("Unknown fromServer")
		}
	}
	return nil, nil
}

// SendMail sends a mail
func SendMail(sendTo string, subject string, content string) error {
	to := strings.Split(sendTo, ";")
	msg := []byte(content)
	auth := NewLoginAuth(e.Username, e.Password)
	addr := fmt.Sprintf("%s:%d", e.Hostname, e.Port)
	c, err := smtp.Dial(addr)
	host, _, _ := net.SplitHostPort(addr)
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer c.Close()

	if ok, _ := c.Extension("STARTTLS"); ok {
		config := &tls.Config{ServerName: host, InsecureSkipVerify: true}
		if err = c.StartTLS(config); err != nil {
			fmt.Println("call start tls")
			return err
		}
	}

	if auth != nil {
		if ok, _ := c.Extension("AUTH"); ok {
			if err = c.Auth(auth); err != nil {
				fmt.Println("check auth with err:", err)
				return err
			}
		}
	}
	if err = c.Mail(e.From); err != nil {
		return err
	}
	for _, addr := range to {
		if err = c.Rcpt(addr); err != nil {
			return err
		}
	}
	w, err := c.Data()

	if err != nil {
		return err
	}

	header := make(map[string]string)
	header["Subject"] = subject
	header["MIME-Version"] = "1.0"
	header["Content-Type"] = "text/plain; charset=\"utf-8\""
	header["Content-Transfer-Encoding"] = "base64"
	message := ""
	for k, v := range header {
		message += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	message += "\r\n" + base64.StdEncoding.EncodeToString(msg)
	_, err = w.Write([]byte(message))

	if err != nil {
		return err
	}
	err = w.Close()
	if err != nil {
		return err
	}
	return c.Quit()
}
