#import "@preview/ori:0.2.3": *
#import "@preview/cuti:0.2.1": show-cn-fakebold

// 伪粗体
#show: show-cn-fakebold
// 设置章节编号样式
#set heading(numbering: numbly("{1:一}、", default: "1.1  "))
#set math.equation(numbering: "(1)")
// 设置图片默认居中
#show image: it => align(center, it)
#show table: it => align(center, it)

#show: ori.with(
  title: "设计报告",
  author: "钟志昊 3230105732",
  semester: "25-26 秋冬",
  date: datetime.today(),
  maketitle: true,
  makeoutline: true,
)

此处用于插入实验报告信息
#v(20em)

= 文档介绍

== 编写目的

本文档描述软件产品功能设计说明书（SRS）的目的是：

1. 定义软件总体要求，作为开发时设计和编码的基础
2. 作为软件总体测试的依据
3. 为软件后续维护和升级提供清晰的参考
4. 为验收系统是否符合预期要求提供判断依据
5. 为需求变更管理提供可追溯的文档支持

== 文档范围

图片管理网站软件产品设计说明书主要是对该系统中功能及功能操作方式进行基本的描述，包含了该模块的总体设计、详细设计、界面设计和数据库设计的详细介绍，包括系统模块图、系统架构图、类图和ER图等内容。

== 读者对象

1. 客户
2. 需求分析人员
3. 开发人员
4. 测试人员
5. 维护人员
6. 项目管理人员

= 项目介绍

== 项目说明

本项目是浙江大学2025-2026秋冬学期《B/S体系软件设计》的课程项目，旨在构建一个图片管理网站，支持用户在登录后上传自己的图片，并可以进行分类、检索、编辑、删除和查看图片等操作。要求网站界面适配PC端和移动端，此外，网站可以扩展对AI识图和AI对话搜索等功能的支持。项目需要提供必要的文档，如设计报告、测试报告和用户手册等（基础的需求文档已由教师提供）。

== 项目背景

随着智能手机、数码相机和社交媒体的普及，个人用户在日常生活中产生的数字图片和照片数量正在以空前的速度增长。这些图片承载着用户珍贵的回忆、重要的工作资料或创作素材。然而，现有的图片管理解决方案普遍存在以下痛点：

1. 数据分散与同步效率低：图片文件分散存储在手机、电脑、云服务等多个平台，缺乏统一、集中的管理入口，导致查找困难和重复存储。
2. 手动分类与组织成本高：传统的管理方式（如文件夹分类）严重依赖用户手动操作。对于海量图片而言，手动整理分类耗时巨大、易错，且往往无法充分利用图片本身包含的丰富元数据。
3. 目前市面上的图片管理网站缺乏对AI技术的集成：如今随着AI技术不断发展，已经可以做到使用AI来分类图片和搜索图片，这样能够极大提升图片管理的效率和准确性。

== 需求概述

=== 功能需求

该项目是一个BS体系软件，按照需求文档的要求，需要实现以下基本功能：

1. 用户注册（需要用户名、密码、邮箱等信息）
2. 用户登录
3. 用户信息修改
4. 图片上传（涉及自动添加信息和手动设置标签）
5. 图片管理（包括编辑、删除、查看）
6. 缩略图生成和展示
7. 多种条件检索图片
8. 支持图片轮播等丰富的展示方式

此外可以扩展以下功能：

1. 支持AI识图
2. 支持AI对话搜索图片

该系统主要用户有两种，*一种是普通用户，另一种是管理员*，其中普通用户仅可以查看自己上传的图片，而管理员可以查看到所有图片，以便进行监管。因此，该项目的功能可以分为两大模块：用户模块和图片模块，如下图所示：

#image("images/模块划分/用户模块图.png")

#image("images/模块划分/图片模块图.png")

其中每个功能可能有不同的具体实现，比如自动添加标签，可以是从图片的EXIF信息中自动提取，也可是由AI识别，或者图片查看可以是查看单张图片，也可以是图片轮播等等，详细内容将在后续章节中介绍。

=== 性能需求

系统应具备良好的总体性能表现，以保证在实际使用过程中提供流畅、稳定的用户体验。具体要求如下：

1. 响应速度：页面加载时间不超过1秒，用户操作响应时间不超过3秒
2. 并发处理：具有一定的并发访问能力，保证用户能够同时上传和查看图片
3. 浏览器兼容性：兼容主流浏览器，包括但不限于手机和电脑端的 Firefox、Chrome、Safari、Edge等，以及微信等应用的内置浏览器
4. 数据准确性：确保前端展示的数据能够实时同步后端数据库最新信息，保证数据的一致性

