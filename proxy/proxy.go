package proxy

import (
	"fmt"
	"github.com/slairmy/balancer/balancer"
	"log"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"sync"
)

// X-Real-IP=客户端的真实IP 修改X-Proxy为代理服务器IP
var (
	XRealIP       = http.CanonicalHeaderKey("X-Real-IP")
	XProxy        = http.CanonicalHeaderKey("X-Proxy")
	XForwardedFor = http.CanonicalHeaderKey("X-Forwarded-For")
)

var (
	ReverseProxy = "Balancer-Reverse-Proxy"
)

// hostMap主机对反向代理的映射
// alive 反向代理主机是否存活

type HTTPProxy struct {
	hostMap map[string]*httputil.ReverseProxy
	lb      balancer.Balancer

	sync.RWMutex
	alive map[string]bool
}

// NewHTTPProxy 传入代理主机切片数据 和 负载均衡算法
func NewHTTPProxy(targetHosts []string, algorithm string) (*HTTPProxy, error) {
	// 第一步: 初始化 hostMap
	hosts := make([]string, 0)
	hostMap := make(map[string]*httputil.ReverseProxy)
	alive := make(map[string]bool)

	for _, targetHost := range targetHosts {
		targetUrl, err := url.Parse(targetHost)
		if err != nil {
			return nil, err
		}

		proxy := httputil.NewSingleHostReverseProxy(targetUrl)
		// 修改请求指向到真实的服务器地址
		originDirector := proxy.Director
		// 重写header
		proxy.Director = func(request *http.Request) {
			originDirector(request)
			request.Header.Set(XProxy, ReverseProxy)
			request.Header.Set(XRealIP, GetIP(request))
		}

		host := GetHost(targetUrl)
		alive[host] = true // 默认主机存活
		hostMap[host] = proxy
		hosts = append(hosts, host)
	}

	// todo 负载均衡策略
	lb, err := balancer.Build(algorithm, hosts)
	if err != nil {
		return nil, err
	}

	return &HTTPProxy{hostMap: hostMap, alive: alive, lb: lb}, nil
}

func (h *HTTPProxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// 主机都加载到了hostMap 如何实现ServerHTTP ? 这里实现 ServerHTTP 是直接接收请求来处理了
	// 执行流程是 ServerHTTP 监听客户端请求 然后通过反向代理选择一个负载均衡算法选择目标服务器

	defer func() {
		// todo 这里没太理解到
		if err := recover(); err != nil {
			log.Printf("proxy causes panic: %s", err)
			w.WriteHeader(http.StatusBadGateway)
			_, _ = w.Write([]byte(err.(error).Error()))
		}
	}()

	host, err := h.lb.Balance(GetIP(r))
	if err != nil {
		// 没有返回值又中断执行怎么办 ? 返回异常的response
		w.WriteHeader(http.StatusBadGateway)
		_, _ = w.Write([]byte(fmt.Sprintf("balanve error: %s", err)))
		return
	}

	// 负载均衡服务器 + 1
	h.lb.Add(host)
	defer h.lb.Done(host)

	h.hostMap[host].ServeHTTP(w, r)
}

// 这里是为了拿到客户端的IP,举个例子
// XForwardedFor 记录的是从客户端地址到最后一个代理服务器的地址 假设客户端IP 是 192.168.1.0 经过2个代理服务器 192.168.1.1 和 192.168.1.2
// 当客户端请求到第一个代理服务器的时候 X-Forwarded-For=192.168.1.0, 192.168.1.1 请求到第二个代理服务器的时候 X-Forwarded-For=192.168.1.0, 192.168.1.1, 192.168.1.2

func GetIP(r *http.Request) string {

	clientIP, _, _ := net.SplitHostPort(r.RemoteAddr)

	if len(r.Header.Get(XForwardedFor)) != 0 { // 如果没有XForwardedFor直接拿XRealIP ? 这里为什么不直接拿 XRealIP ?
		// 这里没太看懂
		xff := r.Header.Get(XForwardedFor) // 这里是字符串
		s := strings.Index(xff, ", ")      // ", "在xff中首次出现的位置
		if s == -1 {
			s = len(r.Header.Get(XForwardedFor))
		}
		clientIP = xff[:s]
	} else if len(r.Header.Get(XRealIP)) != 0 {
		clientIP = r.Header.Get(XRealIP)
	}

	return clientIP
}

// GetHost 获取 host IP:Port的形式
func GetHost(url *url.URL) string {
	// 无法解析请求链接的情况
	if _, _, err := net.SplitHostPort(url.Host); err == nil {
		return url.Host
	}

	if url.Scheme == "http" {
		return fmt.Sprintf("%s:%s", url.Host, "80")
	} else if url.Scheme == "https" {
		return fmt.Sprintf("%s:%s", url.Host, "443")
	}

	return url.Host
}
