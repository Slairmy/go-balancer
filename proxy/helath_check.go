package proxy

import (
	"log"
	"time"
)

/**
只要在相同的 package下面,所有的公用的数据都是可以共享的。
*/

func (h *HTTPProxy) SetAlive(host string, isAlive bool) {
	// map 并发不安全修改加写锁
	h.Lock()
	defer h.Unlock()
	h.alive[host] = isAlive
}

func (h *HTTPProxy) ReadAlive(host string) bool {
	// map 并发不安全读取加读锁
	h.RLock()
	defer h.RUnlock()
	return h.alive[host]
}

func (h *HTTPProxy) HealthCheck(interval uint) {
	// 检测所有的代理目标服务器
	for host, _ := range h.alive {
		log.Printf("开始心跳检测: %s", host)
		go h.healthCheck(host, interval)
	}
}

// 健康检查做的事情就是定时去 ping 当前服务, ping 当前 host
func (h *HTTPProxy) healthCheck(host string, interval uint) {
	// 所以先搞一个定时器
	ticker := time.NewTicker(time.Duration(interval) * time.Second)

	for range ticker.C {
		// 接收到信号
		if h.ReadAlive(host) && !IsBackendAlive(host) {
			log.Printf("site unreachable remove it: %s", host)
			// 设置为不存活
			h.SetAlive(host, false)
			// 从负载均衡中移除
			h.lb.Remove(host)
		} else if IsBackendAlive(host) && !h.ReadAlive(host) {
			// 可以ping通但是alive是false的需要重新加入
			log.Printf("site reachable add it: %s", host)
			h.SetAlive(host, true)
			h.lb.Add(host)
		}
	}

}