=== 安全性需求

用户上传的图片可能会涉及隐私问题，因此系统应具备良好的安全性，具体要求如下：

1. 密码要求6字节以上，在数据库中存储加盐哈希值，避免密码泄露
2. 防范SQL注入等常见攻击，比如设置密码的最大尝试次数
3. 对于长时间误操作的登录会话，应该强制登出，避免账户被盗用

=== 可维护性需求

系统应具备良好的可维护性，具体要求如下：
1. 模块化设计，项目结构清晰
2. 各类文档完备清晰，结合图表展示
3. 遵循代码规范，合理注释
4. 完善测试体系，支持单元测试、集成测试等
5. 使用 git 进行代码管理，确保开发、测试和发布流程规范有序

= 总体设计

== 技术介绍

本项目采用前后端分离的Web开发技术，并在后端编写了一系列api供前端调用，同时自己实现一个mcp server，将前端作为 mcp client，实现对AI查找图片的支持。整个项目的技术选型如下所示：

- 前端：Vue3(ts) + Element Plus 组件库 + vue-router + npm包管理器 + vite构建工具
- 后端：gin框架(go) + gorm
- 数据库：MySQL + Redis
- mcp：MCP TypeScript SDK + MCP Python SDK

=== 前端

本项目前端开发语言主要使用 TypeScript (TS)。TS作为 JavaScript (ES6+) 的超集，引入了静态类型系统，能够在编译阶段和编码过程中提供全面的类型检查，有效避免运行时错误，显著增强了代码的可靠性与可维护性。此外，TS与现有的JavaScript生态完全兼容，不需要担心无法使用某些JS库。

本项目前端框架选用 Vue 3。作为一款渐进式用户界面构建框架，Vue 3 的设计理念使其能够适应不同规模的前端开发需求，支持按需集成各类工具与技术。其核心特性包括高效的声明式渲染、细粒度的响应式系统以及强大的组件化能力。本项目将完全使用组合式 API (Composition API) 进行开发，将特定功能的相关逻辑集中在一起，方便阅读和维护。

组件库选用Element Plus，该组件库基于 Vue3 构建，面向设计师与开发者，提供大量高质量、开箱即用的组件，涵盖数据展示、数据录入、导航、反馈等几乎所有后台系统所需的界面元素，是当前 Vue3 生态中最流行、最成熟的企业级 UI 框架之一。

=== 后端

本项目后端开发语言选用 Go (Golang)，该语言由 Google 开源，设计之初即定位于系统级编程语言，强调高效率和高性能。Go 具备语法简洁、编译快速的特性，能够显著提升开发迭代速度。其核心优势在于对高并发的原生支持，这使得它非常适合构建高吞吐量、低延迟的 Web 服务和微服务架构。此外，Go 语言本身内置了包括完备的包管理工具、高效的单元测试框架以及强大的标准库，为快速构建稳定可靠的后端服务提供了坚实的基础。

本项目后端框架选用 Gin。Gin 是一个基于 Go 语言开发的 高性能、轻量级 Web 框架。它继承了 Go 语言在高并发方面的原生优势，并以其极快的路由速度和低内存占用而闻名。Gin 提供了强大的核心功能集，包括灵活的中间件机制、健壮的 HTTP 请求路由（支持路由分组）、JSON 数据绑定和验证、以及内置的 Crash 崩溃处理。这些特性结合在一起，使得 Gin 成为构建高吞吐量 API 服务和快速开发可维护性强的后端应用的理想选择。

=== 数据库

关系型数据库管理系统选用 MySQL。MySQL 是最流行的关系型数据库管理系统之一，由瑞典 MySQL AB 公司开发，目前属于 Oracle 公司。本项目将使用 gorm 与MySQL交互，gorm 提供了一套高度抽象且开发者友好的 API 接口，能够简化复杂的 SQL 数据库操作，实现类型安全的数据交互，从而显著提升后端开发效率和数据库层的代码质量。

此外，本系统还将引入 Redis 作为MySQL的缓冲。Redis 是一个开源、内存中的键值对存储系统，凭借其极低的延迟和高并发读写能力，主要承担数据缓存、会话管理和分布式锁等核心功能。Redis 灵活的数据结构也使其能够胜任消息中间件的职能。本项目将采用 go-redis 与 Redis 进行交互，实现与 Gin 框架的高效集成。

