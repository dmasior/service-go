package mailing

import "os"

type (
	Mailer interface {
		Send(from, to, subject, body string) error
	}
	MailDriver string
	Config     struct {
		MailDriver   string `env:"MAIL_DRIVER,required"`
		SupportEmail string `env:"SUPPORT_EMAIL,required"`
	}
)

const (
	MailDriverSES MailDriver = "ses"
	MailDriverSpy MailDriver = "spy"
)

func SetupFromEnv() Mailer {
	switch driverFromEnv() {
	case MailDriverSES:
		return NewSES(SESConfigFromEnv())
	case MailDriverSpy:
		return NewSpy()
	default:
		panic("unknown mail driver")
	}
}

func driverFromEnv() MailDriver {
	return MailDriver(os.Getenv("MAIL_DRIVER"))
}
