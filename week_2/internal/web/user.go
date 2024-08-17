package web

import (
	"errors"
	regexp "github.com/dlclark/regexp2"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"go_homework/week_2/internal/domain"
	"go_homework/week_2/internal/service"
	"net/http"
	"time"
)

const (
	// 电子邮箱正则表达式模式，用于校验电子邮箱格式
	emailRegexPattern = "^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\\.[a-zA-Z]{2,}$"
	// 密码正则表达式模式，用于校验密码强度
	passwordRegexPattern = "^(?=.*[a-z])(?=.*[A-Z])(?=.*\\d)(?=.*[!@#$%^&*()-_=+\\[{\\]};:'\",<.>/?])[a-zA-Z0-9!@#$%^&*()-_=+\\[{\\]};:'\",<.>/?]{8,16}$"
)

type UserHandler struct {
	emailRexExp    *regexp.Regexp
	passwordRexExp *regexp.Regexp
	svc            *service.UserService
}

// NewUserHandler 函数创建并返回一个 UserHandler 类型的指针
func NewUserHandler(svc *service.UserService) *UserHandler {
	return &UserHandler{
		// 编译正则表达式来验证邮箱
		emailRexExp: regexp.MustCompile(emailRegexPattern, regexp.None),
		// 编译正则表达式来验证密码
		passwordRexExp: regexp.MustCompile(passwordRegexPattern, regexp.None),
		// 存储 UserService 类型的指针，用于后续用户操作
		svc: svc,
	}
}

// RegisterRoutes 函数为用户处理器注册路由
func (h *UserHandler) RegisterRoutes(server *gin.Engine) {
	ug := server.Group("/users")
	ug.POST("/signup", h.Signup)
	ug.POST("/login", h.Login)
	ug.POST("/edit", h.Edit)
	ug.GET("/profile", h.Profile)
}

func (h *UserHandler) Signup(ctx *gin.Context) {
	type SignupRequest struct {
		Email           string `json:"email"`
		Password        string `json:"password"`
		ConfirmPassword string `json:"confirmPassword"`
	}

	// 定义 SignupRequest 结构体变量 req，用于接收请求参数
	var req SignupRequest
	// 调用 ctx.Bind() 来绑定请求数据到 req 变量
	if err := ctx.Bind(&req); err != nil {
		// 如果绑定过程发生错误，直接返回，不进行后续处理
		return
	}

	// 使用编译好的正则表达式检查邮箱格式是否正确，isEmail 为布尔值
	isEmail, err := h.emailRexExp.MatchString(req.Email)
	// 如果正则表达式匹配过程中发生错误
	if err != nil {
		// 返回 200 OK 状态码，并给出错误信息
		ctx.String(http.StatusOK, "System error")
		// 因为出错，所以结束当前请求的处理流程，return 执行完该方法会终止后续流程
		return
	}
	// 如果邮箱格式不正确
	if !isEmail {
		// 返回 200 OK 状态码，并给出错误信息
		ctx.String(http.StatusOK, "Email is invalid")
		// 因为邮箱格式出错，所以结束当前请求的处理流程
		return
	}

	// 检查密码字段是否匹配确认密码字段
	if req.Password != req.ConfirmPassword {
		// 返回 200 OK 状态码，并给出错误信息
		ctx.String(http.StatusOK, "Password is not same")
		// 因为密码不一致，所以结束当前请求的处理流程
		return
	}

	// 检查密码是否符合密码强度要求
	isPassword, err := h.passwordRexExp.MatchString(req.Password)
	// 如果正则表达式匹配过程中发生错误
	if err != nil {
		// 返回 200 OK 状态码，并给出错误信息
		ctx.String(http.StatusOK, "System error")
		// 因为出错，所以结束当前请求的处理流程
		return
	}
	// 如果密码不符合要求
	if !isPassword {
		// 返回 200 OK 状态码，并给出错误信息，提示密码强度要求
		ctx.String(http.StatusOK, "The password must contain at least one uppercase letter, one lowercase letter, "+
			"one number, and one special character, and must be between 8 and 16 characters long")
		// 因为密码格式出错，所以结束当前请求的处理流程
		return
	}

	// 调用 h.svc 指针，为指定的上下文和用户对象进行注册
	err = h.svc.Signup(ctx, domain.User{
		// 设置用户的电子邮箱
		Email: req.Email,
		// 设置用户的密码
		Password: req.Password,
	})

	// 根据注册返回的错误情况（err），执行不同的操作
	switch {
	// 成功注册
	case err == nil:
		// 在上下文中设置状态代码为 200 OK
		ctx.String(http.StatusOK, "Signup success")
		// 服务层定义的错误，表示该邮箱已存在
	case errors.Is(err, service.ErrDuplicateEmail):
		// 在上下文中设置状态代码为 200 OK，并给出错误消息提示
		ctx.String(http.StatusOK, "Email is already exist")
		// 其他未知的错误
	default:
		// 在上下文中设置状态代码为 200 OK，提示系统错误
		ctx.String(http.StatusOK, "System error")
		// 根据错误类型（err），展示不同的提示
	}

}

