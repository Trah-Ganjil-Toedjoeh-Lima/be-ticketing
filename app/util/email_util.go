package util

import (
	"bytes"
	"fmt"
	"github.com/frchandra/gmcgo/config"
	"gopkg.in/gomail.v2"
	"html/template"
	"strconv"
)

type EmailUtil struct {
	config *config.AppConfig
}

func NewEmailUtil(config *config.AppConfig) *EmailUtil {
	return &EmailUtil{config: config}
}

func (u *EmailUtil) SendEmail(templatePath string, data map[string]any, receiver string, subject string, attachementPath []string) error {
	//prepare template
	var body bytes.Buffer
	t, err := template.ParseFiles(templatePath)
	if err != nil {
		return err
	}
	//populate template with data
	err = t.Execute(&body, data)
	if err != nil {
		return err
	}
	//create mailer
	m := gomail.NewMessage()
	m.SetHeader("From", u.config.MailUsername)
	m.SetHeader("To", receiver)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body.String())
	m.Attach("./storage/picture/polite_cat.jpg")

	port, _ := strconv.Atoi(u.config.MailPort)
	d := gomail.NewDialer(u.config.MailHost, port, u.config.MailUsername, u.config.MailPassword)

	//send mail
	err = d.DialAndSend(m)
	if err != nil {
		fmt.Println("mail not sent!")
		return err
	}
	fmt.Println("mail sent!")
	return nil
}
