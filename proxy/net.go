package proxy

import (
	"fmt"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// 这个文件封装网络请求相关

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

// ConnectTimeout 3s超时时间
const ConnectTimeout = 3 * time.Second

// IsBackendAlive ping 代理的目标服务器是否存活
func IsBackendAlive(host string) bool {
	add, err := net.ResolveTCPAddr("tcp", host)
	if err != nil {
		return false
	}

	resolveAddr := fmt.Sprintf("%s:%d", add.IP, add.Port)
	conn, err := net.DialTimeout("tcp", resolveAddr, ConnectTimeout)

	if err != nil {
		return false
	}
	_ = conn.Close()

	return true
}
