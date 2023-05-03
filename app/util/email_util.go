package util

import (
	"bytes"
	"github.com/frchandra/ticketing-gmcgo/config"
	"gopkg.in/gomail.v2"
	"html/template"
	"io"
	"math/rand"
	"strconv"
	"time"
)

type EmailUtil struct {
	config *config.AppConfig
	log    *LogUtil
}

func NewEmailUtil(config *config.AppConfig, log *LogUtil) *EmailUtil {
	return &EmailUtil{config: config, log: log}
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

	message := gomail.NewMessage() //create mail message
	message.SetHeader("To", receiver)
	message.SetHeader("Subject", "TICKET EMAIL")
	message.SetBody("text/html", body.String())

	if len(attachments) > 0 {
		seatsCount := len(seatsName)
		for i := 0; i < seatsCount; i++ { //This is you from the past: Do not touch this code. Somehow if I use for-range loop for
			filename := seatsName[i] + ".png" //attaching the attachment map ( map["filename"] consist of filename as a key and []byte of e-ticket image as a value)
			message.Attach(                   // this leads to random behaviour that makes the loop won't iterate the attachment variable inside message.Attach()
				filename, //So, I use regular for loop instead. As a result I need to pass a slice of seatsName as an iterator
				gomail.SetCopyFunc(func(writer io.Writer) error {
					_, err = writer.Write(attachments[filename])
					return err
				}),
			)
		}
	}

	rand.Seed(time.Now().UnixNano()) //sending the email, if fail try multiple times
	index := rand.Intn(4)
	for i := 0; i < 5; i++ {
		err = u.SendEmail(message, u.config.Emails[index])
		if err == nil {
			break
		}
		u.log.Logrus.WithField("occurrence", "sending e-ticket email").WithField("receiver", receiver).WithField("seats_name", seatsName).Error(err)
		index = ((index-1)%4 + 4) % 4 //modulus with positive remainder
	}

	if err != nil {
		return err
	}

	return nil
}

func (u *EmailUtil) SendTotpEmail(data map[string]any, receiver string) error {
	var body bytes.Buffer //prepare template
	t, err := template.ParseFiles("./resource/template/totp.gohtml")
	if err != nil {
		return err
	}

	err = t.Execute(&body, data) //populate template with data
	if err != nil {
		return err
	}

	message := gomail.NewMessage() //create mailer
	message.SetHeader("To", receiver)
	message.SetHeader("Subject", "TOTP EMAIL")
	message.SetBody("text/html", body.String())

	rand.Seed(time.Now().UnixNano()) //sending the email, if fail try multiple times
	index := rand.Intn(4)
	for i := 0; i < 5; i++ {
		err = u.SendEmail(message, u.config.Emails[index])
		if err == nil {
			break
		}
		u.log.Logrus.WithField("occurrence", "sending totp email").WithField("receiver", receiver).Error(err)
		index = ((index-1)%4 + 4) % 4 //modulus with positive remainder
	}

	if err != nil {
		return err
	}

	return nil
}

func (u *EmailUtil) SendEmail(message *gomail.Message, emailAccount config.EmailConfig) error {
	message.SetHeader("From", emailAccount.MailUsername)
	port, _ := strconv.Atoi(emailAccount.MailPort) //send the mail
	dialer := gomail.NewDialer(emailAccount.MailHost, port, emailAccount.MailUsername, emailAccount.MailPassword)
	err := dialer.DialAndSend(message)
	if err != nil {
		return err
	}

	return nil
}
