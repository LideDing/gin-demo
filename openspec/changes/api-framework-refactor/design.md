## Context

当前项目是一个基于 Gin 框架的 OIDC 认证应用（对接腾讯 TAI），所有路由定义和处理函数集中在 `cmd/main.go`（约 110 行），OIDC 中间件已独立到 `internal/middleware/oidc.go`（约 322 行）。项目已可以成功运行，但缺乏清晰的分层结构，不利于二次开发和团队协作。

目标架构需要遵循 Go 社区推荐的 [Standard Go Project Layout](https://github.com/golang-standards/project-layout) 模式，同时保持 Gin 框架的惯用风格。

## Goals / Non-Goals

**Goals:**
- 将路由注册与 handler 函数完全解耦，路由在 `router/` 包定义，handler 在 `handler/` 包实现
- 引入结构化配置管理，替代散落的 `getEnv` 调用
- 建立统一的 API 响应格式（code/message/data 结构）
- 提供清晰的模块扩展模式，让二次开发者可以快速添加新业务模块
- 实现优雅关机，确保生产环境下请求不被中断
- 保持所有现有 API 端点行为不变

**Non-Goals:**
- 不引入数据库 ORM 或 Redis 等存储层（保留扩展点即可）
- 不引入接口（interface）抽象的依赖注入框架（如 wire/dig），保持简洁
- 不修改 OIDC 中间件的核心逻辑
- 不引入日志框架（如 zap/logrus），后续可作为独立变更
- 不实现 API 版本管理机制

## Decisions

### 1. 目录结构设计

采用以下分层结构：

```
gin-tai-login/
├── cmd/
│   └── main.go              # 入口：加载配置 → 初始化依赖 → 启动服务器
├── internal/
│   ├── config/
│   │   └── config.go        # 结构化配置加载与校验
│   ├── handler/
│   │   ├── oidc.go          # OIDC 相关 handler
│   │   ├── health.go        # 健康检查 handler（/ping, /hi）
│   │   └── response.go      # 统一响应封装
│   ├── middleware/
│   │   └── oidc.go          # 现有 OIDC 中间件（不变）
│   ├── router/
│   │   └── router.go        # 路由注册，按模块分组
│   └── service/             # 业务逻辑层（预留，暂为空）
├── go.mod
└── go.sum
```

**理由**：使用 `internal/` 包防止外部直接引用内部实现，符合 Go 最佳实践。Handler 与 Router 分离使得路由表清晰可读，handler 可独立测试。

**替代方案**：
- 将 handler 放在项目根目录的 `api/` 包 → 不符合 Go 的 internal 保护约定
- 使用 `controller/` 命名 → Go 社区更常用 `handler`

### 2. 路由注册模式

采用「路由注册函数」模式，每个模块提供 `RegisterXxxRoutes(rg *gin.RouterGroup)` 函数：

```go
// router/router.go
func SetupRouter(cfg *config.Config, oidcMw *middleware.OIDCMiddleware) *gin.Engine {
    r := gin.Default()
    
    // 公开路由
    public := r.Group("/")
    RegisterHealthRoutes(public)
    RegisterOIDCPublicRoutes(public, oidcMw)
    
    // 受保护路由
    protected := r.Group("/")
    protected.Use(oidcMw.RequireOIDC())
    RegisterOIDCProtectedRoutes(protected, oidcMw)
    RegisterProtectedRoutes(protected)
    
    return r
}
```

**理由**：路由注册函数模式在 Gin 社区广泛使用，二次开发者只需添加新的 `RegisterXxxRoutes` 函数并在 `SetupRouter` 中调用即可。

**替代方案**：
- 自动注册（反射扫描） → 过度工程，Go 偏好显式
- 装饰器/注解模式 → Go 不支持，不惯用

### 3. 统一响应格式

定义标准响应结构体：

```go
type Response struct {
    Code    int         `json:"code"`
    Message string      `json:"message"`
    Data    interface{} `json:"data,omitempty"`
}
```

提供 `Success(c, data)` 和 `Error(c, code, message)` 辅助函数。

**理由**：统一响应格式使前端开发者能一致地处理所有 API 返回值，是生产级 API 的标准实践。

### 4. 配置管理

使用结构体 + 环境变量的方式：

```go
type Config struct {
    Server ServerConfig
    OIDC   middleware.OIDCConfig
}
```

通过 `LoadConfig()` 函数从环境变量加载，保持 `.env.example` 作为配置模板。

**理由**：不引入 viper 等重型配置库，保持依赖最小化。结构体方式提供类型安全和 IDE 提示。

### 5. 优雅关机

使用标准库 `os/signal` + `context` 实现：

```go
srv := &http.Server{Addr: addr, Handler: router}
go srv.ListenAndServe()
// 等待中断信号
quit := make(chan os.Signal, 1)
signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
<-quit
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()
srv.Shutdown(ctx)
```

**理由**：Go 标准库已内置支持，无需额外依赖。生产环境必备功能。

## Risks / Trade-offs

- **[迁移风险] 现有端点行为变更** → 通过保持完全相同的路由路径和响应格式来缓解。`/ping` 和 `/hi` 的响应格式会更新为统一格式，属于微小但不兼容的变更，需在文档中说明。
- **[复杂度增加] 文件数量增多** → 虽然文件变多，但每个文件职责单一，更易理解和维护。通过清晰的目录结构和命名约定来缓解。
- **[内存会话] OIDC 会话仍使用内存存储** → 这是现有设计，本次不改变。后续可引入 Redis 会话存储作为独立变更。
