## ADDED Requirements

### Requirement: 应用入口简化
`cmd/main.go` SHALL 仅包含应用启动引导逻辑：加载配置 → 初始化依赖 → 设置路由 → 启动服务器。业务逻辑和路由定义不得出现在 main.go 中。

#### Scenario: main 函数职责
- **WHEN** 查看 `cmd/main.go` 时
- **THEN** main 函数仅包含：调用 `config.LoadConfig()`、创建 `OIDCMiddleware`、调用 `router.SetupRouter()`、启动 HTTP 服务器，总行数不超过 50 行

### Requirement: 优雅关机
应用 SHALL 支持优雅关机。当收到 `SIGINT` 或 `SIGTERM` 信号时，SHALL 等待正在处理的请求完成后再关闭服务器，超时时间为 5 秒。

#### Scenario: 收到 SIGTERM 信号
- **WHEN** 应用收到 `SIGTERM` 信号时
- **THEN** 日志输出关机信息，等待进行中的请求完成（最多 5 秒），然后优雅退出

#### Scenario: 关机超时
- **WHEN** 优雅关机等待超过 5 秒时
- **THEN** 强制关闭服务器并退出

### Requirement: 启动信息输出
应用启动时 SHALL 输出关键配置信息，包括监听端口、OIDC Issuer URL、Client ID、Redirect URL 和 Scopes。

#### Scenario: 启动日志
- **WHEN** 应用成功启动时
- **THEN** 在日志中输出服务器监听地址和 OIDC 配置摘要信息

### Requirement: 依赖初始化顺序
应用启动 SHALL 按以下顺序初始化：1）加载配置；2）创建 OIDC 中间件；3）设置路由；4）启动 HTTP 服务器。任一步骤失败 SHALL 立即终止并输出错误信息。

#### Scenario: OIDC 中间件初始化失败
- **WHEN** OIDC Provider 不可达导致 `NewOIDCMiddleware` 失败时
- **THEN** 应用输出错误信息后立即退出，不启动 HTTP 服务器
