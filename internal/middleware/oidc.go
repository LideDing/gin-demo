package middleware

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"net/http"
	"time"

	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
)

// OIDCMiddleware OIDC 认证中间件
type OIDCMiddleware struct {
	provider     *oidc.Provider
	oauth2Config oauth2.Config
	verifier     *oidc.IDTokenVerifier
	sessions     map[string]*OIDCSession // 简单的内存会话存储
}

// OIDCSession 会话信息
type OIDCSession struct {
	IDToken      string
	AccessToken  string
	RefreshToken string
	UserInfo     map[string]interface{}
	ExpiresAt    time.Time
}

// OIDCConfig OIDC 配置
type OIDCConfig struct {
	IssuerURL    string   // OIDC Provider 的 Issuer URL
	ClientID     string   // 客户端 ID
	ClientSecret string   // 客户端密钥
	RedirectURL  string   // 回调地址
	Scopes       []string // 请求的权限范围
}

// NewOIDCMiddleware 创建新的 OIDC 中间件
func NewOIDCMiddleware(config OIDCConfig) (*OIDCMiddleware, error) {
	ctx := context.Background()

	// 初始化 OIDC Provider
	provider, err := oidc.NewProvider(ctx, config.IssuerURL)
	if err != nil {
		return nil, fmt.Errorf("failed to create OIDC provider: %w", err)
	}

	// 配置 OAuth2
	oauth2Config := oauth2.Config{
		ClientID:     config.ClientID,
		ClientSecret: config.ClientSecret,
		RedirectURL:  config.RedirectURL,
		Endpoint:     provider.Endpoint(),
		Scopes:       config.Scopes,
	}

	// 创建 ID Token 验证器
	verifier := provider.Verifier(&oidc.Config{
		ClientID: config.ClientID,
	})

	return &OIDCMiddleware{
		provider:     provider,
		oauth2Config: oauth2Config,
		verifier:     verifier,
		sessions:     make(map[string]*OIDCSession),
	}, nil
}

// RequireOIDC Gin 中间件函数，要求 OIDC 认证
func (om *OIDCMiddleware) RequireOIDC() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 检查会话 cookie
		sessionID, err := c.Cookie("session_id")
		if err != nil || sessionID == "" {
			// 未登录，重定向到登录
			om.HandleLogin(c)
			c.Abort()
			return
		}

		// 验证会话
		session, exists := om.sessions[sessionID]
		if !exists || session.ExpiresAt.Before(time.Now()) {
			// 会话不存在或已过期
			delete(om.sessions, sessionID)
			om.HandleLogin(c)
			c.Abort()
			return
		}

		// 已认证，将用户信息存储到 context 中
		c.Set("oidc_session", session)
		c.Set("user_info", session.UserInfo)
		c.Next()
	}
}

// HandleLogin 处理登录请求
func (om *OIDCMiddleware) HandleLogin(c *gin.Context) {
	// 生成 state 参数（防止 CSRF 攻击）
	state := generateRandomState()

	// 将 state 存储到 cookie 中
	c.SetCookie("oauth_state", state, 600, "/", "", false, true)

	// 保存原始请求的 URL，登录后重定向回来
	c.SetCookie("redirect_after_login", c.Request.URL.String(), 600, "/", "", false, true)

	// 重定向到 OIDC Provider 的授权页面
	authURL := om.oauth2Config.AuthCodeURL(state)
	
	// 打印调试信息
	fmt.Println("========================================")
	fmt.Println("OIDC Login Request:")
	fmt.Printf("  - Auth URL: %s\n", authURL)
	fmt.Printf("  - Redirect URI in config: %s\n", om.oauth2Config.RedirectURL)
	fmt.Println("========================================")
	
	c.Redirect(http.StatusFound, authURL)
}

