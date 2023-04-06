package util

import (
	"bytes"
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

func (u *EmailUtil) SendInfoEmail(data map[string]any, receiver string) error {
	var body bytes.Buffer //prepare template
	t, err := template.ParseFiles("./resource/template/info.gohtml")
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
	m.SetHeader("Subject", "INFO EMAIL")
	m.SetBody("text/html", body.String())

	port, _ := strconv.Atoi(u.config.MailPort) //send the mail
	d := gomail.NewDialer(u.config.MailHost, port, u.config.MailUsername, u.config.MailPassword)
	err = d.DialAndSend(m)
	if err != nil {
		u.log.BasicLog(err, "EmailUtil@SendInfoEmail: when about to sending an email")
		return err
	}

	return nil
}

func (u *EmailUtil) SendTicketEmail(data map[string]any, receiver string, attachments map[string][]byte, seatsName []string) error {

	var body bytes.Buffer //prepare template
	t, err := template.ParseFiles("./resource/template/ticket.gohtml")
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
	m.SetHeader("Subject", "TICKET EMAIL")
	m.SetBody("text/html", body.String())

	if len(attachments) > 0 {
		seatsCount := len(seatsName)
		for i := 0; i < seatsCount; i++ { //This is you from the past: Do not touch this code. Somehow if I use for-range loop for
			filename := seatsName[i] + ".png" //attaching the attachment map ( map["filename"] consist of filename as a key and []byte of e-ticket image as a value)
			m.Attach(                         // this leads to random behaviour that makes the loop won't iterate the attachment variable inside m.Attach()
				filename, //So, I use regular for loop instead. As a result I need to pass a slice of seatsName as an iterator
				gomail.SetCopyFunc(func(writer io.Writer) error {
					_, err = writer.Write(attachments[filename])
					return err
				}),
			)
		}
	}

	port, _ := strconv.Atoi(u.config.MailPort) //send the mail
	d := gomail.NewDialer(u.config.MailHost, port, u.config.MailUsername, u.config.MailPassword)
	err = d.DialAndSend(m)
	if err != nil {
		u.log.BasicLog(err, "EmailUtil@SendTicketEmail: when about to sending an email")
		return err
	}

	return nil
}

func (u *EmailUtil) SentTotpEmail(data map[string]any, receiver string) error {
	var body bytes.Buffer //prepare template
	t, err := template.ParseFiles("./resource/template/totp.gohtml")
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
	m.SetHeader("Subject", "TOTP EMAIL")
	m.SetBody("text/html", body.String())

	port, _ := strconv.Atoi(u.config.MailPort) //send the mail
	d := gomail.NewDialer(u.config.MailHost, port, u.config.MailUsername, u.config.MailPassword)
	err = d.DialAndSend(m)
	if err != nil {
		u.log.BasicLog(err, "EmailUtil@SendTotpEmail: when about to sending an email")
		return err
	}

	return nil
}
