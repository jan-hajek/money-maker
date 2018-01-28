package mailer

import (
	"github.com/jelito/money-maker/app/log"
	"strings"
)

type BufferFactory struct {
	Mailer *Service
	Log    log.Log
}

func (s *BufferFactory) Create(size int) chan<- BufferItem {
	ch := make(chan BufferItem, size)
	s.Log.Debug("create mail buffer, size:", size)

	go func(<-chan BufferItem) {
		subjects, messages := make([]string, 0), make([]string, 0)

		for item := range ch {
			subjects = append(subjects, item.Subject)
			messages = append(messages, item.Message)
			if len(subjects) == size {
				s.Log.Debug("flush mail buffer, sending mails")
				err := s.Mailer.Send(strings.Join(subjects, ","), strings.Join(messages, "\n"))
				if err != nil {
					s.Log.Error(err)
				}

				subjects, messages = []string{}, []string{}
			}
		}
	}(ch)

	return ch
}

type BufferItem struct {
	Subject string
	Message string
}
