package routes

import (
	"database/sql"
	"encoding/json"
	"mybot/internal/auth"
	"mybot/pkg/config"
	"net/http"
)

// 1. 用户访问 /api/protected/data
// 2. 请求首先被 /api/protected/ 前缀匹配
// 3. 认证中间件验证Token
// 4. 如果认证通过，请求转发到 protectedMux
// 5. protectedMux 根据具体路径调用 protectedDataHandler

// ... existing code ...
func SetupRoutes(db *sql.DB, cfg *config.Config) http.Handler {
	mux := http.NewServeMux()

	// 设置认证路由
	authMux := SetupAuthRoutes(db)
	mux.Handle("/api/", authMux)

	// 受保护的API路由
	protectedMux := http.NewServeMux()
	protectedMux.HandleFunc("/api/protected/data", protectedDataHandler)

	// 管理员专用路由
	adminMux := http.NewServeMux()
	adminMux.HandleFunc("/api/admin/users", adminUsersHandler)

	// 应用认证中间件到受保护的路由
	protectedWithAuth := auth.JWTAuthMiddleware(protectedMux)
	mux.Handle("/api/protected/", protectedWithAuth)

	// 应用角色权限验证到管理员路由
	adminWithAuth := auth.JWTAuthMiddleware(adminMux)
	adminWithRole := auth.RoleBasedAuthMiddleware([]string{"admin"})(adminWithAuth)
	mux.Handle("/api/admin/", adminWithRole)

	// 健康检查路由
	mux.HandleFunc("/health", healthHandler)

	return mux
}

func protectedDataHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "方法不允许", http.StatusMethodNotAllowed)
		return
	}

	response := map[string]interface{}{
		"message": "这是受保护的数据",
		"status":  "success",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "方法不允许", http.StatusMethodNotAllowed)
		return
	}

	response := map[string]string{
		"status": "healthy",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// adminUsersHandler 处理管理员用户管理功能
func adminUsersHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		// 获取用户列表
		getUsersHandler(w, r)
	case http.MethodPost:
		// 创建新用户
		createUserHandler(w, r)
	case http.MethodPut:
		// 更新用户信息
		updateUserHandler(w, r)
	case http.MethodDelete:
		// 删除用户
		deleteUserHandler(w, r)
	default:
		http.Error(w, "方法不允许", http.StatusMethodNotAllowed)
	}
}

// getUsersHandler 获取用户列表
func getUsersHandler(w http.ResponseWriter, r *http.Request) {
	// 模拟用户数据 - 实际项目中应该从数据库获取
	users := []map[string]interface{}{
		{
			"id":        1,
			"username":  "admin",
			"email":     "admin@example.com",
			"role":      "admin",
			"createdAt": "2024-01-01T00:00:00Z",
		},
		{
			"id":        2,
			"username":  "user1",
			"email":     "user1@example.com",
			"role":      "user",
			"createdAt": "2024-01-02T00:00:00Z",
		},
		{
			"id":        3,
			"username":  "user2",
			"email":     "user2@example.com",
			"role":      "user",
			"createdAt": "2024-01-03T00:00:00Z",
		},
	}

	response := map[string]interface{}{
		"users":  users,
		"total":  len(users),
		"status": "success",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// createUserHandler 创建新用户
func createUserHandler(w http.ResponseWriter, r *http.Request) {
	var userData struct {
		Username string `json:"username"`
		Email    string `json:"email"`
		Password string `json:"password"`
		Role     string `json:"role"`
	}

	if err := json.NewDecoder(r.Body).Decode(&userData); err != nil {
		http.Error(w, "无效的请求数据", http.StatusBadRequest)
		return
	}

	// 验证必填字段
	if userData.Username == "" || userData.Email == "" || userData.Password == "" {
		http.Error(w, "用户名、邮箱和密码为必填项", http.StatusBadRequest)
		return
	}

	// 验证角色
	if userData.Role != "admin" && userData.Role != "user" {
		userData.Role = "user" // 默认角色
	}

	// 模拟创建用户 - 实际项目中应该保存到数据库
	newUser := map[string]interface{}{
		"id":        4, // 模拟ID
		"username":  userData.Username,
		"email":     userData.Email,
		"role":      userData.Role,
		"createdAt": "2024-09-30T17:51:23Z",
	}

	response := map[string]interface{}{
		"user":    newUser,
		"message": "用户创建成功",
		"status":  "success",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

// updateUserHandler 更新用户信息
func updateUserHandler(w http.ResponseWriter, r *http.Request) {
	var updateData struct {
		UserID   int    `json:"userId"`
		Username string `json:"username,omitempty"`
		Email    string `json:"email,omitempty"`
		Role     string `json:"role,omitempty"`
	}

	if err := json.NewDecoder(r.Body).Decode(&updateData); err != nil {
		http.Error(w, "无效的请求数据", http.StatusBadRequest)
		return
	}

	if updateData.UserID == 0 {
		http.Error(w, "用户ID为必填项", http.StatusBadRequest)
		return
	}

	// 模拟更新用户 - 实际项目中应该更新数据库
	updatedUser := map[string]interface{}{
		"id":        updateData.UserID,
		"username":  updateData.Username,
		"email":     updateData.Email,
		"role":      updateData.Role,
		"updatedAt": "2024-09-30T17:51:23Z",
	}

	response := map[string]interface{}{
		"user":    updatedUser,
		"message": "用户信息更新成功",
		"status":  "success",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// deleteUserHandler 删除用户
func deleteUserHandler(w http.ResponseWriter, r *http.Request) {
	// 从查询参数获取用户ID
	userID := r.URL.Query().Get("id")
	if userID == "" {
		http.Error(w, "用户ID为必填项", http.StatusBadRequest)
		return
	}

	// 模拟删除用户 - 实际项目中应该从数据库删除
	response := map[string]interface{}{
		"message": "用户删除成功",
		"userId":  userID,
		"status":  "success",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
