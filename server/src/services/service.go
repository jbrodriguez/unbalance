package services

import (
	"github.com/jbrodriguez/pubsub"
)

// MailboxHandler -
type MailboxHandler func(msg *pubsub.Message)

// Service -
type Service struct {
	registry map[string]MailboxHandler
}

func (s *Service) init() {
	s.registry = make(map[string]MailboxHandler)
}

func (s *Service) register(bus *pubsub.PubSub, topic string, handler MailboxHandler) (mbox chan *pubsub.Mailbox) {
	mbox = bus.Sub(topic)
	s.registry[topic] = handler
	return mbox
}

func (s *Service) registerAdditional(bus *pubsub.PubSub, topic string, handler MailboxHandler, mb chan *pubsub.Mailbox) {
	bus.AddSub(mb, topic)
	s.registry[topic] = handler
}

func (s *Service) dispatch(topic string, msg *pubsub.Message) {
	s.registry[topic](msg)
}
