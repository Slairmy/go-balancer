package balancer

import (
	"hash/crc32"
	"sync"
)

type IpHash struct {
	mu sync.RWMutex

	hosts []string
}

// 注册工厂
func init() {
	factories["ip_hash"] = NewIpHash
}

func NewIpHash(hosts []string) Balancer {
	return &IpHash{hosts: hosts}
}

func (i *IpHash) Add(host string) {
	i.mu.RLock()
	defer i.mu.RUnlock()
	// 添加到 hosts
	for _, h := range i.hosts {
		if h == host {
			return
		}
	}

	i.hosts = append(i.hosts, host)
}

func (i *IpHash) Remove(host string) {
	i.mu.RLock()
	defer i.mu.RUnlock()
	// 将host从目标服务器中移除
	for idx, h := range i.hosts {
		if h == host {
			i.hosts = append(i.hosts[idx:], i.hosts[idx+1:]...)
		}
	}
}

func (i *IpHash) Balance(clientIP string) (string, error) {
	// ip hash算法 crc32() 拿到ip(这里是拿客户端的ip) 的hash值然后 % len()
	i.mu.Lock()
	defer i.mu.Unlock()
	if len(i.hosts) == 0 {
		return "", NoHostError
	}

	index := crc32.ChecksumIEEE([]byte(clientIP)) % uint32(len(i.hosts))

	return i.hosts[index], nil
}

func (i *IpHash) Inc(_ string) {}

func (i *IpHash) Done(_ string) {}
