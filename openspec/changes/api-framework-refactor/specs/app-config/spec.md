## ADDED Requirements

### Requirement: 结构化配置定义
应用配置 SHALL 使用 Go 结构体定义，包含 `ServerConfig`（服务器配置）和 `OIDCConfig`（OIDC 配置）两个子结构。`ServerConfig` SHALL 包含 `Port`（监听端口）和 `Mode`（Gin 运行模式）字段。

#### Scenario: 配置结构体定义
- **WHEN** 查看 `internal/config/config.go` 时
- **THEN** 包含 `Config` 结构体，嵌套 `ServerConfig` 和引用 `middleware.OIDCConfig`

### Requirement: 环境变量加载
配置 SHALL 从环境变量加载，每个配置项有明确的环境变量名和默认值。加载函数 SHALL 为 `LoadConfig() *Config`。

#### Scenario: 加载 OIDC 配置
- **WHEN** 调用 `LoadConfig()` 时
- **THEN** 从 `OIDC_ISSUER_URL`、`OIDC_CLIENT_ID`、`OIDC_CLIENT_SECRET`、`OIDC_REDIRECT_URL`、`OIDC_SCOPES` 环境变量加载 OIDC 配置

#### Scenario: 加载服务器配置
- **WHEN** 调用 `LoadConfig()` 时
- **THEN** 从 `PORT`（默认 8080）和 `GIN_MODE`（默认 debug）环境变量加载服务器配置

#### Scenario: 使用默认值
- **WHEN** 环境变量未设置时
- **THEN** 使用代码中定义的默认值，确保应用可以无配置启动（开发模式）

### Requirement: 配置校验
`LoadConfig` SHALL 校验必需的配置项。当 `OIDC_ISSUER_URL`、`OIDC_CLIENT_ID` 或 `OIDC_CLIENT_SECRET` 为空时，SHALL 返回明确的错误信息。

#### Scenario: 缺少必需 OIDC 配置
- **WHEN** `OIDC_CLIENT_ID` 环境变量为空且无默认值时
- **THEN** `LoadConfig()` 返回错误，错误信息包含缺失的配置项名称

### Requirement: 配置与 .env.example 一致
配置加载使用的环境变量名 SHALL 与项目根目录的 `.env.example` 文件保持一致。

#### Scenario: 环境变量名一致性
- **WHEN** 比较 `config.go` 中读取的环境变量名与 `.env.example` 中的变量名时
- **THEN** 两者完全一致
