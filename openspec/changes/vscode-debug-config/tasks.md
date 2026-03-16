## 1. 配置 VSCode 调试文件

- [x] 1.1 编写 `.vscode/launch.json`，添加 Go 调试配置（type: go, request: launch），配置名称为 "启动 Gin 服务"，program 指向 `${workspaceFolder}/cmd`，envFile 指向 `${workspaceFolder}/.env`

## 2. 验证

- [x] 2.1 确认 launch.json 为有效 JSON 格式，且包含所有必要字段（version、configurations、name、type、request、mode、program、envFile）
