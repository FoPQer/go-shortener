package service

import "crypto/rand"

func NewId() string {
	return rand.Text()[0:8]
}
