package router

import (
	"git.woa.com/lideding/gin-tai-login/internal/handler"
	"git.woa.com/lideding/gin-tai-login/internal/middleware"
	"github.com/gin-gonic/gin"
)

// SetupRouter 配置并返回 Gin 路由引擎
func SetupRouter(oidcMw *middleware.OIDCMiddleware) *gin.Engine {
	r := gin.Default()

	// 创建 OIDC Handler
	oidcHandler := handler.NewOIDCHandler(oidcMw)

	// ========================================
	// 公开路由（无需认证）
	// ========================================
	public := r.Group("/")
	RegisterHealthPublicRoutes(public)
	RegisterOIDCPublicRoutes(public, oidcHandler)

	// ========================================
	// 受保护路由（需要 OIDC 认证）
	// ========================================
	protected := r.Group("/")
	// protected.Use(oidcMw.RequireOIDC())
	if oidcMw != nil {
		protected.Use(oidcMw.RequireOIDC())
	}
	RegisterHealthProtectedRoutes(protected)
	RegisterOIDCProtectedRoutes(protected, oidcHandler)

	return r
}

// RegisterHealthPublicRoutes 注册公开的健康检查路由
func RegisterHealthPublicRoutes(rg *gin.RouterGroup) {
	rg.GET("/hi", handler.Hi)
}

// RegisterHealthProtectedRoutes 注册受保护的健康检查路由
func RegisterHealthProtectedRoutes(rg *gin.RouterGroup) {
	rg.GET("/ping", handler.Ping)
}

// RegisterOIDCPublicRoutes 注册公开的 OIDC 路由
func RegisterOIDCPublicRoutes(rg *gin.RouterGroup, h *handler.OIDCHandler) {
	// oidc := rg.Group("/auth")
	// {
	// 	oidc.GET("/login", h.HandleLogin)
	// 	oidc.GET("/logout", h.HandleLogout)
	// }
	auth := rg.Group("/auth")
	{
		auth.GET("/login", h.HandleLogin)
		auth.GET("/logout", h.HandleLogout)
		auth.GET("/callback", h.HandleCallback)
	}
}

// RegisterOIDCProtectedRoutes 注册受保护的 OIDC 路由
func RegisterOIDCProtectedRoutes(rg *gin.RouterGroup, h *handler.OIDCHandler) {
	oidc := rg.Group("/auth")
	{
		oidc.GET("/userinfo", h.HandleUserInfo)
	}
}
