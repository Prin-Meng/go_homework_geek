package dao

import (
	"context"
	"errors"
	"github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
	"time"
)

var (
	ErrDuplicateEmail = errors.New("email already exists")
	ErrRecordNotFound = gorm.ErrRecordNotFound
)

type UserDAO struct {
	db *gorm.DB
}

// NewUserDAO NewUserDao 函数创建并返回一个 UserDAO 类型的指针。该指针包含了一个指向 gorm.DB 类型的指针字段 db。
func NewUserDAO(db *gorm.DB) *UserDAO {
	return &UserDAO{db: db}
}

// Insert 方法插入一条用户记录到数据库中，确保数据完整性和一致性。如果邮箱已存在，返回 ErrDuplicateEmail。
func (dao *UserDAO) Insert(ctx context.Context, u User) error {
	now := time.Now().UnixMilli()
	u.Ctime = now
	u.Utime = now
	err := dao.db.WithContext(ctx).Create(&u).Error
	var me *mysql.MySQLError
	if errors.As(err, &me) {
		const duplicateErr uint16 = 1062
		if me.Number == duplicateErr {
			// 用户冲突，邮箱冲突
			return ErrDuplicateEmail
		}
	}
	return err
}

// FindByEmail 方法根据邮箱地址查询用户信息
func (dao *UserDAO) FindByEmail(ctx context.Context, email string) (User, error) {
	// 定义变量 u，类型为 User
	var u User
	// 通过数据库操作，db.WithContext(ctx) 设置上下文，Where 条件查找邮箱，First 找到第一条匹配记录，并将结果赋给 u
	err := dao.db.WithContext(ctx).Where("email =?", email).First(&u).Error
	// 返回查询到的用户信息 u 和空错误信息
	return u, err
}

// UpdateById 根据给定的用户标识更新数据库中的用户信息
func (dao *UserDAO) UpdateById(ctx context.Context, entity User) error {
	// 使用 GORM 的 Model 函数指定要更新的表和条件
	return dao.db.WithContext(ctx).Model(&entity).Where("id =?", entity.Id).
		// 使用 Updates 函数构建更新字段的映射
		Updates(map[string]any{
			"utime":    time.Now().UnixMilli(), // 更新用户的最后更新时间戳
			"nickname": entity.Nickname,        // 更新用户的昵称
			"birthday": entity.Birthday,        // 更新用户的生日
			"about_me": entity.AboutMe,         // 更新用户的个性签名
		}).Error
}

func (dao *UserDAO) FindById(ctx context.Context, id int64) (User, error) {
	var u User
	err := dao.db.WithContext(ctx).Where("id =?", id).First(&u).Error
	return u, err
}

type User struct {
	Id       int64  `gorm:"primaryKey,autoIncrement"`
	Email    string `gorm:"unique"`
	Password string
	Nickname string `gorm:"type=varchar(128)"`
	Birthday int64
	AboutMe  string `gorm:"type=varchar(4096)"`
	Ctime    int64
	Utime    int64
}
