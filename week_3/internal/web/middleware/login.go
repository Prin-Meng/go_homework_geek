package middleware

import (
	"encoding/gob"
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

type LoginMiddlewareBuilder struct {
}

// CheckLogin 中间件用于检查用户登录状态。
func (m *LoginMiddlewareBuilder) CheckLogin() gin.HandlerFunc {
	// 注册一下这个类型
	gob.Register(time.Now())
	return func(ctx *gin.Context) {
		path := ctx.Request.URL.Path
		if path == "/users/signup" || path == "/users/login" {
			// 不需要登录校验
			return
		}
		sess := sessions.Default(ctx)
		userId := sess.Get("userId")
		if userId == nil {
			// 中断，不要往后执行，也就是不要执行后面的业务逻辑
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		now := time.Now()
		const updateTimeKey = "update_time"
		// 试着拿出上一次刷新时间
		val := sess.Get(updateTimeKey)
		lastUpdateTime, ok := val.(time.Time)
		if val == nil || !ok || now.Sub(lastUpdateTime) > time.Second*10 {
			// 你这是第一次进来
			sess.Set(updateTimeKey, now)
			sess.Set("userId", userId)
			err := sess.Save()
			if err != nil {
				// 打日志
				fmt.Println(err)
			}
		}
	}
}
