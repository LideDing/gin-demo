## Why

当前项目缺少 VSCode 调试配置（`.vscode/launch.json` 为空），开发者无法在 VSCode 中一键启动调试，只能通过命令行手动运行程序。项目依赖 `.env` 文件中的 OIDC 环境变量（如 `OIDC_ISSUER_URL`、`OIDC_CLIENT_ID`、`OIDC_CLIENT_SECRET` 等），调试时需要正确加载这些环境变量才能正常运行。

## What Changes

- 配置 `.vscode/launch.json`，添加 Go 语言调试启动配置
- 程序入口为 `cmd/main.go`，需正确指定 `program` 路径
- 需要在调试配置中加载 `.env` 文件中的环境变量（`OIDC_ISSUER_URL`、`OIDC_CLIENT_ID`、`OIDC_CLIENT_SECRET`、`OIDC_REDIRECT_URL`、`OIDC_SCOPES`、`PORT` 等）
- 确保调试配置支持断点调试、变量查看等标准调试功能

## Capabilities

### New Capabilities
- `vscode-debug`: 提供 VSCode Go 调试启动配置，支持从 `.env` 文件加载环境变量并启动 `cmd/main.go`

### Modified Capabilities

（无）

## Impact

- 新增/修改文件：`.vscode/launch.json`
- 依赖：需要安装 VSCode Go 扩展（`golang.go`）和 Delve 调试器
- 不影响现有代码和运行逻辑，仅为开发工具配置
