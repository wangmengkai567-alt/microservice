package service

import (
	"errors"
	"user-service/internal/model"
	"user-service/internal/pkg"
	"user-service/internal/repository"
)

type UserService struct {
	repo repository.UserRepo
}

func NewUserService(repo repository.UserRepo) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) Register(username, email, password string) error {
	// 输入验证
	if username == "" || email == "" || password == "" {
		return errors.New("username, email and password cannot be empty")
	}
	if len(password) < 6 {
		return errors.New("password must be at least 6 characters")
	}

	hash, err := pkg.HashPassword(password)
	if err != nil {
		return err
	}

	user := &model.User{
		Username: username,
		Email:    email,
		Password: hash,
	}

	return s.repo.Create(user)
}

func (s *UserService) Login(username, password string) (string, error) {
	// 输入验证
	if username == "" || password == "" {
		return "", errors.New("username and password cannot be empty")
	}

	user, err := s.repo.FindByUsername(username)
	if err != nil {
		return "", errors.New("user not found")
	}

	if err := pkg.CheckPassword(user.Password, password); err != nil {
		return "", errors.New("password incorrect")
	}

	return pkg.GenerateToken(user.ID)
}