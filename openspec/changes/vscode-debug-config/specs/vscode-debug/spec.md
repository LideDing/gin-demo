## ADDED Requirements

### Requirement: VSCode 调试配置文件存在且有效
`.vscode/launch.json` 文件 SHALL 包含有效的 JSON 配置，用于 Go 程序调试。

#### Scenario: launch.json 文件包含 Go 调试配置
- **WHEN** 开发者打开项目并查看 `.vscode/launch.json`
- **THEN** 文件 SHALL 包含至少一个 `type` 为 `go`、`request` 为 `launch` 的调试配置

### Requirement: 调试配置指向正确的程序入口
调试配置的 `program` 字段 SHALL 指向 `${workspaceFolder}/cmd`，即 `cmd/main.go` 所在的 main package 目录。

#### Scenario: 启动调试时运行正确的入口文件
- **WHEN** 开发者按 F5 启动调试
- **THEN** 调试器 SHALL 编译并运行 `cmd/main.go` 作为程序入口

### Requirement: 调试配置自动加载 .env 环境变量
调试配置 SHALL 通过 `envFile` 属性指向 `${workspaceFolder}/.env`，使程序运行时能获取 OIDC 和服务器相关的环境变量。

#### Scenario: 程序启动时读取到 OIDC 环境变量
- **WHEN** 开发者按 F5 启动调试且 `.env` 文件中配置了 `OIDC_ISSUER_URL`、`OIDC_CLIENT_ID` 等变量
- **THEN** 调试进程的环境中 SHALL 包含这些变量及其对应的值

#### Scenario: .env 文件不存在时不阻止启动
- **WHEN** 开发者按 F5 启动调试但 `.env` 文件不存在
- **THEN** 调试器 SHALL 仍然启动程序（程序本身可能因缺少配置而报错，但调试器不应崩溃）

### Requirement: 调试配置名称清晰可识别
调试配置的 `name` 字段 SHALL 使用清晰的中文名称，便于开发者在调试配置下拉列表中识别。

#### Scenario: 调试配置下拉列表显示有意义的名称
- **WHEN** 开发者打开 VSCode 调试面板
- **THEN** SHALL 能看到描述性的配置名称（如 "启动 Gin 服务"）
