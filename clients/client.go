package clients

import (
	"github.com/streamdp/ccd/domain"
)

// RestClient interface makes it possible to expand the list of rest data providers
type RestClient interface {
	Get(from string, to string) (*domain.Data, error)
}

// WsClient interface makes it possible to expand the list of wss data providers
type WsClient interface {
	Subscribe(from string, to string) error
	Unsubscribe(from string, to string) error
	ListSubscribes() domain.Subscribes
}
