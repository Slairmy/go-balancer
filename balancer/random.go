package balancer

import (
	"math/rand"
	"sync"
	"time"
)

type Random struct {
	sync.RWMutex

	hosts []string
	rnd   *rand.Rand
}

func init() {
	factories["random"] = NewRandom
}

func NewRandom(hosts []string) Balancer {
	// 这样子写和直接使用 rand.Intn() 有很大的性能影响吗?
	// 根据文档说的 默认的Source是并发安全的，并发安全意味着肯定有锁限制,这里使用自定义的Source是否是因为Bu需要并发安全这一点性能考虑?
	return &Random{hosts: hosts, rnd: rand.New(rand.NewSource(time.Now().UnixNano()))}
}

func (r *Random) Add(host string) {
	r.RLock()
	defer r.RUnlock()

	for _, h := range r.hosts {
		if h == host {
			return
		}
	}

	r.hosts = append(r.hosts, host)
}

func (r *Random) Remove(host string) {
	r.RLock()
	defer r.RUnlock()
	for i, h := range r.hosts {
		if h == host {
			r.hosts = append(r.hosts[:i], r.hosts[i+1:]...)
		}
	}
}

func (r *Random) Balance(_ string) (string, error) {
	// 随机选取一个
	r.Lock()
	defer r.Unlock()

	if len(r.hosts) == 0 {
		return "", NoHostError
	}

	//host := r.hosts[rand.Intn(len(r.hosts))]
	host := r.hosts[r.rnd.Intn(len(r.hosts))]

	return host, nil
}
func (r *Random) Inc(_ string)  {}
func (r *Random) Done(_ string) {}
