# ImageHub

<div align=center>
   <img src="https://img.shields.io/badge/gin-1.10-blue"/>
   <img src="https://img.shields.io/badge/vue-3.5.13-green"/>
   <img src="https://img.shields.io/badge/MCP-integration-yellow"/>
</div>

## 简介

### 功能

搭建了一个存储图片的Web应用，实现的功能如下：

1. 用户注册登录和修改个人信息
2. 支持用户登录后上传图片，会自动从图片中提取拍摄时间、地点和像素等信息
3. 支持用户自定义标签和AI识图添加标签
4. 支持在浏览器中对图片进行编辑
5. 显示图片缩略图，支持图片轮播
6. 支持按照多种条件搜索，支持AI对话搜索

### 技术栈

- 前端：Vue(ts) + Element Plus
- 后端：gin(golang) + gorm
- 数据库：MySQL + Redis
- 鉴权：JWT

## 本地开发指南

### clone 本项目

```bash
git clone git@github.com:cagoya/ImageHub.git
```

### 安装 mysql

**开发环境**的数据库可以跑在docker容器中，这里使用docker部署mysql

```bash
# 或者指定版本
docker pull mysql:latest
docker run -d --name mysql-server -p 3306:3306 -e MYSQL_ROOT_PASSWORD="123456" mysql
```

### 安装Redis

```bash
docker pull redis
docker run -d --name redis-server -p 6379:6379 redis
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

## Docker 部署指南

### 使用 Docker Compose 一键部署

项目已配置好 Docker Compose，可以一键启动所有服务（MySQL、Redis、后端、前端）。

```bash
# 在项目根目录执行
docker-compose up -d

# 查看服务状态
docker-compose ps

# 查看日志
docker-compose logs -f

# 停止所有服务
docker-compose down

# 停止并删除数据卷（注意：会删除数据库数据）
docker-compose down -v
```

### 服务访问地址

- **前端**: http://localhost:5173
- **后端 API**: http://localhost:8888/api
- **MySQL**: localhost:3306
- **Redis**: localhost:6379

### 单独构建镜像

如果需要单独构建前后端镜像：

```bash
# 构建后端镜像
cd server
docker build -t imagehub-backend:latest .

# 构建前端镜像
cd web
docker build -t imagehub-frontend:latest .
```

### 环境变量配置

后端服务支持通过环境变量 `CONFIG_FILE` 指定配置文件，默认为 `server.yaml`。Docker 环境中会自动使用 `docker.yaml` 配置，该配置已设置好容器间的服务名称（mysql、redis）。

### 数据持久化

Docker Compose 已配置数据卷持久化：
- MySQL 数据存储在 `mysql_data` 卷中
- Redis 数据存储在 `redis_data` 卷中
- 后端上传的图片和日志通过挂载目录持久化