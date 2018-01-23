package mailer

import (
	"net/smtp"
)

func Create(enabled bool, addr, from, pass, to, host string) *Service {
	return &Service{
		enabled: enabled,
		addr:    addr,
		from:    from,
		pass:    pass,
		to:      to,
		host:    host,
	}
}

type Service struct {
	enabled bool
	addr    string
	from    string
	pass    string
	to      string
	host    string
}

func (s *Service) Send(subject, message string) error {
	if s.enabled == false {
		return nil
	}

	msg := "From: " + s.from + "\n" +
		"To: " + s.to + "\n" +
		"Subject: " + subject + "\n\n" +
		message

	return smtp.SendMail(
		s.addr,
		smtp.PlainAuth("", s.from, s.pass, s.host),
		s.from,
		[]string{s.to},
		[]byte(msg),
	)
}
