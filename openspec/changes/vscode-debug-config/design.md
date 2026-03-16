## Context

项目是一个基于 Gin 框架的 Go Web 应用，使用 OIDC（TAI）进行身份认证。入口文件为 `cmd/main.go`，运行时依赖 `.env` 文件中的多个环境变量（OIDC 配置、服务端口等）。当前 `.vscode/launch.json` 为空文件，开发者无法使用 VSCode 的调试功能。

项目环境变量（来自 `.env`）：
- `OIDC_ISSUER_URL` - OIDC Provider URL
- `OIDC_CLIENT_ID` - 客户端 ID
- `OIDC_CLIENT_SECRET` - 客户端密钥
- `OIDC_REDIRECT_URL` - 回调地址
- `OIDC_SCOPES` - 权限范围
- `PORT` - 服务端口

## Goals / Non-Goals

**Goals:**
- 配置 VSCode launch.json，使开发者可以一键 F5 启动调试
- 正确加载 `.env` 文件中的环境变量到调试进程
- 支持断点调试、变量查看、调用栈等标准 Go 调试功能

**Non-Goals:**
- 不涉及修改应用代码
- 不配置远程调试
- 不配置测试运行相关的 launch 配置

## Decisions

### 1. 使用 `go` 类型的 launch 配置

**选择**：使用 VSCode Go 扩展提供的 `go` launch 类型，`request` 设为 `launch`。

**理由**：这是 VSCode Go 调试的标准方式，与 Delve 调试器深度集成，支持断点、变量查看等功能。

### 2. 使用 `envFile` 属性加载环境变量

**选择**：使用 launch.json 的 `envFile` 属性指向 `${workspaceFolder}/.env`。

**备选方案**：
- 在 `env` 字段中逐个列出环境变量 → 维护成本高，且需要硬编码敏感信息
- 使用 `.env` 文件中的 `export` 前缀 + shell source → VSCode 不支持

**理由**：`envFile` 是 VSCode 原生支持的方式，自动解析 key=value 格式，无需在 launch.json 中暴露敏感信息。

### 3. `.env` 文件格式适配

**注意点**：当前 `.env` 文件中 `PORT` 变量带有 `export` 前缀（`export PORT=8080`），而其他变量没有 `export`。VSCode 的 `envFile` 解析器**支持**带 `export` 前缀的格式，所以不需要修改 `.env` 文件。

### 4. program 路径配置

**选择**：`program` 设为 `${workspaceFolder}/cmd`，指向 `cmd/main.go` 所在目录。

**理由**：Go 调试器需要指定 package 目录而非单个文件，`cmd` 目录即为 main package 所在位置。

## Risks / Trade-offs

- **[`.env` 格式兼容]** → VSCode envFile 支持 `export KEY=VALUE` 和 `KEY=VALUE` 两种格式，当前 `.env` 文件兼容
- **[敏感信息]** → `.env` 文件已在本地，launch.json 仅引用路径，不暴露实际值。需确保 `.env` 在 `.gitignore` 中
- **[Go 扩展依赖]** → 需要开发者安装 VSCode Go 扩展和 Delve 调试器，属于 Go 开发标配