=== 认证与鉴权

本系统的认证与鉴权采用 JWT (JSON Web Token) 技术。JWT 是一种定义紧凑、自包含的开放标准（RFC 7519），用于在客户端与服务端之间安全地传输信息。由于其无状态性和多平台适用性，JWT 能够简洁而高效地保障用户身份验证与权限分发过程，满足现代 Web 应用的访问控制需求。本项目将使用 `github.com/golang-jwt/jwt/v4` 来实现 JWT 的生成、解析和验证逻辑，确保认证机制的安全性和高效性。

=== 第三方 API

考虑到用户交互和数据呈现的便利性，系统地址信息展示将不只使用原始的经度和纬度。本项目选用了高德地图的逆地理编码（Reverse Geocoding）API。该 API 能够根据输入的经纬度坐标，精确解析并返回完整的国内省市区位置信息。

本系统的 AI 图像识别功能将使用阿里云视觉智能开放平台的图像分类 API 服务。该服务基于阿里云成熟的深度学习模型，能够对上传图片进行高精度、多维度的主题识别和内容分析，并自动生成精确的标签集合。

=== mcp

MCP（Model Context Protocol，模型上下文协议）是一种开放协议，旨在实现大型语言模型（LLM） 应用与外部数据源、工具和服务之间的无缝集成。本项目将在前端使用TS实现mcp client，然后使用python 实现 mcp server ，是系统能够支持用户对话查询网站图片，TS和Python均使用官方的MCP SDK。

== 类图

== 系统总体架构图

系统的总体架构如下图所示，其中客户端支持PC和手机等多种设备访问，后端服务器可以接收前端网页发出的HTTP请求并处理，将需要持久化的数据存储在数据库中。

#image("images/架构图/总体架构图.png")

== 前端架构图

该项目前端主要采用 Vue 3 框架，并采用由饿了么团队维护的 Element Plus 作为组件库。包管理器选择`npm`，打包工具选择`webpack`，构建工具选择`vite`。

#image("images/架构图/前端架构图.png")

== 后端架构图

#image("images/架构图/后端架构图.png")

== 系统运行环境

系统运行环境包括硬件环境和软件环境，以下按照客户端和服务器端列举运行环境

=== 客户端

对用户端的硬件要求如下：

#table(
  columns: (120pt, 250pt),
  align: (center, center),
  table.header(
    [名称],
    [要求]
  ),
  [CPU],
  [以Windows为例 9th COREi5以上],
  [内存],
  [2GB以上],
  [硬盘],
  [256GB以上],
  [显示器],
  [全彩显示],
  [键盘，鼠标],
  [可正常使用]
)

对客户端的软件要求如下：
- 操作系统：Win7以上，Linux或MacOS
- 主流的浏览器：如Edge、Firefox和Chrome等

=== 服务器端

对服务器端的硬件要求如下：

#table(
  columns: (120pt, 250pt),
  align: (center, center),
  table.header(
    [名称],
    [要求]
  ),
  [CPU],
  [以Windows为例11thCOREi5以上],
  [内存],
  [4GB以上],
  [硬盘],
  [512GB以上],
)

对服务器端的软件要求如下：

#table(
  columns: (120pt, 200pt, 80pt),
  align: (center, center, center),
  table.header(
    [类型],
    [名称],
    [版本]
  ),
  [操作系统],
  [Windows10/11且拥有wsl，Linux或MacOS],
  [],
  [网站服务器],
  [Nginx],
  [1.28以上],
  [关系型数据库],
  [MySQL],
  [9.0以上],
  [key-value数据库],
  [Redis],
  [8.0以上]
)

上述软件依赖均通过docker部署，最终项目中会提供docker compose

= 数据库设计与ER图

== 逻辑结构设计

== 表结构设计

本系统的数据库表如下：

=== User

User 用于存储用户信息

#table(
  columns: (100pt, 80pt, 100pt, 150pt),
  align: (center, center, center, center),
  table.header(
    [字段名],
    [类型],
    [描述],
    [备注]
  ),
  [id],
  [int],
  [用户id],
  [主键、非空、自增],
  [name],
  [varchar(64)],
  [用户名],
  [非空，唯一],
  [passward_hash],
  [char(60)],
  [密码哈希],
  [非空],
  [email],
  [varchar(64)],
  [邮箱],
  [非空，唯一],
  [is_admin],
  [bool],
  [是否为管理员],
  [非空，默认false],
  [created_time],
  [datetime],
  [创建时间],
  [grom自动设置],
  [updated_time],
  [datetime],
  [更新时间],
  [grom自动设置]
)

