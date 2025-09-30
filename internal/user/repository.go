package user

import (
	"database/sql"
	"fmt"
)

type Repository interface {
	CreateUser(user *User) error
	GetUserByUsername(username string) (*User, error)
	GetUserByID(id int) (*User, error)
	UpdateUser(user *User) error
}

type userRepository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &userRepository{db: db}
}

func (r *userRepository) CreateUser(user *User) error {
	query := `INSERT INTO users (username, password_hash, email, role) VALUES (?, ?, ?, ?)`
	result, err := r.db.Exec(query, user.Username, user.PasswordHash, user.Email, user.Role)
	if err != nil {
		return fmt.Errorf("创建用户失败: %v", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("获取用户ID失败: %v", err)
	}

	user.ID = int(id)
	return nil
}

func (r *userRepository) GetUserByUsername(username string) (*User, error) {
	query := `SELECT id, username, password_hash, email, role, is_active, created_at, updated_at 
	          FROM users WHERE username = ? AND is_active = true`
	row := r.db.QueryRow(query, username)

	user := &User{}
	err := row.Scan(&user.ID, &user.Username, &user.PasswordHash, &user.Email,
		&user.Role, &user.IsActive, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("用户不存在")
		}
		return nil, fmt.Errorf("查询用户失败: %v", err)
	}

	return user, nil
}

func (r *userRepository) GetUserByID(id int) (*User, error) {
	query := `SELECT id, username, password_hash, email, role, is_active, created_at, updated_at 
	          FROM users WHERE id = ? AND is_active = true`
	row := r.db.QueryRow(query, id)

	user := &User{}
	err := row.Scan(&user.ID, &user.Username, &user.PasswordHash, &user.Email,
		&user.Role, &user.IsActive, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("用户不存在")
		}
		return nil, fmt.Errorf("查询用户失败: %v", err)
	}

	return user, nil
}

func (r *userRepository) UpdateUser(user *User) error {
	query := `UPDATE users SET username = ?, email = ?, role = ?, is_active = ? WHERE id = ?`
	_, err := r.db.Exec(query, user.Username, user.Email, user.Role, user.IsActive, user.ID)
	if err != nil {
		return fmt.Errorf("更新用户失败: %v", err)
	}
	return nil
}
