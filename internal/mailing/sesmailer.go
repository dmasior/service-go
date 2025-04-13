package mailing

import (
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ses"
)

type (
	SESConfig struct {
		AccessKeyID string `env:"AWS_ACCESS_KEY_ID"`
		SecretKey   string `env:"AWS_SECRET_KEY"`
		Region      string `env:"AWS_REGION"`
	}

	SES struct {
		config SESConfig
	}
)

func SESConfigFromEnv() SESConfig {
	return SESConfig{
		Region:      os.Getenv("AWS_REGION"),
		AccessKeyID: os.Getenv("AWS_ACCESS_KEY_ID"),
		SecretKey:   os.Getenv("AWS_SECRET_KEY"),
	}
}

func NewSES(config SESConfig) *SES {
	return &SES{
		config: config,
	}
}

func (m *SES) Send(from, to, subject, body string) error {
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(m.config.Region),
		Credentials: credentials.NewStaticCredentials(m.config.AccessKeyID, m.config.SecretKey, ""),
	},
	)
	if err != nil {
		return err
	}

	svc := ses.New(sess)

	input := &ses.SendEmailInput{
		Destination: &ses.Destination{
			ToAddresses: []*string{
				aws.String(to),
			},
		},
		Message: &ses.Message{
			Body: &ses.Body{
				Text: &ses.Content{
					Data: aws.String(body),
				},
			},
			Subject: &ses.Content{
				Data: aws.String(subject),
			},
		},
		Source: aws.String(from),
	}

	_, err = svc.SendEmail(input)
	if err != nil {
		return err
	}

	return nil
}