=== Image

Image 用于存储图片信息，图片这种大文件不适合直接存数据库，这里认为知晓图片存储的目录，所以只需要存储文件名称，就可以用图片目录+图片名称找到对应的图片文件。

#table(
  columns: (110pt, 100pt, 100pt, 150pt),
  align: (center, center, center, center),
  table.header(
    [字段名],
    [类型],
    [描述],
    [备注]
  ),
  [id],
  [int],
  [图片id],
  [主键、非空、自增],
  [name],
  [varchar(128)],
  [图片名],
  [非空，唯一],
  [time],
  [datetime],
  [拍摄时间],
  [可空],
  [longitude],
  [decimal(10, 6)],
  [经度],
  [可空],
  [latitude],
  [decimal(10, 6)],
  [纬度],
  [可空],
  [province],
  [varchar(16)],
  [省份],
  [可空],
  [city],
  [varchar(16)],
  [市],
  [可空],
  [district],
  [varchar(16)],
  [区],
  [可空],
  [formatted_address],
  [varchar(64)],
  [完整位置],
  [可空],
  [width_pixel],
  [int],
  [宽像素],
  [可空],
  [height_pixel],
  [int],
  [高像素],
  [可空],
  [user_id],
  [int],
  [用户id],
  [非空，外键],
  [created_time],
  [datetime],
  [创建时间],
  [grom自动设置],
  [updated_time],
  [datetime],
  [更新时间],
  [grom自动设置]
)

=== Tag

#table(
  columns: (80pt, 80pt, 100pt, 150pt),
  align: (center, center, center, center),
  table.header(
    [字段名],
    [类型],
    [描述],
    [备注]
  ),
  [id],
  [int],
  [标签id],
  [主键、非空、自增],
  [name],
  [varchar(16)],
  [标签名],
  [非空，唯一],
  [image_id],
  [int],
  [图片id],
  [非空，外键],
  [created_time],
  [datetime],
  [创建时间],
  [grom自动设置],
  [updated_time],
  [datetime],
  [更新时间],
  [grom自动设置]
)

=== thumbnail

本项目中图片存储在后端，在gin框架中配置静态文件目录，以便前端直接根据URL显示图片。本项目会把每位用户的图片放在各自的文件夹（名称为用户id）下存储，原图和缩略图也是分文件夹存储，并且缩略图的名称是在其原图的基础上加上`_thumbnail`后缀，后端大致的目录结构如下所示：

```sh
server/ # 后端目录
|______ public/ # 静态文件目录
|       |______ user_id # 具体的用户id
|               |______ origin # 原图目录
|               |______ thumbnail # 缩略图目录
|______ # 其它目录
```

因为原图和缩略图的路径都是可以由数据库中的图片元信息得到，所以不需要额外在数据库中存储缩略图的路径。此外，由于用户可能上传同名图片，为了避免覆盖，本项目会为每张图片的名称加上一个`uuid`(Universally Unique Identifier)。

== ER 图

#image("images/ER.png")

= 系统接口设计

== 统一说明

系统接口主要用于前后端交互，本项目前端使用`axios`作为http客户端，与后端进行异步交互，而后端则设置一系列api和service，并绑定到相应路由上。这里我们的api一定程度上遵循RESTful规范，目前的设计如下。

注意：由于目前系统仍处于开发阶段，api接口将在编码过程中进行调整，最终的api接口与本文档可能存在出入。

对于入参，本项目在参数较少时会选择使用HTTP的query参数，在参数较多时选择在body中使用json格式传输DTO（数据传输对象）；对于出参，本项目规定统一的返回格式如下所示：

```json
{
  "code": 200, // http 返回码
  "message": "", // 返回的消息
  "data": [] // 返回的数据，可以是任意类型
}
```

其中code是http返回码，与实际返回码保持一致，message是后端返回的提示信息，这两部分的作用比较固定，然后data是后端返回的具体数据，以下接口返回参数说明返回值只解释data部分。

== 用户信息相关接口

用户信息类接口主要使用User表的数据，并对User表进行增删改查，提供的URL均以user开头，具体的api主要有：

=== 用户注册

