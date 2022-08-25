package main

import (
	"Balancer/middleware"
	"Balancer/proxy"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strconv"
)

func main() {
	// 先读取配置
	config, err := ReadConfig("config.yaml")
	if err != nil {
		log.Fatalf("read config error: %s", err)
	}

	// 参数校验
	err = config.Validation()
	if err != nil {
		log.Fatalf("config validation error: %s", err)
	}

	// 开启一个什么路由

	router := mux.NewRouter()

	for _, l := range config.Location {
		httpProxy, err := proxy.NewHTTPProxy(l.ProxyPass, l.BalanceMode)
		if err != nil {
			log.Fatalf("create proxy error: %s", err)
		}

		// 如果配置了健康检查
		if config.HealthCheck {
			httpProxy.HealthCheck(config.HealthCheckInterval)
		}

		router.Handle(l.Pattern, httpProxy)
	}

	if config.MaxAllowed > 0 {
		// 允许最大的并发请求数 -- 设置中间件
		router.Use(middleware.MaxAllowedMiddleware(config.MaxAllowed))
	}

	srv := http.Server{
		Addr:    ":" + strconv.Itoa(config.Port),
		Handler: router,
	}

	config.Print()

	if config.Schema == "http" {
		err := srv.ListenAndServe()
		if err != nil {
			log.Printf("http server listen error: %s", err)
		}
	} else if config.Schema == "https" {
		err := srv.ListenAndServeTLS(config.SSLCertificate, config.SSLCertificateKey)
		if err != nil {
			log.Printf("https server listen error: %s", err)
		}
	}
}
