package balancer

import (
	"hash/crc32"
	"math/rand"
	"sync"
	"time"
)

// p2c = power of 2 random choice
// 原理如下
// 1、如果请求IP为空,p2c将随机选择两个代理主机节点,最后选择其中负载量较小的节点
// 2、如果请求IP不为空,p2c通过对IP地址以及对IP地址加盐进行CRC32哈希计算,则会得到两个32bit的值,将其对主机数量进行取模 也就是 CRC32(IP) % len(hosts)
// 	  CRC32(IP + salt) len(hosts) 最后选择负载数量较小的节点

// 这里需要额外存在盐值salt 以及 额外记录负载量

type host struct {
	name string // 主机名
	load uint64 // 负载量
}

const Salt = "%#!"

type P2c struct {
	sync.RWMutex
	hosts []*host

	rnd     *rand.Rand
	loadMap map[string]*host
}

func init() {
	factories["p2c"] = NewP2c
}

func NewP2c(hosts []string) Balancer {
	p := &P2c{
		hosts:   []*host{},
		loadMap: map[string]*host{},
		rnd:     rand.New(rand.NewSource(time.Now().UnixNano())),
	}
	// 初始化
	for _, h := range hosts {
		p.Add(h)
	}

	return p
}

func (p *P2c) Add(hostName string) {
	p.Lock()
	defer p.Unlock()

	// 有记录负载就不用这么判断了
	if _, ok := p.loadMap[hostName]; ok {
		return
	}
	h := &host{name: hostName, load: 0}

	p.hosts = append(p.hosts, h)
	p.loadMap[hostName] = h
}

func (p *P2c) Remove(hostName string) {
	p.Lock()
	defer p.Unlock()

	// 这里还是可以先判断下 有记录就不用走后续流程了
	if _, ok := p.loadMap[hostName]; ok {
		return
	}

	for i, h := range p.hosts {
		if h.name == hostName {
			p.hosts = append(p.hosts[i:], p.hosts[i+1:]...)
			// 去除负载
			delete(p.loadMap, hostName)
			return
		}
	}
}

func (p *P2c) Balance(clientHost string) (string, error) {
	// 随机选取一个
	p.Lock()
	defer p.Unlock()

	if len(p.hosts) == 0 {
		return "", NoHostError
	}
	n1, n2 := p.hash(clientHost)
	// 判断哪个负载小返回哪个
	if p.loadMap[n1].load < p.loadMap[n2].load {
		return n1, nil
	}

	return n2, nil
}

// 这里写一个函数主要是处理p2c的规则,拿到随机的两个主机host
// 这里不要陷入一个误区导致不知道怎么写,就是随机选择2个即使随机到的两个主机一样那也照样返回一样的
func (p *P2c) hash(ip string) (string, string) {
	var n1, n2 string

	if len(ip) > 0 {
		slatIP := ip + Salt
		n1 := p.hosts[crc32.ChecksumIEEE([]byte(ip))%uint32(len(p.hosts))].name
		n2 := p.hosts[crc32.ChecksumIEEE([]byte(slatIP))%uint32(len(p.hosts))].name

		return n1, n2
	}

	n1 = p.hosts[p.rnd.Intn(len(p.hosts))].name
	n2 = p.hosts[p.rnd.Intn(len(p.hosts))].name

	return n1, n2
}

func (p *P2c) Inc(hostName string) {
	// 这里是修改公共资源写锁
	p.Lock()
	defer p.Unlock()
	h, ok := p.loadMap[hostName]
	if !ok {
		return
	}

	// 还是要判断
	h.load++
}

func (p *P2c) Done(hostName string) {
	// 这里是修改公共资源写锁
	p.Lock()
	defer p.Unlock()
	h, ok := p.loadMap[hostName]
	if !ok {
		return
	}

	// 还是要判断
	h.load--
}