#table(
  columns: (150pt, 300pt),
  align: (center, center),
  [接口URL],
  [/user/register],
  [HTTP请求类型],
  [POST],
  [请求参数类型],
  [application/json],
  [请求参数示例],
  [```json
  {
    "username": "李华",
    "email": "3230106666@zju.edu.cn",
    "password": "123456"
  }
  ```],
  [请求参数说明],
  [username、email和password分别是用户名、邮箱和密码],
  [返回值示例],
  [```json
  {
    "code": 0,
    "message": "string",
    "data": {
        "id": 0,
        "username": "string",
        "email": "string",
        "isAdmin": true,
        "createTime": "string",
        "updateTime": "string"
    }
  }
  ```],
  [返回值说明],
  [id是用户id，username是用户名，email是用户邮箱，isAdmin表示用户是否是管理员，createdTime和updatedTime分别是用户的创建和修改时间],
  [接口简介],
  [用于用户注册，会在数据库表中添加对应信息，默认不是管理员]
)

=== 用户登录

#table(
  columns: (150pt, 300pt),
  align: (center, center),
  [接口URL],
  [/user/login],
  [HTTP请求类型],
  [POST],
  [请求参数类型],
  [application/json],
  [请求参数示例],
  [```json
  {
    "username": "李华",
    "password": "123456"
  }
  ```],
  [请求参数说明],
  [username和password分别是用户名和密码],
  [返回值示例],
  [```json
  {
    "code": 0,
    "message": "string",
    "data": {
        "token": "string",
        "user": {
            "id": 0,
            "username": "string",
            "email": "string",
            "isAdmin": true,
            "createTime": "string",
            "updateTime": "string"
        }
    }
  }
  ```],
  [返回值说明],
  [token是生成的JWT，user中是用户的各项信息],
  [接口简介],
  [用于用户登录，验证用户输入的密码是否正确，若正确则返回JWT和用户信息，这两个东西将存储在浏览器local storage中以便后续使用]
)

=== 根据ID查询用户

#table(
  columns: (150pt, 300pt),
  align: (center, center),
  [接口URL],
  [/user/getUserById],
  [HTTP请求类型],
  [GET],
  [请求参数类型],
  [query],
  [请求参数示例],
  [id:integer],
  [请求参数说明],
  [id是用户id],
  [返回值示例],
  [```json
  {
    "code": 0,
    "message": "string",
    "data": {
        "id": 0,
        "username": "string",
        "email": "string",
        "isAdmin": true,
        "createTime": "string",
        "updateTime": "string"
    }
  }
  ```],
  [返回值说明],
  [用户的各项信息],
  [接口简介],
  [用于根据ID获取用户信息，仅限管理员使用]
)

=== 查询所有用户

#table(
  columns: (150pt, 300pt),
  align: (center, center),
  [接口URL],
  [/user/getAllUsers],
  [HTTP请求类型],
  [GET],
  [请求参数类型],
  [query],
  [请求参数示例],
  [page:integer\
   pasize:integer],
  [请求参数说明],
  [分页查询的页数和每页数量],
  [返回值示例],
  [```json
  {
    "code": 0,
    "message": "string",
    "data": {
        "token": "string",
        "user": {
            "id": 0,
            "username": "string",
            "email": "string",
            "isAdmin": true,
            "createTime": "string",
            "updateTime": "string"
        }
    }
  }
  ```],
  [返回值说明],
  [新的token和更新后的用户信息],
  [接口简介],
  [用于用户更新自己的个人信息]
)

=== 更新用户信息

#table(
  columns: (150pt, 300pt),
  align: (center, center),
  [接口URL],
  [/user/updateUser],
  [HTTP请求类型],
  [PUT],
  [请求参数类型],
  [application/json],
  [请求参数示例],
  [```json
  {
    "id": 2,
    "username": "王华",
    "email": "222333@qq.com",
    "password": "abcdef"
  }
  ```],
  [请求参数说明],
  [id、username、email和password分别是用户id、新用户名、新邮箱和新密码],
  [返回值示例],
  [```json
  {
    "code": 0,
    "message": "string",
    "data": {
        "token": "string",
        "user": {
            "id": 0,
            "username": "string",
            "email": "string",
            "isAdmin": true,
            "createTime": "string",
            "updateTime": "string"
        }
    }
  }
  ```],
  [返回值说明],
  [token是生成的JWT，user中是用户的各项信息],
  [接口简介],
  [用于用户登录，验证用户输入的密码是否正确，若正确则返回JWT和用户信息，这两个东西将存储在浏览器local storage中以便后续使用]
)

