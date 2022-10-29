package impl

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/smtp"

	"github.com/jordan-wright/email"

	"github.com/go-sonic/sonic/model/property"
	"github.com/go-sonic/sonic/service"
	"github.com/go-sonic/sonic/util/xerr"
)

type emailServiceImpl struct {
	OptionService service.OptionService
}

func NewEmailService(optionService service.OptionService) service.EmailService {
	return &emailServiceImpl{
		OptionService: optionService,
	}
}

type emailProperties struct {
	Host     string
	Protocol string
	SSLPort  int
	Username string
	Password string
	FromName string
}

func (e *emailServiceImpl) SendTextEmail(ctx context.Context, to, subject, content string) error {
	emailEnable, err := e.OptionService.GetOrByDefaultWithErr(ctx, property.EmailIsEnabled, false)
	if err != nil {
		return err
	}
	if !emailEnable.(bool) {
		return nil
	}
	emailProperties, err := e.getEmailProperties(ctx)
	if err != nil {
		return err
	}
	email := &email.Email{
		To:      []string{to},
		From:    emailProperties.FromName,
		Subject: subject,
		Text:    []byte(content),
	}
	err = e.sendEmail(email, emailProperties)
	return err
}

func (e *emailServiceImpl) SendTemplateEmail(ctx context.Context, to, subject string, content string) error {
	emailEnable, err := e.OptionService.GetOrByDefaultWithErr(ctx, property.EmailIsEnabled, false)
	if err != nil {
		return err
	}
	if !emailEnable.(bool) {
		return nil
	}
	emailProperties, err := e.getEmailProperties(ctx)
	if err != nil {
		return err
	}
	email := &email.Email{
		To:      []string{to},
		From:    emailProperties.FromName,
		Subject: subject,
		HTML:    []byte(content),
	}
	err = e.sendEmail(email, emailProperties)
	return err
}

func (e *emailServiceImpl) TestConnection() error {
	panic("implement me")
}

func (e *emailServiceImpl) getEmailProperties(ctx context.Context) (emailProperties, error) {
	getOrByDefault := func(ctx context.Context, p property.Property, defaultValue interface{}, err error) (interface{}, error) {
		if err != nil {
			return nil, err
		}
		return e.OptionService.GetOrByDefaultWithErr(ctx, p, defaultValue)
	}
	host, err := getOrByDefault(ctx, property.EmailHost, property.EmailHost.DefaultValue, nil)
	protocol, err := getOrByDefault(ctx, property.EmailProtocol, property.EmailProtocol.DefaultValue, err)
	sslPort, err := getOrByDefault(ctx, property.EmailSSLPort, property.EmailSSLPort.DefaultValue, err)
	username, err := getOrByDefault(ctx, property.EmailUsername, property.EmailUsername.DefaultValue, err)
	password, err := getOrByDefault(ctx, property.EmailPassword, property.EmailPassword.DefaultValue, err)
	fromName, err := getOrByDefault(ctx, property.EmailFromName, property.EmailFromName.DefaultValue, err)
	if err != nil {
		return emailProperties{}, err
	}
	emailProperties := emailProperties{
		Host:     host.(string),
		Protocol: protocol.(string),
		SSLPort:  sslPort.(int),
		Username: username.(string),
		Password: password.(string),
		FromName: fromName.(string),
	}
	return emailProperties, nil
}

func (e *emailServiceImpl) sendEmail(email *email.Email, properties emailProperties) error {
	err := email.SendWithTLS(fmt.Sprintf("%s:%d", properties.Host, properties.SSLPort),
		smtp.PlainAuth("", properties.Username, properties.Password, properties.Host), &tls.Config{ServerName: properties.Host, MinVersion: tls.VersionTLS13})
	if err != nil {
		return xerr.Email.Wrapf(err, "发送邮件错误 emailProperties=%v", properties).WithStatus(xerr.StatusInternalServerError).
			WithMsg("发送邮件错误")
	}
	return nil
}