// Login 函数用于验证用户的登录信息
func (h *UserHandler) Login(ctx *gin.Context) {
	// 定义登录请求结构体
	type LoginRequest struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	var req LoginRequest
	// 解析请求中的 JSON 数据到登录请求结构体
	if err := ctx.Bind(&req); err != nil {
		return
	}
	// 调用服务层的 Login 方法进行登录验证
	u, err := h.svc.Login(ctx, req.Email, req.Password)
	// 根据 Login 方法的返回结果进行不同的响应处理
	switch {
	case err == nil:
		// 登录成功，设置会话信息
		sess := sessions.Default(ctx)
		sess.Set("userId", u.Id)
		sess.Options(sessions.Options{
			MaxAge: 900,
		})
		// 保存会话信息
		err = sess.Save()
		if err != nil {
			// 保存会话信息出错，返回系统错误
			ctx.String(http.StatusOK, "System error")
			return
		}
		// 登录成功，返回成功消息
		ctx.String(http.StatusOK, "Login success")
	case errors.Is(err, service.ErrInvalidUserOrPassword):
		// 登录失败，密码错误，返回相应的错误消息
		ctx.String(http.StatusOK, "Invalid user or password")
	default:
		// 其他错误，返回系统错误消息
		ctx.String(http.StatusOK, "System error")
	}
}

func (h *UserHandler) Edit(ctx *gin.Context) {
	type EditRequest struct {
		Nickname string `json:"nickname"`
		// YYYY-MM-DD 格式的生日日期字符串
		Birthday string `json:"birthday"`
		About    string `json:"about"`
	}
	// 声明一个变量用来接收解析 Web 请求后的表单数据
	var req EditRequest
	// 调用 ctx.Bind 函数对请求参数进行自动解析
	if err := ctx.Bind(&req); err != nil {
		// 如果解析过程中发生错误，则打印错误并返回空
		return
	}
	// 通过 sessions.Default(ctx) 获取当前上下文的会话对象
	sess := sessions.Default(ctx)
	// 尝试从会话中获取用户 ID，如果获取不到，则说明用户未登录
	userId := sess.Get("userId")
	// 判断用户 ID 是否为空，为空则表示用户未登录或登录已过期
	if userId == nil {
		// 返回 401 状态码表示用户未认证
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	// 尝试将请求中的生日字符串转换为 time.Time 类型
	// 使用 time.DateOnly 作为格式模板，确保只解析日期部分
	birthday, err := time.Parse(time.DateOnly, req.Birthday)
	// 发生错误时（如格式不正确），打印错误信息并返回
	if err != nil {
		ctx.String(http.StatusOK, "wrong birthday format")
		return
	}
	// 调用服务层的 UpdateNonSensitiveInfo 方法来更新用户的非敏感信息
	// 通过 domain.User{...} 构造一个用户对象，包含从请求中解析的 ID、昵称、生日和个人简介
	err = h.svc.UpdateNonSensitiveInfo(ctx, domain.User{
		Id:       userId.(int64),
		Nickname: req.Nickname,
		Birthday: birthday,
		AboutMe:  req.About,
	})
	// 如果更新过程中发生错误（如数据库更新失败），打印错误信息并返回
	if err != nil {
		ctx.String(http.StatusOK, "System error")
		return
	}
	// 打印编辑成功的信息，并返回 200 状态码表示请求成功
	ctx.String(http.StatusOK, "Edit success")

}

func (h *UserHandler) Profile(ctx *gin.Context) {
	sess := sessions.Default(ctx)
	userId := sess.Get("userId")
	uid, ok := userId.(int64)
	if !ok {
		ctx.String(http.StatusOK, "System error")
		return
	}

	u, err := h.svc.FindByID(ctx, uid)
	if err != nil {
		ctx.String(http.StatusOK, "System error")
		return
	}
	type User struct {
		Email    string
		Nickname string
		Birthday string
		AboutMe  string
	}

	repData := User{
		Email:    u.Email,
		Nickname: u.Nickname,
		Birthday: u.Birthday.Format(time.DateOnly),
		AboutMe:  u.AboutMe,
	}
	ctx.JSON(http.StatusOK, repData)
}
