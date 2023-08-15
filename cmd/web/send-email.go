package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/pradeepj4u/bookings/cmd/models"
	mail "github.com/xhit/go-simple-mail/v2"
)

func listenToMail() {
	go func() {
		for msg := range app.MailChan {
			sendEmail(msg)
		}
	}()
}

func sendEmail(m models.EmailData) {
	server := mail.NewSMTPClient()
	server.Host = "localhost"
	server.Port = 1025
	server.KeepAlive = false
	server.ConnectTimeout = 10 * time.Second
	server.SendTimeout = 10 * time.Second

	client, err := server.Connect()
	if err != nil {
		errorLog.Println(err)
	}
	email := mail.NewMSG()

	email.SetFrom(m.From).AddTo(m.To).SetSubject(m.Subject)
	if m.Templet == "" {
		email.SetBody(mail.TextHTML, m.Content)

	} else {
		data, err := os.ReadFile(fmt.Sprintf("./email-templet/%s", m.Templet))
		if err != nil {
			errorLog.Panicln(err)
		}
		template := string(data)
		msgSent := strings.Replace(template, "[%body%]", m.Content, 1)
		email.SetBody(mail.TextHTML, msgSent)

	}

	err = email.Send(client)
	if err != nil {
		errorLog.Println(err)
	}
	infoLog.Println("Email Sent")
}
