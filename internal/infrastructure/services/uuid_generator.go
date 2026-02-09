package services

import (
	"github.com/google/uuid"
)

type UUIDv7Generator struct{}

func NewUUIDv7Generator() (*UUIDv7Generator, error) {
	return &UUIDv7Generator{}, nil
}

func (g *UUIDv7Generator) Generate() string {
	return uuid.Must(uuid.NewV7()).String()
}
