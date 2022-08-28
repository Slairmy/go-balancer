package balancer

import (
	fibHeap "github.com/starwander/GoFibonacciHeap"
	"sync"
)

// least load 最小负载算法,将请求定向到负载量最小的目标主机中
// 对于最小负载算法而言，如果把所有主机的负载值动态存入动态数组中，寻找负载最小节点的时间复杂度为O(N)
// 为什么不是把负载动态值存入 map? 因为是对比负载值,即使是存入map也是需要遍历所有的key,寻找的时间复杂度也是O(N)
// 如果把主机的负载值维护成一个红黑树，那么寻找负载最小节点的时间复杂度为O(logN)，我们这里利用的数据结构叫做 斐波那契堆 ，寻找负载最小节点的时间复杂度为O(1)

func init() {
	// factories["least_load"] = NewLeastLoad
}

// 在p2c模式定义了host结构体相同的package不需要重复定义了
//type host struct {
//	name string
//	load uint64
//}

// 这里host 仅为了实现 fibHeap 的Value接口

func (h *host) Tag() interface{} {
	return h.name
}

// Key 为什么这里需要返回浮点型?
func (h *host) Key() float64 {
	return float64(h.load)
}

type LeastLoad struct {
	sync.RWMutex
	heap *fibHeap.FibHeap // 这里应该是斐波那契堆的类型
}

func NewLeastLoad(hosts []string) Balancer {
	l := &LeastLoad{heap: fibHeap.NewFibHeap()}

	for _, h := range hosts {
		l.Add(h)
	}

	return l
}

func (l *LeastLoad) Add(hostName string) {
	// 这里主要是如何使用斐波那契堆初始化负载
	l.Lock()
	defer l.Unlock()

	if ok := l.heap.GetValue(hostName); ok != nil {
		return
	}

	_ = l.heap.InsertValue(&host{name: hostName, load: 0})
}

func (l *LeastLoad) Remove(hostName string) {
	l.Lock()
	defer l.Unlock()

	if ok := l.heap.GetValue(hostName); ok != nil {
		_ = l.heap.Delete(hostName)
	}
}

func (l *LeastLoad) Balance(_ string) (string, error) {
	l.Lock()
	defer l.Unlock()

	// 堆没有数据了
	if l.heap.Num() == 0 {
		return "", NoHostError
	}

	// 这里拿到堆的min数据
	return l.heap.ExtractMinValue().Tag().(string), nil
}

func (l *LeastLoad) Inc(hostName string) {
	l.Lock()
	defer l.Unlock()

	if ok := l.heap.GetValue(hostName); ok == nil {
		return
	}

	h := l.heap.GetValue(hostName)
	// l.heap.GetValue 返回的是 &host{} ,为什么不能直接使用 *h 的形式?
	// 自我解答: 因为GetValue返回的是一个 Value类型的数据,他有2个方法,但是我们这里需要断言出他是 host(因为 host已经实现了 Value) 才能拿到 load 负载量
	h.(*host).load++
	// 感觉这里的命名有点奇怪,反正就是IncreaseKeyValue要先把Increase的量修改好 -- 更对的是堆自身需要维护一定的规则的结构
	_ = l.heap.IncreaseKeyValue(h)
}

func (l *LeastLoad) Done(hostName string) {
	l.Lock()
	defer l.Unlock()

	if ok := l.heap.GetValue(hostName); ok == nil {
		return
	}

	h := l.heap.GetValue(hostName)
	h.(*host).load--
	_ = l.heap.DecreaseKeyValue(h)
}
