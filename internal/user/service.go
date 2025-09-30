package user

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

type Service interface {
	CreateUser(req *CreateUserRequest) (*UserResponse, error)
	AuthenticateUser(username, password string) (*User, error)
	GetUserByID(id int) (*UserResponse, error)
}

type userService struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &userService{repo: repo}
}

func (s *userService) CreateUser(req *CreateUserRequest) (*UserResponse, error) {
	// 检查用户名是否已存在
	existingUser, _ := s.repo.GetUserByUsername(req.Username)
	if existingUser != nil {
		return nil, fmt.Errorf("用户名已存在")
	}

	// 加密密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("密码加密失败: %v", err)
	}

	// 创建用户
	user := &User{
		Username:     req.Username,
		PasswordHash: string(hashedPassword),
		Email:        req.Email,
		Role:         req.Role,
		IsActive:     true,
	}

	if err := s.repo.CreateUser(user); err != nil {
		return nil, err
	}

	return &UserResponse{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		Role:      user.Role,
		CreatedAt: user.CreatedAt,
	}, nil
}

func (s *userService) AuthenticateUser(username, password string) (*User, error) {
	user, err := s.repo.GetUserByUsername(username)
	if err != nil {
		return nil, fmt.Errorf("用户名或密码错误")
	}

	// 验证密码
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return nil, fmt.Errorf("用户名或密码错误")
	}

	return user, nil
}

func (s *userService) GetUserByID(id int) (*UserResponse, error) {
	user, err := s.repo.GetUserByID(id)
	if err != nil {
		return nil, err
	}

	return &UserResponse{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		Role:      user.Role,
		CreatedAt: user.CreatedAt,
	}, nil
}
