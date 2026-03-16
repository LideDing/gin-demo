## Why

当前项目将所有路由定义、处理函数和启动逻辑集中在 `cmd/main.go` 单文件中，OIDC 认证中间件虽然已独立到 `internal/middleware/` 但路由与业务逻辑高度耦合。随着功能增长，这种结构难以维护、难以测试、难以进行团队协作式的二次开发。需要将其重构为路由与功能函数分离的分层架构，使其成为一个可扩展的生产级 API 框架。

## What Changes

- 引入分层目录结构：`router/`（路由注册）、`handler/`（请求处理）、`service/`（业务逻辑）、`model/`（数据模型）、`config/`（配置管理）
- 将路由定义从 `cmd/main.go` 剥离到独立的 `router/` 包，按模块分组注册路由
- 将 handler（控制器）函数从 main.go 剥离到 `handler/` 包，每个功能模块一个文件
- 引入统一的配置管理模块，替代散落在 main.go 中的 `getEnv` 方式
- 引入统一的响应格式封装（成功/失败统一结构体）
- 保持现有 OIDC 中间件不变，但调整其在新架构中的注册方式
- 引入优雅关机（graceful shutdown）机制
- 为二次开发者提供清晰的模块扩展点和示例

## Capabilities

### New Capabilities
- `api-router`: 独立的路由注册层，支持路由分组、模块化注册，公开路由与受保护路由分离
- `api-handler`: Handler 层封装，定义请求处理函数的标准结构与统一响应格式
- `app-config`: 集中式配置管理，支持环境变量加载、配置校验、结构化配置
- `app-bootstrap`: 应用启动引导流程，包含依赖初始化、优雅关机、服务器配置

### Modified Capabilities
（无现有 spec 需要修改）

## Impact

- **代码结构**：`cmd/main.go` 将大幅简化为入口引导代码；路由、处理函数、配置分别迁移到独立包
- **API 兼容性**：所有现有 API 端点（`/oidc/login`、`/oidc/callback`、`/oidc/logout`、`/oidc/userinfo`、`/ping`、`/hi`）保持不变
- **依赖**：无新外部依赖，仅调整内部包结构
- **中间件**：`internal/middleware/oidc.go` 保持不变，仅调整调用方式
- **二次开发**：重构后新增业务模块只需在 `handler/` 添加处理函数、在 `router/` 注册路由，无需修改框架核心代码
