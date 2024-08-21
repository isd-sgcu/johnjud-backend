package email

import (
	"fmt"
	"net/http"

	"github.com/isd-sgcu/johnjud-gateway/config"
	"github.com/pkg/errors"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

type Service interface {
	SendEmail(subject string, toName string, toAddress string, content string) error
}

type serviceImpl struct {
	config config.Sendgrid
	client *sendgrid.Client
}

func NewService(config config.Sendgrid) Service {
	client := sendgrid.NewSendClient(config.ApiKey)
	return &serviceImpl{config: config, client: client}
}

func (s *serviceImpl) SendEmail(subject string, toName string, toAddress string, content string) error {
	from := mail.NewEmail(s.config.Name, s.config.Address)
	to := mail.NewEmail(toName, toAddress)
	message := mail.NewSingleEmail(from, subject, to, content, content)

	resp, err := s.client.Send(message)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusAccepted {
		return errors.New(fmt.Sprintf("%d status code", resp.StatusCode))
	}

	return nil
}
