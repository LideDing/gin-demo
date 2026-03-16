package config

import (
	"fmt"
	"os"
	"strings"

	"git.woa.com/lideding/gin-tai-login/internal/middleware"
)

// ServerConfig 服务器配置
type ServerConfig struct {
	Port string // 监听端口
	Mode string // Gin 运行模式（debug/release/test）
}

// Config 应用配置
type Config struct {
	Server ServerConfig
	OIDC   middleware.OIDCConfig
}

// LoadConfig 从环境变量加载配置，校验必需项
func LoadConfig() (*Config, error) {
	cfg := &Config{
		Server: ServerConfig{
			Port: getEnv("PORT", "8080"),
			Mode: getEnv("GIN_MODE", "debug"),
		},
		OIDC: middleware.OIDCConfig{
			IssuerURL:    getEnv("OIDC_ISSUER_URL", ""),
			ClientID:     getEnv("OIDC_CLIENT_ID", ""),
			ClientSecret: getEnv("OIDC_CLIENT_SECRET", ""),
			RedirectURL:  getEnv("OIDC_REDIRECT_URL", fmt.Sprintf("http://127.0.0.1:%s/auth/callback", getEnv("PORT", "8080"))),
			Scopes:       getScopes(getEnv("OIDC_SCOPES", "openid,profile")),
		},
	}

	// 校验必需的 OIDC 配置项
	var missing []string
	if cfg.OIDC.IssuerURL == "" {
		missing = append(missing, "OIDC_ISSUER_URL")
	}
	if cfg.OIDC.ClientID == "" {
		missing = append(missing, "OIDC_CLIENT_ID")
	}
	if cfg.OIDC.ClientSecret == "" {
		missing = append(missing, "OIDC_CLIENT_SECRET")
	}
	if len(missing) > 0 {
		return nil, fmt.Errorf("缺少必需的配置项: %s", strings.Join(missing, ", "))
	}

	return cfg, nil
}

// getEnv 获取环境变量，如果不存在则返回默认值
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

// getScopes 解析 scopes 字符串为数组
func getScopes(scopesStr string) []string {
	if scopesStr == "" {
		return []string{"openid", "profile"}
	}
	return strings.Split(scopesStr, ",")
}
