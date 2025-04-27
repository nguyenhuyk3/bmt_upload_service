package services

import "context"

type IMessageBrokerWriter interface {
	SendMessage(ctx context.Context, topic, key string, message interface{}) error
	Close()
}
