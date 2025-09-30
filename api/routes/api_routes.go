package routes

import (
	"database/sql"
	"encoding/json"
	"mybot/internal/auth"
	"mybot/pkg/config"
	"net/http"
)

func SetupRoutes(db *sql.DB, cfg *config.Config) http.Handler {
	mux := http.NewServeMux()

	// 设置认证路由
	authMux := SetupAuthRoutes(db)
	mux.Handle("/api/", authMux)

	// 受保护的API路由
	protectedMux := http.NewServeMux()
	protectedMux.HandleFunc("/api/protected/data", protectedDataHandler)

	// 应用认证中间件到受保护的路由
	protectedWithAuth := auth.JWTAuthMiddleware(protectedMux)
	mux.Handle("/api/protected/", protectedWithAuth)

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