=== 找回密码

#table(
  columns: (150pt, 300pt),
  align: (center, center),
  [接口URL],
  [/user/retrievePassword],
  [HTTP请求类型],
  [POST],
  [请求参数类型],
  [application/json],
  [请求参数示例],
  [```json
  {
    "new_pasword": "11223344",
    "id": 0,
    "username": "李华",
    "email":"3230106666@zju.edu.cn"
  }
  ```],
  [请求参数说明],
  [new_password是新的密码，id、username和eamil是对应账户的id、用户名和邮箱],
  [返回值示例],
  [```json
  {
    "code": 0,
    "message": "string",
    "data": null
  }
  ```],
  [返回值说明],
  [根据返回码判断是否成功],
  [接口简介],
  [用于用户找回自己的密码]
)

=== 授予管理员

#table(
  columns: (150pt, 300pt),
  align: (center, center),
  [接口URL],
  [/user/grantAdmin],
  [HTTP请求类型],
  [PUT],
  [请求参数类型],
  [query],
  [请求参数示例],
  [id:integer],
  [请求参数说明],
  [id是用户id],
  [返回值示例],
  [```json
  {
    "code": 0,
    "message": "string",
    "data": null
  }
  ```],
  [返回值说明],
  [根据返回码判断是否成功],
  [接口简介],
  [用于授予管理员权限，仅限管理员使用]
)

=== 撤销管理员

#table(
  columns: (150pt, 300pt),
  align: (center, center),
  [接口URL],
  [/user/revokeAdmin],
  [HTTP请求类型],
  [PUT],
  [请求参数类型],
  [query],
  [请求参数示例],
  [id:integer],
  [请求参数说明],
  [id是用户id],
  [返回值示例],
  [```json
  {
    "code": 0,
    "message": "string",
    "data": null
  }
  ```],
  [返回值说明],
  [根据返回码判断是否成功],
  [接口简介],
  [用于撤销管理员权限，仅限管理员使用]
)

=== 封禁用户

#table(
  columns: (150pt, 300pt),
  align: (center, center),
  [接口URL],
  [/user/suspend],
  [HTTP请求类型],
  [POST],
  [请求参数类型],
  [query],
  [请求参数示例],
  [id:integer\
    duration:integer],
  [请求参数说明],
  [id是用户id，duration是持续时间（按天数计算）],
  [返回值示例],
  [```json
  {
    "code": 0,
    "message": "string",
    "data": null
  }
  ```],
  [返回值说明],
  [根据返回码判断是否成功],
  [接口简介],
  [用于封禁一个用户一定时长，仅限管理员使用]
)

== 图片类接口

图片类接口主要使用Image和Tag表的数据，并对这两张表进行增删改查，提供的URL均以image开头，具体的api主要有：

=== 上传图片

#table(
  columns: (150pt, 300pt),
  align: (center, center),
  [接口URL],
  [/image/uploadImage],
  [HTTP请求类型],
  [POST],
  [请求参数类型],
  [multipart/form-data],
  [请求参数示例],
  [image:string\<binary\>\
  imageNmae:string],
  [请求参数说明],
  [image是图片的二进制数据，imageName是自定义图片名称],
  [返回值示例],
  [```json
  {
    "code": 0,
    "message": "string",
    "data": null
  }
  ```],
  [返回值说明],
  [根据状态码判断是否成功],
  [接口简介],
  [用于上传图片，后端会将图片通过文件系统保存到对应目录，图片exif信息和其他元信息会被保存到数据库中]
)

=== 设置标签

#table(
  columns: (150pt, 300pt),
  align: (center, center),
  [接口URL],
  [/image/setTags],
  [HTTP请求类型],
  [PUT],
  [请求参数类型],
  [application/json],
  [请求参数示例],
  [```json
  {
    "imageID": 0,
    "tags": []
  }
  ```],
  [请求参数说明],
  [imageID是图片id，tags是想要添加的标签（字符串数组）],
  [返回值示例],
  [```json
  {
    "code": 0,
    "message": "string",
    "data": null
  }
  ```],
  [返回值说明],
  [根据状态码判断是否成功],
  [接口简介],
  [用于为指定图片设置标签，是直接覆盖，所以增删标签均使用这个接口]
)

=== 更新图片

