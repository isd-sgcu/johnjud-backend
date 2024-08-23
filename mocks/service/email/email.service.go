package email

type Service interface {
	SendEmail(subject string, toName string, toAddress string, content string) error
}
