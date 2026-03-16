package handler

import (
	"git.woa.com/lideding/gin-tai-login/internal/middleware"
	"github.com/gin-gonic/gin"
)

// OIDCHandler OIDC 相关请求处理器
type OIDCHandler struct {
	oidcMw *middleware.OIDCMiddleware
}

// NewOIDCHandler 创建 OIDC Handler
func NewOIDCHandler(oidcMw *middleware.OIDCMiddleware) *OIDCHandler {
	return &OIDCHandler{oidcMw: oidcMw}
}

// HandleLogin 处理登录请求，委托给 OIDCMiddleware
func (h *OIDCHandler) HandleLogin(c *gin.Context) {
	h.oidcMw.HandleLogin(c)
}

// HandleCallback 处理 OIDC 回调，委托给 OIDCMiddleware
func (h *OIDCHandler) HandleCallback(c *gin.Context) {
	h.oidcMw.HandleCallback(c)
}

// HandleLogout 处理登出请求，委托给 OIDCMiddleware
func (h *OIDCHandler) HandleLogout(c *gin.Context) {
	h.oidcMw.HandleLogout(c)
}

// HandleUserInfo 获取用户信息，委托给 OIDCMiddleware
func (h *OIDCHandler) HandleUserInfo(c *gin.Context) {
	h.oidcMw.GetUserInfo(c)
}
