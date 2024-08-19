package repository

import (
	"context"
	"go_homework/week_3/internal/domain"
	"go_homework/week_3/internal/repository/dao"
	"time"
)

var (
	ErrDuplicateEmail = dao.ErrDuplicateEmail
	ErrUserNotFound   = dao.ErrRecordNotFound
)

type UserRepository struct {
	dao *dao.UserDAO
}

func NewUserRepository(dao *dao.UserDAO) *UserRepository {
	// NewUserRepository 函数用于创建并返回一个 UserRepository 类型的指针
	// 参数 dao 是一个指向 UserDAO 类型的指针，用于初始化 UserRepository 的 dao 字段
	return &UserRepository{dao: dao}
}

// Create 函数用于创建新用户。
func (repo *UserRepository) Create(ctx context.Context, u domain.User) error {
	return repo.dao.Insert(ctx, dao.User{
		Email:    u.Email,
		Password: u.Password,
	})
}

// FindByEmail 方法根据邮箱地址查询用户信息
func (repo *UserRepository) FindByEmail(ctx context.Context, email string) (domain.User, error) {
	// 调用 dao 层的 FindByEmail 方法，传入上下文和邮箱地址作为参数
	u, err := repo.dao.FindByEmail(ctx, email)
	// 如果发生错误，则返回一个空的 domain.User 和 error
	if err != nil {
		return domain.User{}, err
	}
	// 返回转换后的 domain.User 和 nil，表示没有发生错误
	return repo.toDomain(u), nil
}

// toDomain 方法用于将数据访问对象 (DAO) 中的 User 结构转换为域模型中的 User 结构。
func (repo *UserRepository) toDomain(u dao.User) domain.User {
	// 返回转换后的 domain.User 类型对象，包含原始 dao.User 对象的 ID、电子邮件和密码字段。
	return domain.User{
		Id:       u.Id,
		Email:    u.Email,
		Password: u.Password,
		Nickname: u.Nickname,
		Birthday: time.UnixMilli(u.Birthday),
		AboutMe:  u.AboutMe,
	}
}

func (repo *UserRepository) UpdateNonZeroFields(ctx context.Context, u domain.User) error {
	return repo.dao.UpdateById(ctx, dao.User{
		Id:       u.Id,
		Nickname: u.Nickname,
		Birthday: u.Birthday.UnixMilli(),
		AboutMe:  u.AboutMe,
	})
}

func (repo *UserRepository) FindByID(ctx context.Context, uid int64) (domain.User, error) {
	u, err := repo.dao.FindById(ctx, uid)
	if err != nil {
		return domain.User{}, err
	}
	return repo.toDomain(u), nil
}
