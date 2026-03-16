# Gin + OIDC 认证示例

这是一个使用 Gin 框架和 OIDC (OpenID Connect) 实现用户认证的示例项目。

## 功能特性

- ✅ 标准 OIDC 认证流程
- ✅ 支持多种 OIDC Provider（Google, Azure AD, Keycloak, Okta 等）
- ✅ 会话管理
- ✅ 受保护的 API 端点
- ✅ 用户信息获取
- ✅ 登出功能

## 快速开始

### 1. 在 OIDC Provider 中注册应用

首先，你需要在你的 OIDC Provider 中注册一个应用。以下是常见 Provider 的注册指南：

#### Google

1. 访问 [Google Cloud Console](https://console.cloud.google.com/)
2. 创建新项目或选择现有项目
3. 启用 "Google+ API"
4. 创建 OAuth 2.0 凭据
5. 添加授权重定向 URI：`http://localhost:8080/auth/callback`
6. 获取 Client ID 和 Client Secret

#### Azure AD

1. 访问 [Azure Portal](https://portal.azure.com/)
2. 进入 "Azure Active Directory" > "App registrations"
3. 点击 "New registration"
4. 添加重定向 URI：`http://localhost:8080/auth/callback`
5. 创建 Client Secret
6. 记录 Application (client) ID 和 Directory (tenant) ID

#### Keycloak

1. 登录 Keycloak Admin Console
2. 选择或创建 Realm
3. 创建新的 Client
4. 设置 Access Type 为 "confidential"
5. 添加 Valid Redirect URIs：`http://localhost:8080/auth/callback`
6. 获取 Client ID 和 Client Secret

### 2. 配置环境变量

复制 `.env.example` 到 `.env` 并填入你的配置：

```bash
cp .env.example .env
```

编辑 `.env` 文件：

```bash
# 示例：使用 Google
OIDC_ISSUER_URL=https://accounts.google.com
OIDC_CLIENT_ID=your-google-client-id.apps.googleusercontent.com
OIDC_CLIENT_SECRET=your-google-client-secret
OIDC_REDIRECT_URL=http://localhost:8080/auth/callback
OIDC_SCOPES=openid,profile,email
```

### 3. 安装依赖

```bash
go mod tidy
```

### 4. 运行应用

```bash
# 加载环境变量
source .env

# 或者使用 export
export $(cat .env | xargs)

# 运行应用
go run cmd/main.go
```

### 5. 测试

访问 `http://localhost:8080/ping`，应该会自动重定向到 OIDC Provider 的登录页面。

## API 端点

### 公开端点

- **`GET /hi`** - 公开端点，无需认证
  ```bash
  curl http://localhost:8080/hi
  ```

### OIDC 认证端点

- **`GET /oidc/login`** - 发起 OIDC 登录流程
- **`GET /auth/callback`** - OIDC 回调地址（由 Provider 调用）
- **`GET /oidc/logout`** - 登出
  ```bash
  curl http://localhost:8080/oidc/logout
  ```

### 受保护的端点（需要认证）

- **`GET /ping`** - 受保护的端点，返回用户信息
  ```bash
  # 需要在浏览器中访问，因为需要 cookie
  ```

- **`GET /oidc/userinfo`** - 获取当前用户信息
  ```bash
  # 需要在浏览器中访问，因为需要 cookie
  ```

## 认证流程

1. 用户访问受保护的端点（如 `/ping`）
2. 中间件检查是否有有效的会话
3. 如果没有，重定向到 `/oidc/login`
4. 应用重定向到 OIDC Provider 的授权页面
5. 用户在 Provider 完成登录
6. Provider 重定向回 `/auth/callback`，带上授权码
7. 应用交换授权码获取 token
8. 验证 ID token 并创建会话
9. 重定向回原始请求的页面
10. 用户可以访问受保护的资源

## 代码结构

```
.
├── cmd/
│   └── main.go                 # 主程序入口
├── internal/
│   └── middleware/
│       ├── oidc.go            # OIDC 认证中间件
│       └── saml.go            # SAML 认证中间件（旧）
├── .env.example               # 环境变量示例
└── README.md                  # 本文档
```

## 为其他路由添加 OIDC 认证

如果你想为其他路由添加 OIDC 认证，只需在路由定义时添加中间件：

```go
// 单个路由
r.GET("/protected", oidcMiddleware.RequireOIDC(), yourHandler)

// 路由组
protected := r.Group("/api")
protected.Use(oidcMiddleware.RequireOIDC())
{
    protected.GET("/users", getUsersHandler)
    protected.POST("/data", postDataHandler)
}
```

## 获取用户信息

在受保护的路由处理函数中，可以获取用户信息：

```go
r.GET("/profile", oidcMiddleware.RequireOIDC(), func(c *gin.Context) {
    userInfo, exists := c.Get("user_info")
    if exists {
        claims := userInfo.(map[string]interface{})
        
        // 获取常见字段
        email := claims["email"]
        name := claims["name"]
        sub := claims["sub"] // 用户唯一标识
        
        c.JSON(200, gin.H{
            "email": email,
            "name":  name,
            "sub":   sub,
        })
    }
})
```

## 生产环境配置

### 1. 使用 HTTPS

生产环境必须使用 HTTPS：

```bash
OIDC_REDIRECT_URL=https://your-domain.com/auth/callback
```

### 2. 更新 Provider 配置

在 OIDC Provider 中更新重定向 URI 为生产环境地址。

### 3. 使用持久化会话存储

当前实现使用内存存储会话，生产环境建议使用：
- Redis
- 数据库
- 分布式缓存

### 4. 配置会话安全

```go
// 在生产环境设置 secure 和 httpOnly cookie
c.SetCookie("session_id", sessionID, maxAge, "/", "your-domain.com", true, true)
```

## 常见 OIDC Provider 配置

### Google

```bash
OIDC_ISSUER_URL=https://accounts.google.com
OIDC_CLIENT_ID=your-client-id.apps.googleusercontent.com
OIDC_CLIENT_SECRET=your-client-secret
OIDC_REDIRECT_URL=http://localhost:8080/auth/callback
OIDC_SCOPES=openid,profile,email
```

### Azure AD

```bash
OIDC_ISSUER_URL=https://login.microsoftonline.com/{tenant-id}/v2.0
OIDC_CLIENT_ID=your-client-id
OIDC_CLIENT_SECRET=your-client-secret
OIDC_REDIRECT_URL=http://localhost:8080/auth/callback
OIDC_SCOPES=openid,profile,email
```

### Keycloak

```bash
OIDC_ISSUER_URL=https://your-keycloak.com/realms/{realm-name}
OIDC_CLIENT_ID=your-client-id
OIDC_CLIENT_SECRET=your-client-secret
OIDC_REDIRECT_URL=http://localhost:8080/auth/callback
OIDC_SCOPES=openid,profile,email
```

### Okta

```bash
OIDC_ISSUER_URL=https://your-domain.okta.com
OIDC_CLIENT_ID=your-client-id
OIDC_CLIENT_SECRET=your-client-secret
OIDC_REDIRECT_URL=http://localhost:8080/auth/callback
OIDC_SCOPES=openid,profile,email
```

## 故障排查

### 1. "Failed to create OIDC provider"

- 检查 `OIDC_ISSUER_URL` 是否正确
- 确保可以访问 `{OIDC_ISSUER_URL}/.well-known/openid-configuration`

### 2. "Invalid state parameter"

- 这是 CSRF 保护机制
- 确保浏览器启用了 cookie
- 检查是否有代理或负载均衡器修改了请求

### 3. "No id_token in response"

- 确保 `OIDC_SCOPES` 包含 `openid`
- 检查 Provider 配置是否正确

### 4. 会话过期

- 当前实现的会话有效期与 access token 一致
- 可以实现 refresh token 机制来延长会话

## 依赖

```go
require (
    github.com/coreos/go-oidc/v3 v3.x.x
    github.com/gin-gonic/gin v1.x.x
    golang.org/x/oauth2 v0.x.x
)
```

## 许可证

MIT

## 相关资源

- [OpenID Connect 规范](https://openid.net/connect/)
- [OAuth 2.0 规范](https://oauth.net/2/)
- [Gin 框架文档](https://gin-gonic.com/docs/)
- [go-oidc 库文档](https://github.com/coreos/go-oidc)
