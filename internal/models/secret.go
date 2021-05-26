package models

import (
	"errors"
)
// Secret defines model for secret
type Secret struct {
	Key string
	Value string
}

// SecretRESTPost defines model for secret for the rest
type SecretRESTPost struct {
	Key string `json:"key"`
	Value string `json:"value"`
	Secret string `json:"secret"`
}

// Validate provides validating of the data
func (s *SecretRESTPost) Validate() error {
	if s.Key == "" {
		return errors.New("key is not defined")
	}
	if s.Value == "" {
		return errors.New("value is not defined")
	}
	return nil
}