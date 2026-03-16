## 1. Dockerfile

- [x] 1.1 在项目根目录创建 `Dockerfile`，builder 阶段使用 `golang:1.25-alpine`，复制 `go.mod`/`go.sum` 并执行 `go mod download`
- [x] 1.2 在 builder 阶段复制全部源码并执行 `go build` 生成静态链接二进制 `gin-demo`
- [x] 1.3 添加 runtime 阶段，基于 `alpine:3.21`，安装 `ca-certificates`，从 builder 复制 `gin-demo`，设置 `EXPOSE 8080` 和 `ENTRYPOINT`

## 2. .dockerignore

- [x] 2.1 在项目根目录创建 `.dockerignore`，排除 `.git`、`openspec/`、`*.md`、`.env`、`.env.*`、本地二进制 `gin-demo`

## 3. docker-compose.yml

- [x] 3.1 在项目根目录创建 `docker-compose.yml`，定义 `app` 服务，`build: .` 指向本地 Dockerfile
- [x] 3.2 配置 `env_file: .env` 注入 OIDC 环境变量，`ports: "8080:8080"`，`restart: unless-stopped`

## 4. 验证

- [x] 4.1 确认 `docker build -t gin-demo .` 可成功构建镜像，无编译错误
