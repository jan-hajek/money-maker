package mailer

import (
	"strings"

	"time"

	"github.com/jelito/money-maker/src/log"
)

type BufferFactory struct {
	Mailer *Service
	Log    log.Log

	subjects []string
	messages []string
}

func (s *BufferFactory) Create(size int) chan<- BufferItem {
	ch := make(chan BufferItem, size)
	s.Log.Debug("create mail buffer, size:", size)
	s.subjects, s.messages = make([]string, 0), make([]string, 0)

	go func(ch2 <-chan BufferItem) {
		for item := range ch2 {
			s.subjects = append(s.subjects, item.Subject)
			s.messages = append(s.messages, item.Message)
			if len(s.subjects) == size {
				s.Log.Info("buffer is full, sending mails")
				s.sendEmail()
			}
		}
	}(ch)

	go func() {
		t := time.NewTicker(time.Minute * 1)

		for {
			<-t.C

			if len(s.subjects) > 0 {
				s.Log.Info("cycle sending mails")
				s.sendEmail()
			}
		}
	}()

	return ch
}

func (s *BufferFactory) sendEmail() {
	subjects := strings.Join(s.subjects, ",")
	messages := strings.Join(s.messages, "\n")
	s.subjects, s.messages = []string{}, []string{}

	err := s.Mailer.Send(subjects, messages)
	if err != nil {
		s.Log.Error(err)
	}
}

type BufferItem struct {
	Subject string
	Message string
}
