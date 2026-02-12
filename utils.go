package main

import (
	"math/rand"
	"time"
)

const characters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
const urlLength = 8

func init() {
	rand.Seed(time.Now().UnixNano())
}

func generateShortURL() string {
	b := make([]byte, urlLength)
	for i := range b {
		b[i] = characters[rand.Intn(len(characters))]
	}
	return string(b)
}
