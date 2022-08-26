package balancer

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestIpHashBalancer(t *testing.T) {
	type expected struct {
		host string
		err  error
	}

	cases := []struct {
		name     string
		lb       Balancer
		args     string
		expected expected
	}{
		{
			name: "test_balance_1",
			lb:   NewRandom([]string{"http://127.0.0.1: 9501"}),
			args: "",
			expected: expected{
				host: "http://127.0.0.1: 9501",
				err:  nil,
			},
		},
		{
			name: "test_balance_2",
			lb:   NewRandom([]string{}),
			args: "h",
			expected: expected{
				host: "",
				err:  NoHostError,
			},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			host, err := c.lb.Balance(c.args)
			assert.Equal(t, c.expected.host, host)
			assert.Equal(t, c.expected.err, err)
		})
	}

}