#table(
  columns: (150pt, 300pt),
  align: (center, center),
  [接口URL],
  [/image/updateImage],
  [HTTP请求类型],
  [POST],
  [请求参数类型],
  [multipart/form-data],
  [请求参数示例],
  [image:string\<binary\>\
  imageID:integer],
  [请求参数说明],
  [image是新的图片的二进制数据，imageID是图片id],
  [返回值示例],
  [```json
  {
    "code": 0,
    "message": "string",
    "data": null
  }
  ```],
  [返回值说明],
  [根据状态码判断是否成功],
  [接口简介],
  [用于更新图片，主要是由于用户在前端编辑了图片]
)

=== 删除图片

#table(
  columns: (150pt, 300pt),
  align: (center, center),
  [接口URL],
  [/image/deleteImage],
  [HTTP请求类型],
  [DELETE],
  [请求参数类型],
  [query],
  [请求参数示例],
  [imageID:integer\
  userID:interger],
  [请求参数说明],
  [imageID是图片ID，userID是图片拥有者的用户ID],
  [返回值示例],
  [```json
  {
    "code": 0,
    "message": "string",
    "data": null
  }
  ```],
  [返回值说明],
  [根据状态码判断是否成功],
  [接口简介],
  [用于删除图片]
)

=== 查询所有图片

#table(
  columns: (150pt, 300pt),
  align: (center, center),
  [接口URL],
  [/image/getAllImages],
  [HTTP请求类型],
  [GET],
  [请求参数类型],
  [query],
  [请求参数示例],
  [page:integer\
  pageSize:integer],
  [请求参数说明],
  [分页查询的页数和每页数量],
  [返回值示例],
  [```json
  {
    "code": 0,
    "message": "string",
    "data": {
        "images": [
            {
                "image_id": 0,
                "image_name": "string",
                "image_path": "string",
                "tags": []
            }
        ],
        "pageTotal": 0
    }
  }
  ```],
  [返回值说明],
  [总页数+图片核心信息],
  [接口简介],
  [用于分页查询所有图片，仅限管理员使用]
)

=== 查询个人所有图片

#table(
  columns: (150pt, 300pt),
  align: (center, center),
  [接口URL],
  [/image/getImageByID],
  [HTTP请求类型],
  [GET],
  [请求参数类型],
  [query],
  [请求参数示例],
  [page:integer\
  pageSize:integer],
  [请求参数说明],
  [分页查询的页数和每页数量],
  [返回值示例],
  [```json
  {
    "code": 0,
    "message": "string",
    "data": {
        "images": [
            {
                "image_id": 0,
                "image_name": "string",
                "image_path": "string",
                "tags": []
            }
        ],
        "pageTotal": 0
    }
  }
  ```],
  [返回值说明],
  [总页数+图片核心信息],
  [接口简介],
  [用于分页查询用户自己的所有图片]
)

=== 根据标签查询图片

#table(
  columns: (150pt, 300pt),
  align: (center, center),
  [接口URL],
  [/image/getImageByTag],
  [HTTP请求类型],
  [GET],
  [请求参数类型],
  [query],
  [请求参数示例],
  [page:integer\
  pageSize:integer\
  tag:string],
  [请求参数说明],
  [分页查询的页数和每页数量，tag是查询的标签],
  [返回值示例],
  [```json
  {
    "code": 0,
    "message": "string",
    "data": {
        "images": [
            {
                "image_id": 0,
                "image_name": "string",
                "image_path": "string",
                "tags": []
            }
        ],
        "pageTotal": 0
    }
  }
  ```],
  [返回值说明],
  [总页数+图片核心信息],
  [接口简介],
  [用于分页根据tag查询用户自己的所有图片]
)

=== 根据时间区间查询图片

#table(
  columns: (150pt, 300pt),
  align: (center, center),
  [接口URL],
  [/image/getImageByTimeRange],
  [HTTP请求类型],
  [GET],
  [请求参数类型],
  [query],
  [请求参数示例],
  [page:integer\
  pageSize:integer\
  startTime:string\
  endTime:string],
  [请求参数说明],
  [分页查询的页数和每页数量，以及时间区间的始末],
  [返回值示例],
  [```json
  {
    "code": 0,
    "message": "string",
    "data": {
        "images": [
            {
                "image_id": 0,
                "image_name": "string",
                "image_path": "string",
                "tags": []
            }
        ],
        "pageTotal": 0
    }
  }
  ```],
  [返回值说明],
  [总页数+图片核心信息],
  [接口简介],
  [用于分页根据时间区间查询用户自己的所有图片]
)

