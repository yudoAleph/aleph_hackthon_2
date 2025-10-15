package redis

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPingRedisFail(t *testing.T) {
	client := NewRedisClient("localhost:9999", "", 0)
	err := PingRedis(client)
	assert.Error(t, err)
}
