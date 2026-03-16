## ADDED Requirements

### Requirement: Handler 独立于路由定义
所有请求处理函数 SHALL 定义在 `internal/handler/` 包中，每个功能模块一个文件。Handler 函数 SHALL 接收 `*gin.Context` 参数，不直接依赖路由注册逻辑。

#### Scenario: OIDC Handler 独立文件
- **WHEN** 查看 `internal/handler/oidc.go` 时
- **THEN** 包含 `HandleLogin`、`HandleCallback`、`HandleLogout`、`HandleUserInfo` 等 handler 函数，这些函数委托给 `OIDCMiddleware` 的对应方法

#### Scenario: 健康检查 Handler 独立文件
- **WHEN** 查看 `internal/handler/health.go` 时
- **THEN** 包含 `Ping` 和 `Hi` handler 函数

### Requirement: 统一 API 响应格式
所有 API 端点 SHALL 使用统一的响应结构体返回数据。响应结构体 SHALL 包含 `code`（整型状态码）、`message`（描述信息）和 `data`（响应数据，可选）三个字段。

#### Scenario: 成功响应格式
- **WHEN** API 请求成功时
- **THEN** 返回 `{"code": 0, "message": "success", "data": <具体数据>}` 格式

#### Scenario: 错误响应格式
- **WHEN** API 请求失败时
- **THEN** 返回 `{"code": <错误码>, "message": "<错误描述>"}` 格式，HTTP 状态码与错误类型匹配

### Requirement: 响应辅助函数
`internal/handler/response.go` SHALL 提供 `Success(c *gin.Context, data interface{})` 和 `Error(c *gin.Context, httpCode int, bizCode int, message string)` 两个辅助函数，所有 handler 中 SHALL 使用这些函数返回响应。

#### Scenario: 使用 Success 辅助函数
- **WHEN** handler 需要返回成功响应时
- **THEN** 调用 `response.Success(c, data)` 即可，自动封装为标准格式

#### Scenario: 使用 Error 辅助函数
- **WHEN** handler 需要返回错误响应时
- **THEN** 调用 `response.Error(c, http.StatusBadRequest, 40001, "参数错误")` 即可，自动封装为标准格式

### Requirement: Handler 结构体封装
OIDC 相关的 handler SHALL 封装在 `OIDCHandler` 结构体中，该结构体持有 `*middleware.OIDCMiddleware` 引用。非 OIDC 的通用 handler 可以使用包级函数。

#### Scenario: OIDCHandler 初始化
- **WHEN** 创建 `OIDCHandler` 时
- **THEN** 接收 `*middleware.OIDCMiddleware` 参数，并在内部方法中委托调用
