# ImageHub

<div align=center>
   <img src="https://img.shields.io/badge/gin-1.10-blue"/>
   <img src="https://img.shields.io/badge/vue-3.5.13-green"/>
   <img src="https://img.shields.io/badge/MCP-integration-yellow"/>
   <img src="https://img.shields.io/badge/Qdrant-Multimodality-purple"/>
</div>

## 简介

### 功能

搭建了一个存储图片的Web应用，实现的功能如下：

1. 用户注册、登录和修改个人信息
2. 支持用户登录后上传图片，会自动从图片中提取拍摄时间、地点和像素等信息
3. 支持用户自定义标签和AI识图添加标签
4. 支持在浏览器中对图片进行编辑
5. 显示图片缩略图，支持图片轮播
6. 支持按照多种条件搜索，支持AI对话搜索
7. 可以删除和下载图片

### 技术栈

- 前端：Vue(ts) + Element Plus
- 后端：gin(golang) + gorm + viper + zap
- 数据库：PostgreSQL + Redis + Qdrant
- 鉴权：JWT

## 本地开发指南

### 环境推荐

- **go：1.24+**
- **node：v22 / npm：10**

### clone 本项目

```bash
git clone git@github.com:cagoya/ImageHub.git
```

### 安装 PostgreSQL

> 这里的docker部署都没启用映射卷，如认为有必要请自行添加。

**开发和测试环境**的数据库可以选择使用docker部署

```bash
# 或者指定版本
docker pull postgres
docker run -d --name postgre-server -p 5432:5432 -e POSTGRES_PASSWORD="123456" postgres
```

为了方便查询，可一并拉取PgAdmin，这里的邮箱随意，只起到一个用户名的作用。

```bash
docker pull dpage/pgadmin4
docker run -d --name pgadmin4 -p 5433:80 -e 'PGADMIN_DEFAULT_EMAIL=user@example.com' -e 'PGADMIN_DEFAULT_PASSWORD=123456' dpage/pgadmin4
```

### 安装Redis

还是使用docker部署，需要启动键空间通知(Key-space Notifications)机制

```bash
docker pull redis
docker run -d --name redis-server -p 6379:6379 redis redis-server --notify-keyspace-events Ex
```

### 安装Qdrant

```bash
docker pull qdrant/qdrant
docker run -p 6333:6333 -p 6334:6334 qdrant/qdrant
```

### 启动后端

```bash
# 切换到后端目录
cd ./server
# 下载、清理相关依赖
go mod tidy
# 运行后端
go run main.go
```

### 启动前端

```bash
# 切换到前端目录
cd ./web
# 下载相关依赖
npm install
# 运行前端
npm run dev
```

### 启动mcp server

```bash
cd ./mcp
go run main.go
```

## Docker 部署指南

### 使用 Docker Compose

项目已配置好 Docker Compose，可以一键启动所有服务

```bash
# 在项目根目录执行
docker-compose up -d

# 查看服务状态
docker-compose ps

# 查看日志
docker-compose logs -f

# 停止所有服务
docker-compose down
```

### 服务访问

对于部署机，可以使用下述地址访问各项服务：

- **前端**: http://localhost:5173
- **后端 API**: http://localhost:8888/api
- **PostgreSQL**: localhost:5432
- **Redis**: localhost:6379
- **Qdrant(gRpc)**：localhost:6334

对于其它设备（其它PC或移动设备），需要将localhost换成部署机的ip地址，并且只能访问前端（使用了`Nginx`代理）

### 配置

后端服务支持通过环境变量 `CONFIG_FILE` 指定配置文件:

- 开发环境为 `server.yaml`
- 测试环境为 `test.yaml`
- 生产环境(docker)为 `docker.yaml`

前端在开发和测试环境会硬编码后端服务的URL，在生产环境(docker)使用`Nginx`代理，故不需要担心部署机ip的变化