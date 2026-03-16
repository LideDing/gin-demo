## Context

项目是一个基于 Gin + OIDC 的 Go Web 应用，入口为 `cmd/main.go`，模块路径为 `git.woa.com/lideding/gin-tai-login`，构建产物为 `gin-demo` 二进制文件。当前没有任何容器化配置，运行需要本地 Go 环境并手动配置环境变量。

## Goals / Non-Goals

**Goals:**
- 使用多阶段构建（builder + runtime）最小化生产镜像体积
- 运行时镜像基于 `alpine`，包含必要的 CA 证书（OIDC 需要 HTTPS 通信）
- `docker-compose.yml` 支持从 `.env` 文件读取 OIDC 配置，开箱即用
- 提供 `.dockerignore` 加速构建，避免将源码、openspec 文档等打入镜像

**Non-Goals:**
- 不配置 nginx 反向代理或 TLS 终止
- 不配置多环境（staging/prod）compose 文件
- 不引入数据库或 Redis 等额外服务（本项目无持久化需求）

## Decisions

**多阶段构建**
- Builder 阶段：`golang:1.25-alpine`，执行 `go build`，利用层缓存
- Runtime 阶段：`alpine:3.21`，仅复制编译产物，镜像约 20MB
- 替代方案：`scratch` 镜像更小，但无 shell 和 CA 证书，OIDC HTTPS 调用会失败；`debian` 系镜像过大

**环境变量注入**
- docker-compose 通过 `env_file: .env` 直接加载项目已有的 `.env` 文件，无需重复配置
- 不将任何 secret 硬编码到镜像或 compose 文件中

**Go 模块缓存优化**
- 先单独 `COPY go.mod go.sum` 并 `go mod download`，再复制源码，使依赖层可被 Docker 缓存复用

## Risks / Trade-offs

- [私有 Go 模块] `git.woa.com` 模块路径在构建时可能需要 GONOSUMCHECK / GOFLAGS 配置 → 由于 `go.sum` 已提交，`go mod download` 可离线使用缓存；若在 CI 中构建需配置 GOPROXY/GONOSUMCHECK
- [OIDC Provider 网络] 容器内需能访问 OIDC Issuer URL → 确保容器网络可达外网，或在 compose 中配置 `extra_hosts`
- [时区] alpine 镜像默认 UTC，若日志需要本地时区需手动安装 tzdata → 当前项目无时区强依赖，暂不处理
