package email

import (
	"log"
	"os"
	"strconv"

	"gopkg.in/gomail.v2"
)

type Client interface {
	Send(subject string, message string, destination string) error
}

type Email struct {
}

func New() Client {
	return &Email{}
}

func (s *Email) Send(subject string, message string, destination string) error {

	host := os.Getenv("EMAIL_SMTP_HOST")
	port, err := strconv.Atoi(os.Getenv("EMAIL_SMTP_PORT"))
	if err != nil {
		log.Print(err)
		log.Print("fail to send email")
	}
	username := os.Getenv("EMAIL_SMTP_EMAIL")
	password := os.Getenv("EMAIL_SMTP_PASSWORD")

	mailer := gomail.NewMessage()
	mailer.SetHeader("From", "test@gmail.com")
	mailer.SetHeader("To", destination)
	mailer.SetHeader("Subject", subject)
	mailer.SetBody("text/html", message)

	// err = sendMail([]string{"kenichixie@gmail.com"}, []string{}, subject, message, username, password, host, port)
	// if err != nil {
	// 	log.Print(err)
	// }

	dialer := gomail.NewDialer(
		host, port, username, password,
	)

	err = dialer.DialAndSend(mailer)
	if err != nil {
		log.Print(err)
		log.Print("fail to send email 3")
	}

	return nil
}

// func sendMail(to []string, cc []string, subject string, message string, email string, password string, host string, port int) error {
// 	body := "From: phantomdev.sindomas@gmail.com \n" +
// 		"To: " + strings.Join(to, ",") + "\n" +
// 		"Cc: " + strings.Join(cc, ",") + "\n" +
// 		"Subject: " + subject + "\n\n" +
// 		message

// 	auth := smtp.PlainAuth("", email, password, host)
// 	smtpAddr := fmt.Sprintf("%s:%d", host, port)

// 	err := smtp.SendMail(smtpAddr, auth, email, append(to, cc...), []byte(body))
// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }

func (s *Email) SendTest(phoneNumber string, message string) error {
	log.Printf("sending Email to: %s with message: %s", phoneNumber, message)
	return nil
}
