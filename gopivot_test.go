package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetConfig(t *testing.T) {
	config := GetConfig()
	assert.NotEqual(t, config.Token, "", "Expected token to have some length to it")
}
