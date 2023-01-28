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

func (u *EmailUtil) SendGomail(templatePath string, data map[string]any, reciever string) error {
	var body bytes.Buffer
	t, err := template.ParseFiles(templatePath)
	if err != nil {
		return err
	}
	err = t.Execute(&body, data)
	if err != nil {
		return err
	}

	m := gomail.NewMessage()
	m.SetHeader("From", u.config.MailUsername)
	m.SetHeader("To", "nismara.chandra@gmail.com")
	m.SetHeader("Subject", "Hello!")
	m.SetBody("text/html", body.String())
	m.Attach("./storage/picture/polite_cat.jpg")

	port, _ := strconv.Atoi(u.config.MailPort)
	d := gomail.NewDialer(u.config.MailHost, port, u.config.MailUsername, u.config.MailPassword)

	err = d.DialAndSend(m)
	if err != nil {
		fmt.Println("mail not sent!")
		return err
	}
	fmt.Println("mail sent!")
	return nil
}
