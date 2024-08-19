package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"go_homework/week_3/internal/web"
	"log"
	"net/http"
	"strings"
	"time"
)

// LoginJWTMiddlewareBuilder 定义了一个名为 LoginJWTMiddlewareBuilder 的结构体，用于构建处理 JWT 登录校验的中间件
type LoginJWTMiddlewareBuilder struct {
}

// CheckLogin 定义了一个名为 CheckLogin 的方法，返回一个 gin.HandlerFunc 类型的函数
func (m *LoginJWTMiddlewareBuilder) CheckLogin() gin.HandlerFunc {
	// 返回的函数，用于处理登录校验逻辑
	return func(ctx *gin.Context) {
		// 获取当前请求的路径
		path := ctx.Request.URL.Path
		// 判断当前请求路径是否为注册或登录的接口
		if path == "/users/signup" || path == "/users/login" || path == "/hello" {
			// 如果是注册或登录的接口，不需要进行登录校验，直接返回
			return
		}
		// 从请求头中获取 Authorization 字段的值 like Bearer XXXX
		authCode := ctx.GetHeader("Authorization")
		// 如果请求头中没有 Authorization 字段，或者字段值为空
		if authCode == "" {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		// 将令牌字符串进行分割，获取令牌部分
		segs := strings.Split(authCode, " ")
		// 没登录，Authorization 中的内容是乱传的
		if len(segs) != 2 {
			// 返回 401 状态码（未授权）
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		// 获取分割后的令牌字符串
		tokenStr := segs[1]
		// 声明一个 UserClaims 结构体指针变量用于解析 JWT 中的声明信息
		var uc web.UserClaims
		token, err := jwt.ParseWithClaims(tokenStr, &uc, func(token *jwt.Token) (interface{}, error) {
			// 返回用于验证 JWT 签名的密钥
			return web.JWTKey, nil
		})

		// token 是伪造的
		if err != nil {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		// 如果令牌无效或已过期
		if token == nil || !token.Valid {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		// 如果 JWT 中的用户代理信息与请求头中的 User-Agent 信息不匹配
		if uc.UserAgent != ctx.GetHeader("User-Agent") {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		// 获取 JWT 中的过期时间
		expireTime := uc.ExpiresAt

		// 判断剩余过期时间是否小于 10 分钟
		if expireTime.Sub(time.Now()) < time.Minute*10 {
			// 更新 JWT 中的过期时间为当前时间加 1 小时
			uc.ExpiresAt = jwt.NewNumericDate(time.Now().Add(time.Hour))
			// 使用更新后的声明信息生成新的令牌字符串
			tokenStr, err = token.SignedString(web.JWTKey)
			// 在响应头中添加 x-jwt-token 字段，用于提供新的令牌
			ctx.Header("x-jwt-token", tokenStr)
			// 如果生成新令牌的过程中发生错误
			if err != nil {
				// 不要终止请求处理流程，因为即使令牌的过期时间没有被刷新，用户仍然是登录状态
				log.Println(err)
			}
		}
		// 将解析后的声明信息存储到上下文中，以便后续处理使用
		ctx.Set("user", uc)
	}

}
