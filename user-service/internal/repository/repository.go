package repository

import "user-service/internal/model"

// UserRepo 定义用户仓储接口，方便测试时 mock
type UserRepo interface {
	Create(user *model.User) error
	FindByUsername(username string) (*model.User, error)
}
