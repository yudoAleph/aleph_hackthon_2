package coreclient

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetStatusFail(t *testing.T) {
	client := NewClient("http://localhost:9999")
	_, err := client.GetStatus("/test")
	assert.Error(t, err)
}
