package service

import (
	"context"
	"errors"
	"go_homework/week_2/internal/domain"
	"go_homework/week_2/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrDuplicateEmail        = repository.ErrDuplicateEmail
	ErrInvalidUserOrPassword = errors.New("the user does not exist or the password is incorrect")
)

type UserService struct {
	repo *repository.UserRepository
}

// NewUserService 函数创建并返回一个 UserService 实例
func NewUserService(repo *repository.UserRepository) *UserService {
	return &UserService{repo: repo}
}

// Signup 函数处理用户的注册流程
func (svc *UserService) Signup(ctx context.Context, u domain.User) error {
	// 使用 bcrypt 算法对用户输入的密码进行哈希
	hash, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		// 如果哈希过程中出现错误，则返回该错误
		return err
	}
	// 将哈希后的密码转换为字符串并赋值给 u.Password
	u.Password = string(hash)
	// 调用仓储层的 Create 方法，将用户数据保存到数据库中
	return svc.repo.Create(ctx, u)
}

// Login 函数用于验证用户的登录信息。
// 它接受上下文、邮箱和密码作为参数。
func (svc *UserService) Login(ctx context.Context, email string, password string) (domain.User, error) {
	// 调用仓库层的 FindByEmail 方法，根据邮箱查找用户。
	u, err := svc.repo.FindByEmail(ctx, email)
	// 如果发生错误，则返回一个空的 domain.User 和错误信息。
	if err != nil {
		return domain.User{}, err
	}
	// 如果没有找到用户，则返回 ErrUserNotFound 错误。
	if errors.Is(err, repository.ErrUserNotFound) {
		return domain.User{}, ErrDuplicateEmail
	}
	// 使用 bcrypt.CompareHashAndPassword 方法，对用户输入的密码进行验证。
	err = bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	// 如果发生错误，则返回一个空的 domain.User 和错误信息。
	if err != nil {
		return domain.User{}, ErrInvalidUserOrPassword
	}
	// 如果登录成功，返回用户信息和 nil，表示没有发生错误。
	return u, nil
}

// UpdateNonSensitiveInfo 在给定的上下文中更新用户的非敏感信息
func (svc *UserService) UpdateNonSensitiveInfo(ctx context.Context,
	user domain.User) error {
	return svc.repo.UpdateNonZeroFields(ctx, user)
}

func (svc *UserService) FindByID(ctx context.Context, uid int64) (domain.User, error) {
	return svc.repo.FindByID(ctx, uid)
}
