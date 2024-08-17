package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"go_homework/week_2/internal/repository"
	"go_homework/week_2/internal/repository/dao"
	"go_homework/week_2/internal/service"
	"go_homework/week_2/internal/web"
	"go_homework/week_2/internal/web/middleware"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"strings"
	"time"
)

// main 函数是应用的启动点，负责初始化数据库、配置服务器和启动服务
func main() {
	// 初始化数据库连接
	db := initDB()
	// 初始化 Web 服务器
	server := initWebServer()
	// 初始化用户处理器，主要负责实现用户相关的路由和逻辑
	initUserHdl(db, server)
	// 启动 Web 服务器，并监听 8080 端口，返回错误信息，err!= nil 时退出
	err := server.Run(":8080")
	if err != nil {
		return
	}
}

func initUserHdl(db *gorm.DB, server *gin.Engine) {
	// 将 GORM 数据库实例传入 UserDAO
	ud := dao.NewUserDAO(db)
	// 在 UserRepository 中，通过之前创建的 UserDAO 初始化
	ur := repository.NewUserRepository(ud)
	// 实例化 UserService，并注入 UserRepository
	us := service.NewUserService(ur)
	// 创建 UserHandler 实例以便处理用户相关的请求，其中包含用户服务对象
	hdl := web.NewUserHandler(us)
	// 调用 UserHandler 的 RegisterRoutes 方法，向引擎注册用户相关的路由
	hdl.RegisterRoutes(server)
}

func initDB() *gorm.DB {
	// 打开一个数据库连接，使用了 root 用户，密码为 root，连接到 13316 端口上运行的 MySQL 服务，数据库名为 webook
	db, err := gorm.Open(mysql.Open("root:root@tcp(localhost:13316)/webook"))
	// 如果打开数据库连接时发生错误，使用 panic 抛出一个错误，以确保程序能够停止执行，并提醒开发者处理这个错误
	if err != nil {
		panic(err)
	}
	// 初始化数据库表，使用 dao 包中的 InitTables 函数
	err = dao.InitTables(db)
	// 如果初始化表结构时发生错误，也使用 panic 抛出错误
	if err != nil {
		panic(err)
	}
	// 返回初始化后的数据库连接对象
	return db
}

// initWebServer 函数用于初始化 Web 服务器，设置路由和中间件
func initWebServer() *gin.Engine {
	// 创建 gin.Default() 对象，该对象默认包含了 Logger 和 Recovery 中间件
	server := gin.Default()
	// 使用 CORS 中间件，允许跨域请求，并指定了允许的请求头，以及允许的来源
	server.Use(cors.New(cors.Config{
		// 是否允许在跨域请求中携带用户凭证（如 cookies、HTTP 认证）
		AllowCredentials: true,
		// 允许的请求头列表
		AllowHeaders: []string{"Content-Type"},
		// 验证来源的函数，根据函数的返回值决定是否允许该来源
		AllowOriginFunc: func(origin string) bool {
			// 如果来源是以 http://localhost 开头的，就允许该来源
			if strings.HasPrefix(origin, "http://localhost") {
				//if strings.Contains(origin, "localhost") {
				return true
			}
			// 如果来源包含 bt.com，就允许该来源
			return strings.Contains(origin, "bt.com")
		},
		// 将预检请求的结果缓存 12 小时，减少预检请求的次数，提高效率
		MaxAge: 12 * time.Hour,
	}))

	login := &middleware.LoginMiddlewareBuilder{}
	// 存储数据的，也就是你 userId 存哪里
	// 直接存 cookie
	store := cookie.NewStore([]byte("secret"))
	server.Use(sessions.Sessions("ssid", store), login.CheckLogin())
	return server
}
