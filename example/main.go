package main

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
)

// HTTPProxy 反向代理
type HTTPProxy struct {
	proxy *httputil.ReverseProxy
}

// NewHTTPProxy 构造函数
func NewHTTPProxy(target string) (*HTTPProxy, error) {
	u, err := url.Parse(target)
	if err != nil {
		return nil, err
	}

	return &HTTPProxy{proxy: httputil.NewSingleHostReverseProxy(u)}, nil
}

// 必须实现 ServeHTTP 方法
func (h *HTTPProxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.proxy.ServeHTTP(w, r)
}

func main() {
	//// 反向代理 ip
	//proxy, err := NewHTTPProxy("http://127.0.0.1:9100")
	//if err != nil {
	//	fmt.Println(err)
	//}
	//http.Handle("/", proxy)
	//http.ListenAndServe(":9101", nil)

	//xff := "192.168.1.1, 192.168.1.2, 1"
	//s := strings.Index(xff, ", ")
	//
	//fmt.Println(xff)
	//fmt.Println(s)
	//fmt.Println(xff[:s])

	targetUrl, _ := url.Parse("https://zhuanlan.zhihu.com/p/506415782")
	fmt.Println(targetUrl.Host)

}