=== 根据省份获取图片

#table(
  columns: (150pt, 300pt),
  align: (center, center),
  [接口URL],
  [/image/getImageByProvince],
  [HTTP请求类型],
  [GET],
  [请求参数类型],
  [query],
  [请求参数示例],
  [page:integer\
  pageSize:integer\
  province:string],
  [请求参数说明],
  [分页查询的页数和每页数量],
  [返回值示例],
  [```json
  {
    "code": 0,
    "message": "string",
    "data": {
        "images": [
            {
                "image_id": 0,
                "image_name": "string",
                "image_path": "string",
                "tags": []
            }
        ],
        "pageTotal": 0
    }
  }
  ```],
  [返回值说明],
  [总页数+图片核心信息],
  [接口简介],
  [用于根据省份分页查询用户自己的所有图片，直辖市算作省]
)

=== 根据城市查询图片

#table(
  columns: (150pt, 300pt),
  align: (center, center),
  [接口URL],
  [/image/getImageByCity],
  [HTTP请求类型],
  [GET],
  [请求参数类型],
  [query],
  [请求参数示例],
  [page:integer\
  pageSize:integer\
  district:string],
  [请求参数说明],
  [分页查询的页数和每页数量],
  [返回值示例],
  [```json
  {
    "code": 0,
    "message": "string",
    "data": {
        "images": [
            {
                "image_id": 0,
                "image_name": "string",
                "image_path": "string",
                "tags": []
            }
        ],
        "pageTotal": 0
    }
  }
  ```],
  [返回值说明],
  [总页数+图片核心信息],
  [接口简介],
  [用于根据城市分页查询用户自己的所有图片]
)

=== 根据区县查询图片

#table(
  columns: (150pt, 300pt),
  align: (center, center),
  [接口URL],
  [/image/getImageByDistrict],
  [HTTP请求类型],
  [GET],
  [请求参数类型],
  [query],
  [请求参数示例],
  [page:integer\
  pageSize:integer\
  district:string],
  [请求参数说明],
  [分页查询的页数和每页数量],
  [返回值示例],
  [```json
  {
    "code": 0,
    "message": "string",
    "data": {
        "images": [
            {
                "image_id": 0,
                "image_name": "string",
                "image_path": "string",
                "tags": []
            }
        ],
        "pageTotal": 0
    }
  }
  ```],
  [返回值说明],
  [总页数+图片核心信息],
  [接口简介],
  [用于根据区县分页查询用户自己的所有图片]
)

=== 查询图片详情

#table(
  columns: (150pt, 300pt),
  align: (center, center),
  [接口URL],
  [/image/getImageDetail],
  [HTTP请求类型],
  [GET],
  [请求参数类型],
  [query],
  [请求参数示例],
  [imageID：integer],
  [请求参数说明],
  [图片id],
  [返回值示例],
  [```json
  {
    "code": 0,
    "message": "string",
    "data": {
        "id": 0,
        "name": "string",
        "path": "string",
        "captured_time": "string",
        "latitude": 0,
        "longitude": 0,
        "formatted_address": "string",
        "province": "string",
        "city": "string",
        "district": "string",
        "width_pixel": 0,
        "height_pixel": 0,
        "user_id": 0,
        "user": {
            "id": 0,
            "username": "string",
            "email": "string",
            "isAdmin": true,
            "createTime": "string",
            "updateTime": "string"
        },
        "created_time": "string",
        "updated_time": "string"
    }
  }
  ```],
  [返回值说明],
  [图片的详细信息，包括tag和exif信息等],
  [接口简介],
  [用于根据图片id查询图片的详细信息]
)

=== 获取图片文件

#table(
  columns: (150pt, 300pt),
  align: (center, center),
  [接口URL],
  [/image/getImageFile],
  [HTTP请求类型],
  [GET],
  [请求参数类型],
  [query],
  [请求参数示例],
  [imageID：integer],
  [请求参数说明],
  [图片id],
  [返回值示例],
  [string\<binary\>],
  [返回值说明],
  [真正的图片文件],
  [接口简介],
  [用于获取图片文件，主要是下载图片时需要]
)

= 系统出错设计

== 出错信息

== 补救措施

= 界面原型

此处的部分界面可能还比较简单，后续会进行优化

== 登录界面

== 注册界面

== 网站首页

= 附录

== 项目进度安排



== 备注