#!/bin/bash

# OIDC 配置测试脚本

echo "=== OIDC 配置测试 ==="
echo ""

# 检查环境变量
echo "1. 检查环境变量..."
echo ""

check_env() {
    local var_name=$1
    local var_value=$(eval echo \$$var_name)
    
    if [ -z "$var_value" ]; then
        echo "  ❌ $var_name: 未设置"
        return 1
    else
        echo "  ✅ $var_name: $var_value"
        return 0
    fi
}

all_set=true

check_env "OIDC_ISSUER_URL" || all_set=false
check_env "OIDC_CLIENT_ID" || all_set=false
check_env "OIDC_CLIENT_SECRET" || all_set=false
check_env "OIDC_REDIRECT_URL" || all_set=false

echo ""

if [ "$all_set" = false ]; then
    echo "❌ 缺少必需的环境变量"
    echo ""
    echo "请设置环境变量："
    echo "  source .env"
    echo ""
    echo "或者："
    echo "  export OIDC_ISSUER_URL=https://your-provider.com"
    echo "  export OIDC_CLIENT_ID=your-client-id"
    echo "  export OIDC_CLIENT_SECRET=your-client-secret"
    echo "  export OIDC_REDIRECT_URL=http://localhost:8080/auth/callback"
    exit 1
fi

# 测试 OIDC Provider 连接
echo "2. 测试 OIDC Provider 连接..."
echo ""

WELL_KNOWN_URL="${OIDC_ISSUER_URL}/.well-known/openid-configuration"
echo "  测试 URL: $WELL_KNOWN_URL"
echo ""

if curl -s -f "$WELL_KNOWN_URL" > /dev/null 2>&1; then
    echo "  ✅ OIDC Provider 可访问"
    echo ""
    
    # 显示 Provider 信息
    echo "  Provider 信息："
    curl -s "$WELL_KNOWN_URL" | python3 -m json.tool 2>/dev/null | grep -E "(issuer|authorization_endpoint|token_endpoint)" | head -3
    echo ""
else
    echo "  ❌ 无法访问 OIDC Provider"
    echo ""
    echo "  请检查："
    echo "    1. OIDC_ISSUER_URL 是否正确"
    echo "    2. 网络连接是否正常"
    echo "    3. Provider 是否在线"
    exit 1
fi

# 检查依赖
echo "3. 检查 Go 依赖..."
echo ""

if go list -m github.com/coreos/go-oidc/v3 > /dev/null 2>&1; then
    echo "  ✅ go-oidc 已安装"
else
    echo "  ❌ go-oidc 未安装"
    echo "  运行: go mod tidy"
    exit 1
fi

if go list -m golang.org/x/oauth2 > /dev/null 2>&1; then
    echo "  ✅ oauth2 已安装"
else
    echo "  ❌ oauth2 未安装"
    echo "  运行: go mod tidy"
    exit 1
fi

echo ""

# 配置摘要
echo "4. 配置摘要"
echo ""
echo "  Issuer:      $OIDC_ISSUER_URL"
echo "  Client ID:   $OIDC_CLIENT_ID"
echo "  Redirect:    $OIDC_REDIRECT_URL"
echo "  Scopes:      ${OIDC_SCOPES:-openid,profile,email}"
echo ""

# 提示下一步
echo "✅ 配置检查完成！"
echo ""
echo "下一步："
echo "  1. 确保在 OIDC Provider 中配置了正确的 Redirect URI"
echo "  2. 运行应用: go run cmd/main.go"
echo "  3. 访问: http://localhost:8080/ping"
echo ""
echo "测试端点："
echo "  - 公开端点: http://localhost:8080/hi"
echo "  - 受保护端点: http://localhost:8080/ping"
echo "  - 用户信息: http://localhost:8080/oidc/userinfo"
echo "  - 登出: http://localhost:8080/oidc/logout"
echo ""
