package balancer

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestConsistentHash_Balance(t *testing.T) {

	type expected struct {
		reply string // 最终返回的主机
		err   error
	}

	cases := []struct {
		name     string
		clientIP string
		lb       Balancer
		expected expected
	}{
		{
			name:     "test_consistent_hash_balance_1",
			clientIP: "http://192.168.1.1",
			lb:       NewConsistentHash([]string{"http://127.0.0.1:9501", "http://127.0.0.2:9501", "http://127.0.0.3:9501"}),
			expected: expected{reply: "http://127.0.0.1:9501", err: nil},
		},
		{
			name:     "test_consistent_hash_balance_2",
			clientIP: "http://192.168.1.1",
			lb:       NewConsistentHash([]string{}),
			expected: expected{reply: "", err: NoHostError},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			h, err := c.lb.Balance(c.clientIP)
			assert.Equal(t, c.expected.reply, h)
			assert.Equal(t, c.expected.err, err)
		})
	}

}
