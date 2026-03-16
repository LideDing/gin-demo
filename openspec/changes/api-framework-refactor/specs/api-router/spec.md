## ADDED Requirements

### Requirement: 路由注册与 Handler 分离
路由注册 SHALL 在独立的 `internal/router/` 包中完成，不得在 `cmd/main.go` 中直接定义路由。所有路由 SHALL 通过 `SetupRouter()` 函数统一注册，该函数返回配置完毕的 `*gin.Engine` 实例。

#### Scenario: SetupRouter 返回完整路由引擎
- **WHEN** 调用 `router.SetupRouter(cfg, oidcMiddleware)` 时
- **THEN** 返回的 `*gin.Engine` 包含所有已注册的路由（公开路由和受保护路由）

### Requirement: 路由按模块分组
路由 SHALL 按功能模块进行分组注册。每个模块 SHALL 提供独立的注册函数 `RegisterXxxRoutes(rg *gin.RouterGroup)`，在 `SetupRouter` 中调用。

#### Scenario: OIDC 路由分组注册
- **WHEN** 调用 `SetupRouter` 时
- **THEN** `/oidc/login`、`/oidc/callback`、`/oidc/logout` 注册为公开路由，`/oidc/userinfo` 注册为受保护路由

#### Scenario: 健康检查路由分组注册
- **WHEN** 调用 `SetupRouter` 时
- **THEN** `/hi` 注册为公开路由，`/ping` 注册为受保护路由

### Requirement: 公开路由与受保护路由分离
路由 SHALL 明确区分公开路由（无需认证）和受保护路由（需要 OIDC 认证）。受保护路由组 SHALL 统一应用 `RequireOIDC()` 中间件。

#### Scenario: 访问公开路由无需认证
- **WHEN** 未认证用户访问 `/hi` 时
- **THEN** 直接返回 200 响应，不触发 OIDC 登录流程

#### Scenario: 访问受保护路由需认证
- **WHEN** 未认证用户访问 `/ping` 时
- **THEN** 触发 OIDC 登录流程（重定向到 OIDC Provider）

### Requirement: 二次开发路由扩展
开发者 SHALL 能够通过以下步骤添加新路由模块：1）在 `handler/` 中创建 handler 文件；2）在 `router/` 中创建或复用路由注册函数；3）在 `SetupRouter` 中调用注册函数。无需修改框架核心代码。

#### Scenario: 添加新业务模块路由
- **WHEN** 开发者需要添加 `/api/v1/users` 路由时
- **THEN** 只需在 `handler/` 添加 `user.go`，在 `router/router.go` 中添加 `RegisterUserRoutes` 调用，无需修改 `cmd/main.go` 或中间件代码
