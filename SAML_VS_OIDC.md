# SAML vs OIDC 对比指南

## 快速决策

**推荐使用 OIDC，除非：**
- 你的组织强制要求使用 SAML
- 你已经有现成的 SAML 基础设施
- 你需要与传统企业系统集成

## 主要区别

| 特性 | OIDC | SAML |
|------|------|------|
| **协议基础** | OAuth 2.0 + OpenID | XML-based |
| **数据格式** | JSON | XML |
| **配置复杂度** | ⭐⭐ 简单 | ⭐⭐⭐⭐ 复杂 |
| **需要证书** | ❌ 不需要 | ✅ 需要 |
| **移动端支持** | ✅ 优秀 | ⚠️ 一般 |
| **API 友好** | ✅ 是 | ❌ 否 |
| **现代化程度** | ✅ 现代 | ⚠️ 传统 |
| **企业采用** | ✅ 广泛 | ✅ 广泛 |
| **学习曲线** | ⭐⭐ 容易 | ⭐⭐⭐⭐ 困难 |

## 配置对比

### OIDC 配置（简单）

```bash
# 只需要 4 个必需参数
OIDC_ISSUER_URL=https://accounts.google.com
OIDC_CLIENT_ID=your-client-id
OIDC_CLIENT_SECRET=your-client-secret
OIDC_REDIRECT_URL=http://localhost:8080/oidc/callback
```

**优点：**
- ✅ 不需要生成证书
- ✅ 不需要上传 metadata 文件
- ✅ 配置项少
- ✅ 5 分钟即可完成

### SAML 配置（复杂）

```bash
# 需要更多参数和文件
SAML_ENTITY_ID=http://localhost:8080
SAML_ACS_URL=http://localhost:8080/saml/acs
SAML_METADATA_URL=https://your-idp.com/metadata
SAML_CERTIFICATE_PATH=./certs/saml.crt
SAML_PRIVATE_KEY_PATH=./certs/saml.key
```

**额外步骤：**
1. 生成证书和私钥
2. 生成 SP metadata
3. 上传 metadata 到 IdP
4. 下载 IdP metadata
5. 配置证书路径

**缺点：**
- ❌ 需要生成和管理证书
- ❌ 需要交换 metadata 文件
- ❌ 配置项多
- ❌ 可能需要 30 分钟以上

## 代码对比

### OIDC 中间件（简洁）

```go
// 创建中间件 - 简单
oidcMiddleware, err := middleware.NewOIDCMiddleware(middleware.OIDCConfig{
    IssuerURL:    "https://accounts.google.com",
    ClientID:     "your-client-id",
    ClientSecret: "your-client-secret",
    RedirectURL:  "http://localhost:8080/oidc/callback",
    Scopes:       []string{"openid", "profile", "email"},
})

// 使用中间件
r.GET("/ping", oidcMiddleware.RequireOIDC(), handler)
```

### SAML 中间件（复杂）

```go
// 需要先加载证书
cert, err := samlConfig.LoadCertificate()
privateKey, err := samlConfig.LoadPrivateKey()

// 创建中间件 - 复杂
samlMiddleware, err := middleware.NewSAMLMiddleware(middleware.SAMLConfig{
    EntityID:          "http://localhost:8080",
    ACSURL:            "http://localhost:8080/saml/acs",
    MetadataURL:       "https://your-idp.com/metadata",
    Certificate:       cert,
    PrivateKey:        privateKey,
    AllowIDPInitiated: true,
})

// 使用中间件
r.GET("/ping", samlMiddleware.RequireSAML(), handler)
```

## 认证流程对比

### OIDC 流程（标准 OAuth 2.0）

```
1. 用户访问 /ping
2. 重定向到 /oidc/login
3. 重定向到 Provider 授权页面
4. 用户登录
5. 重定向到 /oidc/callback?code=xxx
6. 交换 code 获取 token
7. 验证 ID token
8. 创建会话
9. 重定向回 /ping
```

**特点：**
- ✅ 流程清晰
- ✅ 使用标准 HTTP 重定向
- ✅ 支持 refresh token
- ✅ 易于调试

### SAML 流程（XML 签名）

```
1. 用户访问 /ping
2. 生成 SAML Request（XML）
3. 签名 SAML Request
4. 重定向到 IdP（带 SAMLRequest）
5. 用户登录
6. IdP 生成 SAML Response（XML）
7. IdP 签名 SAML Response
8. POST 到 /saml/acs
9. 验证签名
10. 解析 XML
11. 创建会话
12. 重定向回 /ping
```

**特点：**
- ⚠️ 流程复杂
- ⚠️ 需要 XML 处理
- ⚠️ 需要签名验证
- ⚠️ 调试困难

## 支持的 Provider

### OIDC Provider（广泛）

