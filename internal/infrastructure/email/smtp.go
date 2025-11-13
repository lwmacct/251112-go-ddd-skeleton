package email

import (
	"fmt"
	"net/smtp"
)

// Config SMTP配置
type Config struct {
	Host     string
	Port     int
	Username string
	Password string
	From     string
}

// SMTPSender SMTP邮件发送器
type SMTPSender struct {
	config Config
}

// NewSMTPSender 创建SMTP发送器
func NewSMTPSender(config Config) *SMTPSender {
	return &SMTPSender{config: config}
}

// Send 发送邮件
func (s *SMTPSender) Send(to, subject, body string) error {
	auth := smtp.PlainAuth("", s.config.Username, s.config.Password, s.config.Host)

	msg := []byte(fmt.Sprintf("To: %s\r\n"+
		"Subject: %s\r\n"+
		"\r\n"+
		"%s\r\n", to, subject, body))

	addr := fmt.Sprintf("%s:%d", s.config.Host, s.config.Port)
	return smtp.SendMail(addr, auth, s.config.From, []string{to}, msg)
}
