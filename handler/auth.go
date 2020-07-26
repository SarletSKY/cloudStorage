package handler

import (
	"net/http"
)

// token 拦截器
func HTTPInterceptor(h http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(
		func(writer http.ResponseWriter, request *http.Request) {
			request.ParseForm()

			username := request.Form.Get("username")
			token := request.Form.Get("token")

			// 验证token
			if len(username) < 3 || !ValidToToken(token) {
				writer.WriteHeader(http.StatusForbidden)
				return
			}

			h(writer, request)
		})
}

// 跨域请求
func ReceiveClientRequest(h http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(
		func(writer http.ResponseWriter, request *http.Request) {
			writer.Header().Set("Access-Control-Allow-Origin", "*") //允许访问所有域

			request.ParseForm()
			h(writer, request)
		})
}
