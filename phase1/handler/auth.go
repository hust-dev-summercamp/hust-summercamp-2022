package handler

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

// HTTPInterceptor : http请求拦截器
//func HTTPInterceptor(h http.HandlerFunc) http.HandlerFunc {
//	return http.HandlerFunc(
//		func(w http.ResponseWriter, r *http.Request) {
//			r.ParseForm()
//			username := r.Form.Get("username")
//			token := r.Form.Get("token")
//
//			//验证登录token是否有效
//			if len(username) < 3 || !IsTokenValid(token, username) {
//				// token校验失败则直接返回失败提示
//				resp := util.NewRespMsg(
//					int(common.StatusInvalidToken),
//					"token无效",
//					nil,
//				)
//				w.Write(resp.JSONBytes())
//				return
//			}
//			h(w, r)
//		})
//}

// Authorize : http请求拦截器
func Authorize() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Request.ParseForm()
		username := c.Request.FormValue("username") // 用户名
		token := c.Request.FormValue("token")       // 访问令牌
		fmt.Println(username)
		fmt.Println(len(token))

		if len(username) < 3 || !IsTokenValid(token, username) {
			// 验证不通过，不再调用后续的函数处理
			c.Abort()
			c.JSON(http.StatusUnauthorized, gin.H{"message": "访问未授权"})
			// return可省略, 只要前面执行Abort()就可以让后面的handler函数不再执行
			return
		}
		c.Next()
	}
}
