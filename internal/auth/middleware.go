package auth

import (
	"context"
	"mybot/pkg/utils"
	"net/http"
	"strings"
)

// RoleBasedAuthMiddleware 角色权限中间件
func RoleBasedAuthMiddleware(allowedRoles []string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// 从上下文中获取用户信息
			claims, ok := r.Context().Value("user").(*utils.Claims)
			if !ok {
				http.Error(w, "用户信息未找到", http.StatusUnauthorized)
				return
			}

			// 检查用户角色是否在允许的角色列表中
			hasPermission := false
			for _, role := range allowedRoles {
				if claims.Role == role {
					hasPermission = true
					break
				}
			}

			if !hasPermission {
				http.Error(w, "权限不足", http.StatusForbidden)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// JWTAuthMiddleware 现有的JWT认证中间件
func JWTAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenString := extractTokenFromHeader(r)
		if tokenString == "" {
			http.Error(w, "未提供认证令牌", http.StatusUnauthorized)
			return
		}

		claims, err := utils.ValidateJWT(tokenString)
		if err != nil {
			http.Error(w, "无效的认证令牌", http.StatusUnauthorized)
			return
		}

		// 将用户信息存入上下文
		ctx := context.WithValue(r.Context(), "user", claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func extractTokenFromHeader(r *http.Request) string {
	authHeader := r.Header.Get("Authorization")
	return strings.TrimPrefix(authHeader, "Bearer ")
}
