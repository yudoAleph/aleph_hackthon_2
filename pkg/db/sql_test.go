package db

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewSQLConnection(t *testing.T) {
	_, err := NewSQLConnection("invalid-dsn")
	assert.Error(t, err)
}
