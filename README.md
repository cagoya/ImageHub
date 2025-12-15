# 部署指南

## 项目简介

### 技术栈

本项目采用前后端分离的开发模式：

- 前端：Vue(TS)
- 后端：Gin(Go)
- 关系型数据库：PostgreSQL
- Key-Value数据库：Redis
- 向量数据库：Qdrant

### 目录说明

- `.github`：github action(CICD)配置文件
- `data`：参考测试数据，包括有或无`exif`信息的`jpg`和`png`格式图片
- `mcp`：mcp server代码实现
- `script`：插入用户信息的脚本
- `server`：后端代码实现
- `web`：前端代码实现

## Docker 一键部署

### 使用 Docker Compose

- 已配置好 Docker Compose，可以一键启动所有服务
- PostgreSQL 中需要的Database在启动镜像时已经创建

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

### 运行导入脚本

因为管理员用户无法直接在前端创建，所以提供了脚本用于在数据库中插入三位用户（包括一位管理员），脚本位于`script`目录下

信息如下：

|用户名|密码|权限|
|-|-|-|
|张三|123456|普通用户|
|李四|123456|普通用户|
|王五|123456|**管理员**|

运行脚本的方法如下，需要**在docker对应的环境下执行**，因为访问的是docker中的数据库容器

1. Windows 下使用`.\setup.bash`
2. Linux（或wsl） 下使用`./setup.sh`

**管理员**比普通用户多了两大功能：

1. 查看其它用户的图片
2. 查看所有用户的信息，并且可以封禁/解封用户

图片数据由于有复杂的约束（关系型数据库、文件系统和向量数据库中都要有），实在不好用脚本插入，还麻烦老师通过前端插入。

### 服务访问

对于部署机，可以使用下述地址访问各项服务：

- **前端**: http://localhost:5173
- 后端 API: http://localhost:8888/api
- PostgreSQL: localhost:5432
- Redis: localhost:6379
- Qdrant(gRpc)：localhost:6334

由于使用了`Nginx`代理，对于非部署机（其它PC或移动设备），**需要将localhost换成部署机的局域网ip地址**

## 可能遇到的问题

### docker 镜像拉取慢

麻烦久等一下，或者对docker进行换源，我使用的镜像源如下：

```yaml
"registry-mirrors": [
   "https://docker.1ms.run",
   "https://docker.xuanyuan.me"
]
```

### 秘钥失效

本项目中共有四个外部api调用：

1. 高德地图/逆地理解析：将图片`exif`信息中的经纬度解析成具体地址
2. 百度/图片标签生成：识别图片并生成标签，用于AI生成标签功能
3. 阿里/多模态向量生成：编码文本和图片为向量，用于（AI对话）图片搜索
4. Deepseek/llm模型api调用：用于AI对话搜索图片

秘钥已一并提交，在`server/config/docker.yaml`中，提交时确认是有效的。**如果以上功能出现问题，应该是由于秘钥失效，麻烦联系下我**

### 端口冲突

部署时用到的端口有：

- 前端：`5173`
- 后端：`8888`
- mcp server：`7777`
- `qdrant`（向量数据库）：`6333`和`6334`
- `postgresql`（关系型数据库）：`5432`
- `Redis`（key-value数据库）：`6379`

**如果出现端口冲突，麻烦暂时关闭导致冲突的服务**