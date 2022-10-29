package service

import (
	"context"
)

type EmailService interface {
	SendTextEmail(ctx context.Context, to, subject, content string) error
	SendTemplateEmail(ctx context.Context, to, subject string, content string) error
	TestConnection() error
}
