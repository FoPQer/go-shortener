package service

import "crypto/rand"

func NewID() string {
	return rand.Text()[0:8]
}
