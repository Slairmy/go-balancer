package middleware

import (
	"github.com/gorilla/mux"
	"net/http"
)

// 一般来说 mux middleware 是使用方式是 r.Use()
// func (r *Router) Use(mwf ...MiddlewareFunc) {
//	for _, fn := range mwf {
//		r.middlewares = append(r.middlewares, fn)
//	}
//}
// 所以Use参数传一个 MiddlewareFunc 类型就好了
// MiddlewareFunc 的类型定义是这样的 type MiddlewareFunc func(http.Handler) http.Handler
// 所以当返回中间件 返回 MiddlewareFunc 类型 也是就是 返回一个 func(http.Handler) http.Handler

func MaxAllowedMiddleware(n uint) mux.MiddlewareFunc {

	sem := make(chan struct{}, n)
	acquire := func() { sem <- struct{}{} } // 写入管道
	release := func() { <-sem }             // 写出

	// 被套娃套晕了
	return func(next http.Handler) http.Handler { // 这里看上面的注释
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// 到这里计数
			acquire()
			defer release() // handle 执行完了就释放计数
			next.ServeHTTP(w, r)
		})
	}
}
