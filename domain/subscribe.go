package domain

import (
	"strings"
)

type Subscribe struct {
	From string `json:"from"`
	To   string `json:"to"`
	id   int64
}
type Subscribes map[string]*Subscribe

func NewSubscribe(from, to string, id int64) *Subscribe {
	return &Subscribe{
		From: strings.ToUpper(from),
		To:   strings.ToUpper(to),
		id:   id,
	}
}

func (s *Subscribe) Id() int64 {
	return s.id
}
