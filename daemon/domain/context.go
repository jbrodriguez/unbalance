package domain

import "github.com/cskr/pubsub"

type Context struct {
	Config

	Port string
	Hub  *pubsub.PubSub
}
