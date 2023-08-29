package cron

import (
	"bytes"
	"crypto/tls"
	"github.com/TUM-Dev/Campus-Backend/server/model"
	log "github.com/sirupsen/logrus"
	"gopkg.in/gomail.v2"
	htmlTemplate "html/template"
	"os"
	"strconv"
	textTemplate "text/template"
	"time"
)
import _ "embed"

// Iterate is necessary, as go otherwise cannot count up in a for loop inside templates
func Iterate(count int32) []int32 {
	var Items []int32
	var i int32
	for i = 0; i < count; i++ {
		Items = append(Items, i)
	}
	return Items
}

//go:embed emailTemplates/feedbackBody.gohtml
var htmlFeedbackBody string

//go:embed emailTemplates/feedbackBody.txt.tmpl
var txtFeedbackBody string

func (c *CronService) feedbackEmailCron() error {
	var results []model.Feedback
	if err := c.db.Find(&results, "processed = false").Scan(&results).Error; err != nil {
		log.WithError(err).Fatal("could not get unprocessed feedback")
		return err
	}
	funcMap := textTemplate.FuncMap{"Iterate": Iterate}
	parsedHtmlBody, err := htmlTemplate.New("htmlFeedbackBody").Funcs(funcMap).Parse(htmlFeedbackBody)
	if err != nil {
		log.WithError(err).Fatal("htmlFeedbackBody is not a valid template")
		return err
	}
	parsedTxtBody, err := textTemplate.New("txtFeedbackBody").Funcs(funcMap).Parse(txtFeedbackBody)
	if err != nil {
		log.WithError(err).Fatal("txtFeedbackBody is not a valid template")
		return err
	}

	smtpPort, err := strconv.Atoi(os.Getenv("SMTP_PORT"))
	if err != nil {
		log.WithError(err).Fatal("SMTP_PORT is not an integer")
		return err
	}
	d := gomail.NewDialer(os.Getenv("SMTP_URL"), smtpPort, os.Getenv("SMTP_USERNAME"), os.Getenv("SMTP_PASSWORD"))
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}
	for i, feedback := range results {
		m := gomail.NewMessage()
		// set message-headers
		m.SetAddressHeader("From", os.Getenv("SMTP_USERNAME"), "TUM Campus App")
		if feedback.Receiver.Valid {
			m.SetHeader("To", feedback.Receiver.String)
		} else {
			m.SetHeader("To", "app@tum.de")
		}
		if feedback.ReplyTo.Valid {
			m.SetHeader("Reply-To", feedback.ReplyTo.String)
		}
		if feedback.Timestamp.Valid {
			m.SetDateHeader("Date", feedback.Timestamp.Time)
		} else {
			m.SetDateHeader("Date", time.Time{})
		}
		m.SetHeader("Subject", "Feedback via Tum Campus App")

		// attach a body
		var txtBodyBuffer bytes.Buffer
		if err := parsedTxtBody.Execute(&txtBodyBuffer, feedback); err != nil {
			return err
		}
		m.SetBody("text/plain", txtBodyBuffer.String())

		var htmlBodyBuffer bytes.Buffer
		if err := parsedHtmlBody.Execute(&htmlBodyBuffer, feedback); err != nil {
			return err
		}
		m.AddAlternative("text/html", htmlBodyBuffer.String())

		// send mail
		if err := d.DialAndSend(m); err != nil {
			log.WithError(err).Error("could not send mail")
			continue
		}
		log.Trace("sending feedback %d to %s successfull", i, feedback.Receiver)

		// prevent the message being send the next time around
		if err := c.db.Find(model.Feedback{}, "id = ?", feedback.Id).Update("processed", "true").Error; err != nil {
			log.WithError(err).Error("could not prevent mail from being send again")
		}
	}
	return nil
}
