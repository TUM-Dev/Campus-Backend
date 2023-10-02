package cron

import (
	"bytes"
	htmlTemplate "html/template"
	"os"
	"strconv"
	textTemplate "text/template"
	"time"

	"github.com/TUM-Dev/Campus-Backend/server/model"
	log "github.com/sirupsen/logrus"
	"gopkg.in/gomail.v2"

	_ "embed"
)

//go:embed email_templates/feedback_body.gohtml
var htmlFeedbackBody string

//go:embed email_templates/feedback_body.txt.tmpl
var txtFeedbackBody string

// iterate is a template helper to make counting possible
func iterate(count int32) []int32 {
	var items []int32
	var i int32
	for i = 0; i < count; i++ {
		items = append(items, i)
	}
	return items
}

func parseTemplates() (*htmlTemplate.Template, *textTemplate.Template, error) {
	funcMap := textTemplate.FuncMap{"iterate": iterate}
	parsedHtmlBody, err := htmlTemplate.New("htmlFeedbackBody").Funcs(funcMap).Parse(htmlFeedbackBody)
	if err != nil {
		return nil, nil, err
	}
	parsedTxtBody, err := textTemplate.New("txtFeedbackBody").Funcs(funcMap).Parse(txtFeedbackBody)
	if err != nil {
		return nil, nil, err
	}
	return parsedHtmlBody, parsedTxtBody, nil

}

type MailHeaders struct {
	From      string
	To        string
	ReplyTo   string //optional
	Timestamp time.Time
	Subject   string
}

func messageWithHeaders(feedback *model.Feedback) *gomail.Message {
	m := gomail.NewMessage()
	// From
	m.SetAddressHeader("From", os.Getenv("SMTP_USERNAME"), "TUM Campus App")
	// To
	if feedback.Receiver.Valid {
		m.SetHeader("To", feedback.Receiver.String)
	} else {
		m.SetHeader("To", "app@tum.de")
	}
	// ReplyTo
	if feedback.ReplyTo.Valid {
		m.SetHeader("Reply-To", feedback.ReplyTo.String)
	}
	// Timestamp
	if feedback.Timestamp.Valid {
		m.SetDateHeader("Date", feedback.Timestamp.Time)
	} else {
		m.SetDateHeader("Date", time.Now())
	}
	// Subject
	m.SetHeader("Subject", "Feedback via Tum Campus App")
	return m
}

func generateTemplatedMail(parsedHtmlBody *htmlTemplate.Template, parsedTxtBody *textTemplate.Template, feedback *model.Feedback) (string, string, error) {
	var htmlBodyBuffer bytes.Buffer
	if err := parsedHtmlBody.Execute(&htmlBodyBuffer, feedback); err != nil {
		return "", "", err
	}
	var txtBodyBuffer bytes.Buffer
	if err := parsedTxtBody.Execute(&txtBodyBuffer, feedback); err != nil {
		return "", "", err
	}
	return htmlBodyBuffer.String(), txtBodyBuffer.String(), nil
}

func (c *CronService) feedbackEmailCron() error {

	var results []model.Feedback
	if err := c.db.Find(&results, "processed = false").Scan(&results).Error; err != nil {
		log.WithError(err).Fatal("could not get unprocessed feedback")
		return err
	}
	parsedHtmlBody, parsedTxtBody, err := parseTemplates()
	if err != nil {
		log.WithError(err).Fatal("could not parse email templates")
		return err
	}

	dialer, err := setupSMTPDialer()
	if err != nil {
		return err
	}
	for i, feedback := range results {
		m := messageWithHeaders(&feedback)

		// attach a body
		htmlBodyBuffer, txtBodyBuffer, err := generateTemplatedMail(parsedHtmlBody, parsedTxtBody, &feedback)
		if err != nil {
			log.WithError(err).Error("Could not template mail body")
			return err
		}
		m.SetBody("text/plain", txtBodyBuffer)
		m.AddAlternative("text/html", htmlBodyBuffer)

		// send mail
		if err := dialer.DialAndSend(m); err != nil {
			log.WithError(err).Error("could not send mail")
			continue
		}
		log.Tracef("sending feedback %dialer to %v successfull", i, feedback.Receiver)

		// prevent the message being send the next time around
		if err := c.db.Find(model.Feedback{}, "id = ?", feedback.Id).Update("processed", "true").Error; err != nil {
			log.WithError(err).Error("could not prevent mail from being send again")
		}
	}
	return nil
}

// setupSMTPDialer sets up the SMTP dialer
func setupSMTPDialer() (*gomail.Dialer, error) {
	smtpPort, err := strconv.Atoi(os.Getenv("SMTP_PORT"))
	if err != nil {
		log.WithError(err).Fatal("SMTP_PORT is not an integer")
		return nil, err
	}
	d := gomail.NewDialer(os.Getenv("SMTP_URL"), smtpPort, os.Getenv("SMTP_USERNAME"), os.Getenv("SMTP_PASSWORD"))
	return d, nil
}