// HandleCallback 处理 OIDC 回调
func (om *OIDCMiddleware) HandleCallback(c *gin.Context) {
	ctx := context.Background()

	// 检查是否有错误参数
	if errParam := c.Query("error"); errParam != "" {
		errDesc := c.Query("error_description")
		fmt.Printf("⚠️  OIDC 认证失败:\n")
		fmt.Printf("   错误类型: %s\n", errParam)
		fmt.Printf("   错误描述: %s\n", errDesc)
		c.JSON(http.StatusBadRequest, gin.H{
			"error":             errParam,
			"error_description": errDesc,
			"message":           "OIDC 认证失败，请检查配置",
		})
		return
	}

	// 验证 state 参数
	state := c.Query("state")
	savedState, err := c.Cookie("oauth_state")
	if err != nil || state != savedState {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid state parameter"})
		return
	}

	// 清除 state cookie
	c.SetCookie("oauth_state", "", -1, "/", "", false, true)

	// 获取授权码
	code := c.Query("code")
	if code == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No authorization code"})
		return
	}

	// 交换授权码获取 token
	oauth2Token, err := om.oauth2Config.Exchange(ctx, code)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to exchange token"})
		return
	}

	// 提取 ID Token
	rawIDToken, ok := oauth2Token.Extra("id_token").(string)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "No id_token in response"})
		return
	}

	// 验证 ID Token
	idToken, err := om.verifier.Verify(ctx, rawIDToken)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to verify ID token"})
		return
	}

	// 提取用户信息
	var claims map[string]interface{}
	if err := idToken.Claims(&claims); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse claims"})
		return
	}

	// 标准化用户信息（处理字段映射）
	userInfo := normalizeUserInfo(claims)

	// 打印用户信息用于调试
	fmt.Println("========================================")
	fmt.Println("用户认证成功:")
	fmt.Printf("  - Sub (用户唯一标识): %v\n", userInfo["sub"])
	fmt.Printf("  - Username: %v\n", userInfo["username"])
	fmt.Printf("  - Name: %v\n", userInfo["name"])
	if email, ok := userInfo["email"]; ok {
		fmt.Printf("  - Email: %v\n", email)
	}
	fmt.Println("========================================")

	// 创建会话
	sessionID := generateRandomState()
	session := &OIDCSession{
		IDToken:      rawIDToken,
		AccessToken:  oauth2Token.AccessToken,
		RefreshToken: oauth2Token.RefreshToken,
		UserInfo:     userInfo,
		ExpiresAt:    oauth2Token.Expiry,
	}

	// 存储会话
	om.sessions[sessionID] = session

	// 设置会话 cookie
	c.SetCookie("session_id", sessionID, int(time.Until(oauth2Token.Expiry).Seconds()), "/", "", false, true)

	// 获取登录前的 URL
	redirectURL, err := c.Cookie("redirect_after_login")
	if err != nil || redirectURL == "" {
		redirectURL = "/"
	}
	c.SetCookie("redirect_after_login", "", -1, "/", "", false, true)

	// 重定向回原始页面
	c.Redirect(http.StatusFound, redirectURL)
}

// HandleLogout 处理登出
func (om *OIDCMiddleware) HandleLogout(c *gin.Context) {
	// 获取会话 ID
	sessionID, err := c.Cookie("session_id")
	if err == nil && sessionID != "" {
		// 删除会话
		delete(om.sessions, sessionID)
	}

	// 清除 cookie
	c.SetCookie("session_id", "", -1, "/", "", false, true)

	c.JSON(http.StatusOK, gin.H{"message": "Logged out successfully"})
}

// GetUserInfo 获取用户信息（可选，用于获取更多用户信息）
func (om *OIDCMiddleware) GetUserInfo(c *gin.Context) {
	session, exists := c.Get("oidc_session")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Not authenticated"})
		return
	}

	oidcSession := session.(*OIDCSession)
	
	// 构造返回的用户信息，包含标准化字段和原始字段
	response := gin.H{
		"user_info": oidcSession.UserInfo,
		"standardized_fields": gin.H{
			"sub":      oidcSession.UserInfo["sub"],      // 用户唯一标识
			"username": oidcSession.UserInfo["username"], // 用户名
			"name":     oidcSession.UserInfo["name"],     // 显示名称
		},
		"field_mapping": gin.H{
			"description": "TAI OIDC 字段映射关系",
			"mappings": map[string]string{
				"user_name": "username (TAI 特有字段)",
				"sub":       "sub (用户唯一标识，标准 OIDC 字段)",
			},
		},
	}
	
	c.JSON(http.StatusOK, response)
}

// generateRandomState 生成随机 state 字符串
func generateRandomState() string {
	b := make([]byte, 32)
	rand.Read(b)
	return base64.URLEncoding.EncodeToString(b)
}

// normalizeUserInfo 标准化用户信息，处理不同 OIDC Provider 的字段映射
func normalizeUserInfo(claims map[string]interface{}) map[string]interface{} {
	userInfo := make(map[string]interface{})

	// 复制所有原始字段
	for k, v := range claims {
		userInfo[k] = v
	}

	// 标准化 username 字段
	// TAI 使用 user_name，标准 OIDC 使用 preferred_username
	if username, ok := claims["user_name"].(string); ok {
		userInfo["username"] = username
	} else if username, ok := claims["preferred_username"].(string); ok {
		userInfo["username"] = username
	} else if sub, ok := claims["sub"].(string); ok {
		// 如果没有 username，使用 sub 作为备用
		userInfo["username"] = sub
	}

	// 确保 sub 字段存在（用户唯一标识）
	if _, ok := userInfo["sub"]; !ok {
		if username, ok := claims["user_name"].(string); ok {
			userInfo["sub"] = username
		}
	}

	// 标准化 name 字段
	if _, ok := userInfo["name"]; !ok {
		if displayName, ok := claims["display_name"].(string); ok {
			userInfo["name"] = displayName
		} else if username, ok := userInfo["username"].(string); ok {
			userInfo["name"] = username
		}
	}

	return userInfo
}
