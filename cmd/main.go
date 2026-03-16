package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"git.woa.com/lideding/gin-tai-login/internal/config"
	"git.woa.com/lideding/gin-tai-login/internal/middleware"
	"git.woa.com/lideding/gin-tai-login/internal/router"
	"github.com/gin-gonic/gin"
)

func main() {
	// 1. 加载配置
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("配置加载失败: %v", err)
	}

	// 2. 设置 Gin 运行模式
	gin.SetMode(cfg.Server.Mode)

	// 3. 创建 OIDC 中间件
	oidcMiddleware, err := middleware.NewOIDCMiddleware(cfg.OIDC)
	if err != nil {
		log.Fatalf("OIDC 中间件初始化失败: %v", err)
	}

	// 4. 设置路由
	r := router.SetupRouter(oidcMiddleware)

	// 5. 输出启动信息
	addr := fmt.Sprintf("0.0.0.0:%s", cfg.Server.Port)
	log.Println("========================================")
	log.Println("Server starting on :" + cfg.Server.Port)
	log.Println("OIDC Configuration:")
	log.Println("  - Issuer URL:", cfg.OIDC.IssuerURL)
	log.Println("  - Client ID:", cfg.OIDC.ClientID)
	log.Println("  - Redirect URL:", cfg.OIDC.RedirectURL)
	log.Println("  - Scopes:", strings.Join(cfg.OIDC.Scopes, ", "))
	log.Println("========================================")
	log.Println("")
	log.Println("⚠️  请确保在 TAI 后台配置的回调地址为:")
	log.Println("   ", cfg.OIDC.RedirectURL)
	log.Println("")

	// 6. 启动 HTTP 服务器（支持优雅关机）
	srv := &http.Server{
		Addr:    addr,
		Handler: r,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("服务器启动失败: %v", err)
		}
	}()

	// 7. 等待中断信号，优雅关机
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("正在关闭服务器...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("服务器关闭异常: %v", err)
	}
	log.Println("服务器已安全退出")
}
