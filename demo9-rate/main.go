package main

import (
	"fmt"
	"net/http"
	"time"

	"golang.org/x/time/rate"
)

func main() {
	// ServeMux类型是HTTP请求的多路转接器。它会将每一个接收的请求的URL与一个注册模式的列表进行匹配，并调用和URL最匹配的模式的处理器。
	mux := http.NewServeMux()
	mux.HandleFunc("/", defaultHttp)
	http.ListenAndServe(":8080", middlewareLimit(mux))
}

// 限流桶，每2s一个请求
var limiter = rate.NewLimiter(rate.Every(time.Second), 1)

func middlewareLimit(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if limiter.Allow() == false {
			fmt.Println("---limit")
			delay := limiter.Reserve().Delay()
			time.Sleep(delay)
			fmt.Println("---over", delay)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// 默认http处理
func defaultHttp(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path

	if path == "/" {
		w.Write([]byte("index"))
		fmt.Println("index")
		return
	}

	// 自定义404
	http.Error(w, "you lost???", http.StatusNotFound)
}
