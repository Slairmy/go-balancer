package balancer

import (
	"github.com/stretchr/testify/assert"
	"math/rand"
	"testing"
	"time"
)

var rnd = rand.New(rand.NewSource(time.Now().UnixNano()))

// 这里需要提前生成host,实时生成会导致地址不一样而断言失败
var h1 = &host{name: "http://127.0.0.1:9501", load: 0}
var h2 = &host{name: "http://127.0.0.2:9501", load: 0}
var h3 = &host{name: "http://127.0.0.3:9501", load: 0}
var h4 = &host{name: "http://127.0.0.4:9501", load: 1}

func TestP2c_Add(t *testing.T) {

	cases := []struct {
		name string
		lb   Balancer

		hostName string   // 添加的host
		expected Balancer // 期望的结果
	}{
		{
			name: "test_p2c_add_1",
			lb: &P2c{
				hosts:   []*host{h1, h2},
				rnd:     rnd,
				loadMap: map[string]*host{"http://127.0.0.1:9501": h1, "http://127.0.0.2:9501": h2}},
			hostName: "http://127.0.0.2:9501",
			expected: &P2c{
				hosts:   []*host{h1, h2},
				rnd:     rnd,
				loadMap: map[string]*host{"http://127.0.0.1:9501": h1, "http://127.0.0.2:9501": h2},
			},
		},
		{
			name: "test_p2c_add_2",
			lb: &P2c{
				hosts:   []*host{h1, h2},
				loadMap: map[string]*host{"http://127.0.0.1:9501": h1, "http://127.0.0.2:9501": h2}},
			hostName: "http://127.0.0.3:9501",
			expected: &P2c{
				hosts:   []*host{h1, h2, h3},
				loadMap: map[string]*host{"http://127.0.0.1:9501": h1, "http://127.0.0.2:9501": h2, "http://127.0.0.3:9501": h3},
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

func TestP2c_Balance(t *testing.T) {
	// 这里是断言返回的目标资源服务器
	type excepted struct {
		err   error
		reply string
	}

	c := struct {
		name     string
		lb       Balancer
		key      string
		excepted excepted
	}{
		name: "test_p2c_balancer",
		lb:   NewP2c([]string{"http://127.0.0.1:9501", "http://127.0.0.2:9501"}),
		key:  "key1",
		excepted: excepted{
			err:   nil,
			reply: "http://127.0.0.1:9501",
		},
	}

	t.Run(c.name, func(t *testing.T) {
		// inc
		c.lb.Inc(c.key)
		actual, err := c.lb.Balance(c.key)
		assert.Equal(t, c.excepted.reply, actual)
		assert.Equal(t, c.excepted.err, err)
	})

}

func TestP2c_Inc(t *testing.T) {
	lb := &P2c{
		hosts:   []*host{h1},
		rnd:     rnd,
		loadMap: map[string]*host{"http://127.0.0.1:9501": h1},
	}

	hostName := "http://127.0.0.1:9501"

	t.Run("test_p2c_inc", func(t *testing.T) {
		lb.Inc(hostName)
		assert.Equal(t, uint64(1), lb.loadMap[hostName].load)
	})

}

func TestP2c_Done(t *testing.T) {
	lb := &P2c{
		hosts:   []*host{h1},
		rnd:     rnd,
		loadMap: map[string]*host{"http://127.0.0.4:9501": h4},
	}

	hostName := "http://127.0.0.4:9501"

	t.Run("test_p2c_done", func(t *testing.T) {
		lb.Done(hostName)
		assert.Equal(t, uint64(0), lb.loadMap[hostName].load)
	})
}
