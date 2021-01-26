package mail

import (
	"log"
	"net/smtp"
	"strings"
	"time"
)

type Emailer struct {
	FromEmail    string
	FromPassword string
	Host         string
	Port         string
}

func NewEmailer(fromEmail string, fromPassword string, host string, port string) *Emailer {
	if fromEmail == "" || fromPassword == "" || host == "" || port == "" {
		log.Fatalln("missing email arguments")
	}

	return &Emailer{
		FromEmail:    fromEmail,
		FromPassword: fromPassword,
		Host:         host,
		Port:         port,
	}
}

func (e *Emailer) ConstructMessage(content string, toEmail string) []byte {
	// Subject
	subject := "RSS Feeds for " + time.Now().Format(time.RFC822)

	// Message
	msg := strings.Builder{}
	msg.WriteString("From: " + e.FromEmail + "\n")
	msg.WriteString("To: " + toEmail + "\n")
	msg.WriteString("Subject: " + subject + "\n")
	msg.WriteString("MIME-version: 1.0;\n")
	msg.WriteString("Content-Type: text/html;charset=\"UTF-8\";\n")
	msg.WriteString("\n")
	msg.WriteString(content)

	return []byte(msg.String())
}

func (e *Emailer) Send(content string, toEmail string) error {
	// Auth
	auth := smtp.PlainAuth("", e.FromEmail, e.FromPassword, e.Host)
	msgBytes := e.ConstructMessage(content, toEmail)

	// Send
	err := smtp.SendMail(e.Host+":"+e.Port, auth, e.FromEmail, []string{toEmail}, msgBytes)
	if err != nil {
		log.Println(err)
		return err
	}

	log.Println("email sent successfully")
	return nil
}
