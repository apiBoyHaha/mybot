package auth

import (
	"encoding/json"
	"mybot/internal/user"
	"mybot/pkg/utils"
	"net/http"
)

type Handler struct {
	userService user.Service
}

func NewHandler(userService user.Service) *Handler {
	return &Handler{userService: userService}
}

type LoginResponse struct {
	Token string            `json:"token"`
	User  user.UserResponse `json:"user"`
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "方法不允许", http.StatusMethodNotAllowed)
		return
	}

	var req user.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "无效的请求数据", http.StatusBadRequest)
		return
	}

	// 认证用户
	authenticatedUser, err := h.userService.AuthenticateUser(req.Username, req.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	// 生成JWT令牌
	token, err := utils.GenerateJWT(authenticatedUser.ID, authenticatedUser.Username, authenticatedUser.Role)
	if err != nil {
		http.Error(w, "令牌生成失败", http.StatusInternalServerError)
		return
	}

	// 返回响应
	response := LoginResponse{
		Token: token,
		User: user.UserResponse{
			ID:        authenticatedUser.ID,
			Username:  authenticatedUser.Username,
			Email:     authenticatedUser.Email,
			Role:      authenticatedUser.Role,
			CreatedAt: authenticatedUser.CreatedAt,
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "方法不允许", http.StatusMethodNotAllowed)
		return
	}

	var req user.CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "无效的请求数据", http.StatusBadRequest)
		return
	}

	// 创建用户
	createdUser, err := h.userService.CreateUser(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(createdUser)
}

func (h *Handler) RefreshToken(w http.ResponseWriter, r *http.Request) {
	// 这里可以实现令牌刷新逻辑
	// 暂时返回未实现状态
	http.Error(w, "功能暂未实现", http.StatusNotImplemented)
}
