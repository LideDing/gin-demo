## 1. 配置管理模块

- [x] 1.1 创建 `internal/config/config.go`，定义 `Config`、`ServerConfig` 结构体，实现 `LoadConfig()` 函数从环境变量加载配置
- [x] 1.2 实现配置校验逻辑，确保必需的 OIDC 配置项（IssuerURL、ClientID、ClientSecret）非空时返回明确错误
- [x] 1.3 验证配置中使用的环境变量名与 `.env.example` 一致

## 2. 统一响应封装

- [x] 2.1 创建 `internal/handler/response.go`，定义 `Response` 结构体（Code/Message/Data）
- [x] 2.2 实现 `Success(c *gin.Context, data interface{})` 辅助函数
- [x] 2.3 实现 `Error(c *gin.Context, httpCode int, bizCode int, message string)` 辅助函数

## 3. Handler 层实现

- [x] 3.1 创建 `internal/handler/oidc.go`，定义 `OIDCHandler` 结构体，封装 `HandleLogin`、`HandleCallback`、`HandleLogout`、`HandleUserInfo` 方法，内部委托给 `OIDCMiddleware`
- [x] 3.2 创建 `internal/handler/health.go`，实现 `Ping`（受保护，返回用户信息）和 `Hi`（公开）handler 函数，使用统一响应格式

## 4. 路由注册层

- [x] 4.1 创建 `internal/router/router.go`，实现 `SetupRouter(cfg *config.Config, oidcMw *middleware.OIDCMiddleware) *gin.Engine` 函数
- [x] 4.2 实现公开路由组注册：`/hi`、`/oidc/login`、`/oidc/callback`、`/oidc/logout`
- [x] 4.3 实现受保护路由组注册：`/ping`、`/oidc/userinfo`，统一应用 `RequireOIDC()` 中间件

## 5. 应用入口重构

- [x] 5.1 重构 `cmd/main.go`：按顺序调用 `config.LoadConfig()` → `middleware.NewOIDCMiddleware()` → `router.SetupRouter()` → 启动服务器
- [x] 5.2 实现优雅关机：监听 `SIGINT`/`SIGTERM` 信号，使用 `http.Server.Shutdown()` 带 5 秒超时
- [x] 5.3 保留启动信息输出（监听端口、OIDC 配置摘要、回调地址提醒）
- [x] 5.4 删除 main.go 中不再需要的 `getEnv`、`getScopes` 辅助函数（已迁移到 config 模块）

## 6. 验证与清理

- [x] 6.1 确保 `go build ./...` 编译通过，无未使用的导入
- [x] 6.2 验证所有现有端点（`/oidc/login`、`/oidc/callback`、`/oidc/logout`、`/oidc/userinfo`、`/ping`、`/hi`）在新架构下正常工作
- [x] 6.3 创建 `internal/service/` 目录（空包，预留业务逻辑层扩展点）
