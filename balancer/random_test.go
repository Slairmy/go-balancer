package balancer

import (
	"github.com/stretchr/testify/assert"
	"math/rand"
	"testing"
	"time"
)

func TestRandomAdd(t *testing.T) {
	rand1 := rand.New(rand.NewSource(time.Now().UnixNano()))

	// 单元测试实际对比的是 Balancer
	cases := []struct {
		name     string
		lb       Balancer
		host     string
		expected Balancer
	}{
		{
			name: "test_add_1",
			lb: &Random{
				hosts: []string{"http://127.0.0.1: 9501", "http://127.0.0.1: 9502"},
				rnd:   rand1,
			},
			host: "http://127.0.0.1: 9503",
			expected: &Random{
				hosts: []string{"http://127.0.0.1: 9501", "http://127.0.0.1: 9502", "http://127.0.0.1: 9503"},
				rnd:   rand1,
			},
		},
		{
			name: "test_add_2",
			lb: &Random{
				hosts: []string{"http://127.0.0.1: 9501", "http://127.0.0.1: 9502"},
				rnd:   rand1,
			},
			host: "http://127.0.0.1: 9504",
			expected: &Random{
				hosts: []string{"http://127.0.0.1: 9501", "http://127.0.0.1: 9502", "http://127.0.0.1: 9504"},
				rnd:   rand1,
			},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			c.lb.Add(c.host)
			assert.Equal(t, c.expected, c.lb)
		})
	}
}

func TestRandomRemove(t *testing.T) {

	rand1 := rand.New(rand.NewSource(time.Now().UnixNano()))

	// 单元测试实际对比的是 Balancer
	cases := []struct {
		name     string
		lb       Balancer
		host     string
		expected Balancer
	}{
		{
			name: "test_remove_1",
			lb: &Random{
				hosts: []string{"http://127.0.0.1: 9501", "http://127.0.0.1: 9502", "http://127.0.0.1: 9503"},
				rnd:   rand1,
			},
			host: "http://127.0.0.1: 9503",
			expected: &Random{
				hosts: []string{"http://127.0.0.1: 9501", "http://127.0.0.1: 9502"},
				rnd:   rand1,
			},
		},
		{
			name: "test_remove_2",
			lb: &Random{
				hosts: []string{"http://127.0.0.1: 9501", "http://127.0.0.1: 9502", "http://127.0.0.1: 9504"},
				rnd:   rand1,
			},
			host: "http://127.0.0.1: 9504",
			expected: &Random{
				hosts: []string{"http://127.0.0.1: 9501", "http://127.0.0.1: 9502"},
				rnd:   rand1,
			},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			c.lb.Remove(c.host)
			assert.Equal(t, c.expected, c.lb)
		})
	}
}

// 我以为的单元测试 -- 缺少了拿取不到host的情况
func TestRandomBalancer(t *testing.T) {
	rand1 := rand.New(rand.NewSource(time.Now().UnixNano()))

	lb := &Random{
		hosts: []string{"http://127.0.0.1: 9501", "http://127.0.0.1: 9502", "http://127.0.0.1: 9503"},
		rnd:   rand1,
	}

	selected, _ := lb.Balance("")
	assert.Contains(t, lb.hosts, selected)
}

func TestRandomBalancer1(t *testing.T) {
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
