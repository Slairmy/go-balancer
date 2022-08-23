package balancer

import (
	"errors"
	"sync"
)

// 负载均衡 -- 轮询算法

var (
	NoHostError = errors.New("no host error")
)

type RoundRobin struct {
	sync.RWMutex
	i     uint64
	hosts []string
}

// 加入工厂算法集 -- balancer.factories
func init() {
	factories["round-robin"] = NewRoundRobin
}

// NewRoundRobin 初始化把主机映射加进去
func NewRoundRobin(hosts []string) Balancer {
	return &RoundRobin{i: 0, hosts: hosts}
}

// 先实现 Balancer interface

func (r *RoundRobin) Add(host string) {
	r.Lock()
	defer r.Unlock()
	// 如果请求的目标host在切片数组中直接返回否则加入代理集
	for _, h := range r.hosts {
		if h == host {
			return
		}
	}

	r.hosts = append(r.hosts, host)
}

func (r *RoundRobin) Remove(host string) {
	r.Lock()
	defer r.Unlock()
	for i, h := range r.hosts {
		if h == host {
			r.hosts = append(r.hosts[:i], r.hosts[i+1:]...)
			return
		}
	}
}

func (r *RoundRobin) Balance(_ string) (string, error) {
	r.Lock()
	defer r.Unlock()
	if len(r.hosts) == 0 {
		return "", NoHostError
	}

	host := r.hosts[r.i%uint64(len(r.hosts))]
	r.i++
	return host, nil
}

func (r *RoundRobin) Inc(_ string) {}

func (r *RoundRobin) Done(_ string) {}
