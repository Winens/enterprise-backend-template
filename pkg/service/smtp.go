package service

import (
	"bytes"
	"context"
	"crypto/tls"
	"fmt"
	"runtime"
	"time"

	"github.com/knadh/smtppool"
	"github.com/spf13/viper"
	"github.com/winens/enterprise-backend-template/pkg/service/interfaces"
	"github.com/winens/enterprise-backend-template/pkg/templates"
)

func init() {
	viper.SetDefault("smtp.settings.idle_timeout", time.Second*10)
	viper.SetDefault("smtp.settings.pool_wait_timeout", time.Second*3)
	viper.SetDefault("smtp.settings.max_conns", runtime.NumCPU()*2)
}

type smtpService struct {
	smtp *smtppool.Pool
}

func NewSMTPService() (interfaces.SMTPService, error) {

	opts := smtppool.Opt{
		Host: viper.GetString("smtp.host"),
		Port: viper.GetInt("smtp.port"),

		MaxConns:        viper.GetInt("smtp.settings.max_conns"),
		IdleTimeout:     viper.GetDuration("smtp.settings.idle_timeout"),
		PoolWaitTimeout: viper.GetDuration("smtp.settings.pool_wait_timeout"),
	}

	// set start_tls to true to enable TLS.
	if viper.GetBool("smtp.start_tls") {
		opts.TLSConfig = &tls.Config{
			InsecureSkipVerify: viper.GetBool("smtp.insecure_skip_verify"),
			ServerName:         viper.GetString("smtp.host"),
		}
	}

	// set username and password if provided.
	if viper.GetString("smtp.username") != "" {
		opts.Auth = &smtppool.LoginAuth{
			Username: viper.GetString("smtp.username"),
			Password: viper.GetString("smtp.password"),
		}
	}

	smtp, err := smtppool.New(opts)

	if err != nil {
		return nil, err
	}

	return &smtpService{
		smtp: smtp,
	}, nil
}

func (s *smtpService) SendUserVerificationEmail(email, firstName, token string) error {
	link := viper.GetString("website.url") + "/auth/verify?token=" + token

	var buf bytes.Buffer
	if err := templates.UserVerificationEmail(firstName, link).Render(context.Background(), &buf); err != nil {
		return fmt.Errorf("template rendering failed: %w", err)
	}

	e := smtppool.Email{
		From:    viper.GetString("smtp.from"),
		To:      []string{email},
		Subject: fmt.Sprintf("Welcome %s! Please verify your email address", firstName),
		HTML:    buf.Bytes(),
	}

	return s.smtp.Send(e)
}
