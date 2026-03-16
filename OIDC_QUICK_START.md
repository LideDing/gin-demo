# TAI OIDC 字段映射关系

## 概述

TAI OIDC Provider 返回的用户信息字段与标准 OIDC 规范略有不同。本文档说明字段映射关系和标准化处理。

## 字段映射表

| TAI 字段 | 标准 OIDC 字段 | 标准化后的字段 | 说明 |
|---------|---------------|--------------|------|
| `user_name` | `preferred_username` | `username` | 用户名 |
| `sub` | `sub` | `sub` | 用户唯一标识（Subject Identifier） |
| - | `name` | `name` | 用户显示名称 |
| - | `email` | `email` | 用户邮箱（如果申请了 email scope） |

## 字段说明

### 核心字段

#### `sub` (Subject Identifier)
- **必需字段**
- 用户的唯一标识符
- 在同一个 OIDC Provider 中，每个用户的 `sub` 是唯一且不变的
- 应该用于关联用户到应用的内部用户系统

#### `user_name` → `username`
- TAI 特有字段
- 包含用户的登录名
- 应用中统一映射为 `username` 字段

#### `name`
- 用户的显示名称
- 用于界面展示
- 如果不存在，可以使用 `username` 作为备用

## 代码实现

### 标准化处理

```go
func normalizeUserInfo(claims map[string]interface{}) map[string]interface{} {
    userInfo := make(map[string]interface{})
    
    // 复制所有原始字段
    for k, v := range claims {
        userInfo[k] = v
    }
    
    // 标准化 username 字段
    if username, ok := claims["user_name"].(string); ok {
        userInfo["username"] = username
    } else if username, ok := claims["preferred_username"].(string); ok {
        userInfo["username"] = username
    } else if sub, ok := claims["sub"].(string); ok {
        userInfo["username"] = sub
    }
    
    // 确保 sub 字段存在
    if _, ok := userInfo["sub"]; !ok {
        if username, ok := claims["user_name"].(string); ok {
            userInfo["sub"] = username
        }
    }
    
    return userInfo
}
```

### 使用示例

```go
// 获取标准化后的用户信息
userInfo, _ := c.Get("user_info")
userMap := userInfo.(map[string]interface{})

// 访问标准化字段
sub := userMap["sub"]           // 用户唯一标识
username := userMap["username"] // 用户名
name := userMap["name"]         // 显示名称

// 原始 TAI 字段仍然保留
userName := userMap["user_name"] // TAI 原始字段
```

## API 响应示例

### `/tai/userinfo` 接口响应

```json
{
  "user_info": {
    "sub": "user123",
    "user_name": "zhangsan",
    "username": "zhangsan",
    "name": "张三",
    "iat": 1707734400,
    "exp": 1707738000
  },
  "standardized_fields": {
    "sub": "user123",
    "username": "zhangsan",
    "name": "张三"
  },
  "field_mapping": {
    "description": "TAI OIDC 字段映射关系",
    "mappings": {
      "user_name": "username (TAI 特有字段)",
      "sub": "sub (用户唯一标识，标准 OIDC 字段)"
    }
  }
}
```

### `/ping` 接口响应

```json
{
  "message": "pong",
  "user": {
    "sub": "user123",
    "username": "zhangsan",
    "name": "张三"
  },
  "full_user_info": {
    "sub": "user123",
    "user_name": "zhangsan",
    "username": "zhangsan",
    "name": "张三"
  }
}
```

## 最佳实践

1. **使用 `sub` 作为用户唯一标识**
   - 在数据库中关联用户时，使用 `sub` 字段
   - 不要使用 `username` 作为唯一标识，因为用户名可能会变更

2. **使用 `username` 展示和查询**
   - 界面展示时优先使用 `name`，其次使用 `username`
   - 搜索和查询时使用 `username`

3. **保留原始字段**
   - 标准化处理后仍保留所有原始字段
   - 方便调试和兼容性处理

4. **处理字段缺失**
   - 始终检查字段是否存在
   - 提供合理的默认值或备用方案

## 配置 Scopes

TAI OIDC 支持的 scopes：

- `openid` - **必需**，标识这是 OIDC 请求
- `profile` - 获取用户基本信息（`user_name`、`name` 等）
- ~~`email`~~ - TAI 目前不支持此 scope

当前配置：
```go
Scopes: []string{"openid", "profile"}
```

## 常见问题

### Q: 为什么需要字段映射？
A: 不同的 OIDC Provider 可能使用不同的字段名。标准化处理可以让应用代码更统一，方便切换不同的 OIDC Provider。

### Q: `sub` 和 `username` 有什么区别？
A: 
- `sub` 是用户的唯一标识符，永远不会改变
- `username` 是用户名，可能会变更
- 应该用 `sub` 来关联用户到应用的内部系统

### Q: 如何获取更多用户信息？
A: 
1. 检查 TAI 管理后台支持哪些 scopes
2. 在配置中添加相应的 scopes
3. 在回调中会自动获取对应的用户信息

## 参考资料

- [OpenID Connect Core 1.0](https://openid.net/specs/openid-connect-core-1_0.html)
- [Standard Claims](https://openid.net/specs/openid-connect-core-1_0.html#StandardClaims)