- ✅ Google
- ✅ Microsoft Azure AD
- ✅ Okta
- ✅ Auth0
- ✅ Keycloak
- ✅ AWS Cognito
- ✅ GitHub
- ✅ GitLab
- ✅ Apple
- ✅ Facebook
- ✅ 几乎所有现代身份提供商

### SAML Provider（企业为主）

- ✅ Microsoft Azure AD
- ✅ Okta
- ✅ OneLogin
- ✅ Ping Identity
- ✅ ADFS
- ✅ Shibboleth
- ⚠️ Google（需要 Google Workspace）
- ❌ GitHub（不支持）
- ❌ 很多消费级服务不支持

## 使用场景

### 适合使用 OIDC

- ✅ 新项目
- ✅ 移动应用
- ✅ SPA（单页应用）
- ✅ API 服务
- ✅ 微服务架构
- ✅ 需要快速集成
- ✅ 需要支持多种 Provider
- ✅ 需要 refresh token

### 适合使用 SAML

- ✅ 企业内部系统
- ✅ 已有 SAML 基础设施
- ✅ 需要与传统系统集成
- ✅ 企业强制要求
- ✅ 需要 IdP 发起的登录
- ✅ 需要复杂的属性映射

## 安全性对比

### OIDC

- ✅ 使用 JWT（JSON Web Token）
- ✅ 内置签名验证
- ✅ 支持 PKCE（增强安全性）
- ✅ 标准化的安全最佳实践
- ✅ 易于实现和审计

### SAML

- ✅ 使用 XML 签名
- ✅ 成熟的安全模型
- ✅ 支持加密
- ⚠️ 配置错误风险高
- ⚠️ XML 签名复杂

**结论：两者都很安全，但 OIDC 更容易正确实现**

## 性能对比

| 指标 | OIDC | SAML |
|------|------|------|
| **Token 大小** | 小（JSON） | 大（XML） |
| **解析速度** | 快 | 慢 |
| **网络传输** | 少 | 多 |
| **CPU 使用** | 低 | 高 |

## 开发体验

### OIDC

```go
// 简单直观
userInfo, _ := c.Get("user_info")
claims := userInfo.(map[string]interface{})
email := claims["email"].(string)
```

### SAML

```go
// 需要处理 XML 和属性映射
session, _ := c.Get("saml_session")
samlSession := session.(samlsp.SessionWithAttributes)
email := samlSession.GetAttribute("email")
```

## 迁移建议

### 从 SAML 迁移到 OIDC

**步骤：**
1. 在 Provider 中注册 OIDC 应用
2. 部署 OIDC 版本到新端点
3. 逐步迁移用户
4. 保持 SAML 端点一段时间
5. 最终关闭 SAML

**优点：**
- ✅ 简化维护
- ✅ 提升性能
- ✅ 更好的开发体验

### 从 OIDC 迁移到 SAML

**通常不推荐，除非：**
- 企业要求
- 需要与传统系统集成

## 实际项目建议

### 新项目（强烈推荐 OIDC）

```bash
# 使用 OIDC
✅ 配置简单
✅ 开发快速
✅ 维护容易
✅ 用户体验好
```

### 企业项目（根据需求选择）

```bash
# 如果企业已有 SAML 基础设施
→ 使用 SAML

# 如果是新系统或有选择权
→ 使用 OIDC
```

### 混合场景（同时支持）

```go
// 可以同时支持两种认证方式
r.GET("/ping", 
    middleware.AuthRequired(oidcMiddleware, samlMiddleware),
    handler)
```

## 总结

| 场景 | 推荐 |
|------|------|
| 新项目 | **OIDC** ⭐⭐⭐⭐⭐ |
| 移动应用 | **OIDC** ⭐⭐⭐⭐⭐ |
| API 服务 | **OIDC** ⭐⭐⭐⭐⭐ |
| 企业内部系统 | OIDC ⭐⭐⭐⭐ 或 SAML ⭐⭐⭐ |
| 传统系统集成 | **SAML** ⭐⭐⭐⭐ |
| 快速原型 | **OIDC** ⭐⭐⭐⭐⭐ |

## 快速开始

### 使用 OIDC（推荐）

```bash
# 查看快速开始指南
cat OIDC_QUICK_START.md

# 或查看详细文档
cat README.md
```

### 使用 SAML

```bash
# 查看 SAML 配置指南
cat SAML_SETUP_GUIDE.md

# 或查看 SAML 文档
cat README_SAML.md
```

## 结论

**对于大多数项目，我们强烈推荐使用 OIDC：**

- ✅ 更简单
- ✅ 更现代
- ✅ 更快速
- ✅ 更易维护
- ✅ 更好的开发体验

**只有在以下情况下才使用 SAML：**

- 企业强制要求
- 已有 SAML 基础设施
- 需要与传统系统集成

---

**当前项目已经从 SAML 改为 OIDC！** 🎉

查看 [OIDC_QUICK_START.md](OIDC_QUICK_START.md) 开始使用。
