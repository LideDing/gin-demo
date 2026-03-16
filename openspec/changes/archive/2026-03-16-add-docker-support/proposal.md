## Why

目前项目缺少容器化支持，开发者需要在本地安装 Go 环境才能运行应用。添加 Dockerfile 和 docker-compose 脚本可以实现一键容器化构建和运行，简化部署流程，统一开发与生产环境。

## What Changes

- 新增 `Dockerfile`：多阶段构建，生产镜像基于 `alpine`，最小化镜像体积
- 新增 `docker-compose.yml`：定义应用服务，支持通过 `.env` 文件注入 OIDC 环境变量，映射端口
- 新增 `.dockerignore`：排除不必要文件，加快构建速度

## Capabilities

### New Capabilities
- `docker-build`: Dockerfile 多阶段构建规范，定义镜像构建方式与运行时配置
- `docker-compose-run`: docker-compose 服务编排规范，定义本地开发与部署的服务配置

### Modified Capabilities
<!-- 无现有 spec 需要修改 -->

## Impact

- 新增根目录文件：`Dockerfile`、`docker-compose.yml`、`.dockerignore`
- 不修改任何现有 Go 源码
- 依赖 `.env` 文件传递 OIDC 配置（已在项目中存在）
- 构建产物为 `gin-demo` 二进制，监听端口由 `PORT` 环境变量控制（默认 `8080`）
