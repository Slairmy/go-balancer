package balancer

import "github.com/lafikl/consistent"

// 一致性hash负载均衡法
// 首先通过将所有代理主机节点通过 CRC32散列算法 对主机IP或者UUID计算hash值,将得到的所有值连成环(哈希环)
// 当请求到来的时候,通过CRC32对请求IP计算hash值,然后沿着环顺时针查找,找到第一个主机节点(可以通过2分法寻找)
// 可能遇到主机节点比较少,节点分布不均匀,这种情况可以引入虚拟节点比如节点 127.0.0.1 可以引入虚拟节点 127.0.0.1#1, 127.0.0.1#2 等

// 一致性hash有实现方式 github.com/lafikl/consistent

type ConsistentHash struct {
	ch *consistent.Consistent
	// 这里应该有一个结构存储代理主机节点(环) 而且这里不需要加锁 consistent.Consistent 有锁处理
}

func init() {

	factories["consistent_hash"] = NewConsistentHash

}

func NewConsistentHash(hosts []string) Balancer {
	c := &ConsistentHash{ch: consistent.New()}

	for _, h := range hosts {
		c.ch.Add(h)
	}

	return c
}

func (c *ConsistentHash) Add(hostName string) {
	c.ch.Add(hostName)
}

func (c *ConsistentHash) Remove(hostName string) {
	c.ch.Remove(hostName)
}

func (c *ConsistentHash) Balance(clientIP string) (string, error) {
	if len(c.ch.Hosts()) == 0 {
		return "", NoHostError
	}
	return c.ch.Get(clientIP)
}

// 传统一致性hash算法不需要计算负载量 -- 有界一致性hash需要计算负载量 consistent 使用的负载因子ß=0.25

func (c *ConsistentHash) Inc(_ string) {
}

func (c *ConsistentHash) Done(_ string) {
}
