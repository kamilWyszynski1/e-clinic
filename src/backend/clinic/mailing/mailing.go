package mailing

import (
	mailjet "github.com/mailjet/mailjet-apiv3-go"
)

type Mailer interface {
	SendEmail(email, name, subject, content string) (*mailjet.ResultsV31, error)
}

type Client struct {
	cli *mailjet.Client
}

func NewClient(cli *mailjet.Client) *Client {
	return &Client{cli: cli}
}

func (c Client) SendEmail(email, name, subject, content string) (*mailjet.ResultsV31, error) {
	messagesInfo := []mailjet.InfoMessagesV31{
		mailjet.InfoMessagesV31{
			From: &mailjet.RecipientV31{
				Email: "c-clinic@sinoqwerty.com",
				Name:  "e-clinic",
			},
			To: &mailjet.RecipientsV31{
				mailjet.RecipientV31{
					Email: email,
					Name:  name,
				},
			},
			Subject:  subject,
			TextPart: "My first Mailjet email",
			HTMLPart: content,
			CustomID: "AppGettingStartedTest",
		},
	}
	messages := mailjet.MessagesV31{Info: messagesInfo}
	res, err := c.cli.SendMailV31(&messages)
	if err != nil {
		return nil, err
	}
	return res, nil
}
