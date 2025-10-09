# ImageHub

<div align=center>
   <img src="https://img.shields.io/badge/gin-1.10-blue"/>
   <img src="https://img.shields.io/badge/vue-3.5.13-green"/>
   <img src="https://img.shields.io/badge/MCP-integration-yellow"/>
</div>

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