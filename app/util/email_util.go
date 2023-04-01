package util

import (
	"bytes"
	"fmt"
	"github.com/frchandra/ticketing-gmcgo/config"
	"gopkg.in/gomail.v2"
	"html/template"
	"io"
	"strconv"
)

type EmailUtil struct {
	config *config.AppConfig
	log    *LogUtil
}

func NewEmailUtil(config *config.AppConfig, log *LogUtil) *EmailUtil {
	return &EmailUtil{config: config, log: log}
}

func (u *EmailUtil) SendEmail(templatePath string, data map[string]any, receiver string, subject string, attachments map[string][]byte, seatsName []string) error {

	var body bytes.Buffer //prepare template
	t, err := template.ParseFiles(templatePath)
	if err != nil {
		return err
	}

	err = t.Execute(&body, data) //populate template with data
	if err != nil {
		return err
	}

	m := gomail.NewMessage() //create mailer
	m.SetHeader("From", u.config.MailUsername)
	m.SetHeader("To", receiver)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body.String())

	if len(attachments) > 0 { //attaching e-ticket attachment if there is any
		/*		//for filename, _ := range attachments {
				fmt.Println(len(attachments["31.png"]))
				m.Attach(
					"31.png",
					gomail.SetCopyFunc(func(writer io.Writer) error {
						_, err = writer.Write(attachments["31.png"])

						fmt.Println(len(attachments["31.png"]))
						return err
					}),
				)
				fmt.Println(len(attachments["32.png"]))
				m.Attach(
					"32.png",
					gomail.SetCopyFunc(func(writer io.Writer) error {
						_, err = writer.Write(attachments["32.png"])

						fmt.Println(len(attachments["32.png"]))
						return err
					}),
				)
				fmt.Println(len(attachments["33.png"]))
				m.Attach(
					"33.png",
					gomail.SetCopyFunc(func(writer io.Writer) error {
						_, err = writer.Write(attachments["33.png"])

						fmt.Println(len(attachments["33.png"]))
						return err
					}),
				)

				//}
				fmt.Println(">>>>>>>>")*/

		seatsCount := len(seatsName)

		for i := 0; i < seatsCount; i++ {
			filename := seatsName[i] + ".png"
			fmt.Println(len(attachments[filename]))
			m.Attach(
				filename,
				gomail.SetCopyFunc(func(writer io.Writer) error {
					_, err = writer.Write(attachments[filename])

					fmt.Println(len(attachments[filename]))
					return err
				}),
			)
		}

	}

	port, _ := strconv.Atoi(u.config.MailPort) //send the mail
	d := gomail.NewDialer(u.config.MailHost, port, u.config.MailUsername, u.config.MailPassword)
	err = d.DialAndSend(m)
	if err != nil {
		u.log.BasicLog(err, "EmailUtil@SendEmail: when about to sending an email")
		return err
	}

	return nil
}
