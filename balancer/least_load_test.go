package balancer

import (
	fibHeap "github.com/starwander/GoFibonacciHeap"
	"github.com/stretchr/testify/assert"
	"testing"
)

// 这里需要提前生成host,实时生成会导致地址不一样而断言失败
var l1 = &host{name: "http://127.0.0.1:9501", load: 0}
var l2 = &host{name: "http://127.0.0.2:9501", load: 0}
var l3 = &host{name: "http://127.0.0.3:9501", load: 0}
var l4 = &host{name: "http://127.0.0.4:9501", load: 1}

func TestLeastLoad_Add(t *testing.T) {

	var heap1 = fibHeap.NewFibHeap()
	_ = heap1.InsertValue(l1)
	_ = heap1.InsertValue(l2)
	_ = heap1.InsertValue(l3)

	cases := []struct {
		name string
		lb   Balancer

		hostName string   // 添加的host
		expected Balancer // 期望的结果
	}{
		{
			name: "test_least_load_add_1",
			lb: &LeastLoad{
				heap: heap1,
			},
			hostName: "http://127.0.0.1:9501",
			expected: &LeastLoad{
				heap: heap1,
			},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			c.lb.Add(c.hostName)
			assert.Equal(t, c.expected, c.lb)
		})
	}
}

func TestLeastLoad_Balance(t *testing.T) {
	balancer := NewLeastLoad([]string{"http://127.0.0.1:9501", "http://127.0.0.1:9502", "http://127.0.0.1:9503", "http://127.0.0.1:9504"})

	balancer.Remove("http://127.0.0.1:9504")
	balancer.Inc("http://127.0.0.1:9502")
	balancer.Inc("http://127.0.0.1:9502")
	balancer.Inc("http://127.0.0.1:9501")
	balancer.Inc("http://127.0.0.1:9503")
	balancer.Done("http://127.0.0.1:9501")
	h, _ := balancer.Balance("")

	assert.Equal(t, "http://127.0.0.1:9501", h)
}
