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

#set table(
  fill: (x, y) => {
    if calc.even(y){
      gray.lighten(50%)
    }
  }
)

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

本软件产品功能设计说明书（SRS）旨在实现以下目标：

1. 明确定义软件系统的总体要求，为开发阶段的设计和编码工作提供坚实基础
2. 为软件系统的全面测试提供可靠的依据和标准
3. 为系统后续的维护和升级工作提供清晰、完整的参考文档
4. 为系统验收工作提供明确的判断标准，确保系统符合预期要求
5. 建立完善的需求变更管理体系，确保所有变更都有据可查

== 文档范围

本图片管理网站软件产品设计说明书全面描述了系统的功能特性及操作方式，详细阐述了系统的架构设计、接口设计、界面设计和数据库设计等核心内容，并提供了系统模块图、系统架构图和ER图等可视化设计文档。

== 读者对象

本设计说明书面向以下读者群体：

1. 项目客户和利益相关者
2. 需求分析人员
3. 系统开发人员
4. 软件测试人员
5. 系统维护人员
6. 项目管理人员

= 项目介绍

== 项目说明

本项目是浙江大学2025-2026秋冬学期《B/S体系软件设计》课程的实践项目，旨在构建一个功能完善的图片管理网站。该系统支持用户注册登录后上传个人图片，并提供分类、检索、编辑、删除和查看等完整的图片管理功能。网站界面采用响应式设计，能够良好适配PC端和移动端设备。此外，系统还预留了AI图像识别和AI对话搜索等先进功能的扩展接口。项目要求提供完整的技术文档，包括设计报告、测试报告和用户手册等。

== 项目背景

随着智能手机、数码相机和社交媒体的快速发展，个人用户在日常生活中的数字图片和照片数量呈指数级增长。这些图片承载着用户珍贵的回忆、重要的工作资料或创作素材，具有重要的保存价值。然而，现有的图片管理解决方案普遍存在以下问题：

1. 数据分散与同步效率低下：图片文件分散存储在手机、电脑、云服务等多个平台，缺乏统一、集中的管理入口，导致查找困难且容易造成重复存储。
2. 手动分类与组织成本高昂：传统的管理方式（如文件夹分类）严重依赖用户手动操作。面对海量图片，手动整理分类不仅耗时巨大、容易出错，而且往往无法充分利用图片本身包含的丰富元数据信息。
3. AI技术应用不足：目前市面上的图片管理网站普遍缺乏对AI技术的深度集成。随着AI技术的不断发展，已经可以实现使用AI进行图片分类和智能搜索，这将极大提升图片管理的效率和准确性。

== 需求概述

=== 功能需求

本项目是一个基于B/S架构的软件系统，根据需求文档的要求，需要实现以下核心功能：

1. 用户注册（需要用户名、密码、邮箱等基本信息）
2. 用户登录认证
3. 用户信息修改
4. 图片上传（支持自动提取信息和手动设置标签）
5. 图片管理（包括编辑、删除、查看等操作）
6. 缩略图自动生成和展示
7. 多条件组合检索图片
8. 支持图片轮播等丰富的展示方式

此外，系统还支持以下扩展功能：

1. AI图像识别
2. AI对话搜索图片

该系统主要面向两类用户：*普通用户*和*管理员*。普通用户只能查看和管理自己上传的图片，而管理员可以查看和管理所有用户的图片，以便进行系统监管。因此，项目的功能可以分为两大核心模块：用户管理模块和图片管理模块，如下图所示：

#image("images/模块划分/用户模块图.png")

#image("images/模块划分/图片模块图.png")

其中每个功能都有多种具体的实现方式，例如自动添加标签可以从图片的EXIF信息中自动提取，也可以由AI智能识别；图片查看可以是单张图片查看，也可以是图片轮播展示等。详细的技术实现方案将在后续章节中逐一介绍。

=== 性能需求

系统应具备良好的整体性能表现，确保在实际使用过程中为用户提供流畅、稳定的体验。具体要求如下：

1. 响应速度：页面加载时间控制在3秒以内，用户操作响应时间不超过5秒
2. 并发处理：具备一定的并发访问能力，支持多用户同时上传和查看图片
3. 浏览器兼容性：兼容主流浏览器，包括但不限于PC端和移动端的Firefox、Chrome、Safari、Edge等，以及微信等应用的内置浏览器
4. 数据准确性：确保前端展示的数据能够实时同步后端数据库的最新信息，保证数据的一致性和准确性

=== 安全性需求

考虑到用户上传的图片可能涉及个人隐私，系统必须具备完善的安全防护机制。具体要求如下：

1. 密码安全：要求密码长度不少于6位，在数据库中存储加盐哈希值，防止密码泄露
2. 攻击防护：防范SQL注入等常见网络攻击，设置密码错误尝试次数限制
3. 会话管理：对于长时间无操作的登录会话，系统应自动强制登出，防止账户被盗用
4. 输入验证：对前端用户输入进行严格验证，拦截非法字符和恶意代码

=== 可维护性需求

系统应具备良好的可维护性，确保长期稳定运行和持续优化。具体要求如下：

1. 模块化设计：采用清晰的模块化架构，项目结构层次分明
2. 文档完备：提供完整、清晰的技术文档，结合图表进行详细说明
3. 代码规范：严格遵循代码规范，提供充分的注释说明
4. 测试体系：建立完善的测试体系，支持单元测试、集成测试等多种测试方式
5. 版本管理：使用Git进行代码版本管理，确保开发、测试和发布流程的规范有序

= 总体设计

== 技术介绍

本项目采用前后端分离的Web开发架构，后端提供RESTful API供前端调用，同时实现了MCP（Model Context Protocol）服务器，将前端作为MCP客户端，实现对AI智能查找图片功能的支持。系统处理流程如下图所示：

#image("images/流程图.png")

整个项目的技术选型如下：

- 前端：Vue3(TypeScript) + Element Plus组件库 + Vue Router + npm包管理器 + Vite构建工具
- 后端：Gin框架(Go) + GORM
- 数据库：MySQL + Redis
- MCP：MCP TypeScript SDK + MCP Python SDK

=== 前端

本项目前端开发语言主要使用TypeScript（TS）。TS作为JavaScript（ES6+）的超集，引入了静态类型系统，能够在编译阶段和编码过程中提供全面的类型检查，有效避免运行时错误，显著提升代码的可靠性与可维护性。此外，TS与现有的JavaScript生态系统完全兼容，无需担心无法使用某些JS库的问题。

本项目前端框架选用Vue 3。作为一款渐进式用户界面构建框架，Vue 3的设计理念使其能够适应不同规模的前端开发需求，支持按需集成各类工具与技术。其核心特性包括高效的声明式渲染、细粒度的响应式系统以及强大的组件化能力。本项目将完全使用组合式API（Composition API）进行开发，将特定功能的相关逻辑集中在一起，便于阅读和维护。

组件库选用Element Plus，该组件库基于Vue 3构建，面向设计师与开发者，提供大量高质量、开箱即用的组件，涵盖数据展示、数据录入、导航、反馈等几乎所有后台系统所需的界面元素，是当前Vue 3生态中最流行、最成熟的企业级UI框架之一。

=== 后端

本项目后端开发语言选用Go（Golang），该语言由Google开源，设计之初即定位于系统级编程语言，强调高效率和高性能。Go具备语法简洁、编译快速的特性，能够显著提升开发迭代速度。其核心优势在于对高并发的原生支持，这使得它非常适合构建高吞吐量、低延迟的Web服务和微服务架构。此外，Go语言本身内置了完备的包管理工具、高效的单元测试框架以及强大的标准库，为快速构建稳定可靠的后端服务提供了坚实的基础。

本项目后端框架选用Gin。Gin是一个基于Go语言开发的高性能、轻量级Web框架。它继承了Go语言在高并发方面的原生优势，并以其极快的路由速度和低内存占用而闻名。Gin提供了强大的核心功能集，包括灵活的中间件机制、健壮的HTTP请求路由（支持路由分组）、JSON数据绑定和验证、以及内置的Crash崩溃处理。这些特性结合在一起，使得Gin成为构建高吞吐量API服务和快速开发可维护性强的后端应用的理想选择。

=== 数据库

关系型数据库管理系统选用MySQL。MySQL是最流行的关系型数据库管理系统之一，由瑞典MySQL AB公司开发，目前属于Oracle公司。本项目将使用GORM与MySQL进行交互，GORM提供了一套高度抽象且开发者友好的API接口，能够简化复杂的SQL数据库操作，实现类型安全的数据交互，从而显著提升后端开发效率和数据库层的代码质量。

此外，本系统还将引入Redis作为MySQL的缓存层。Redis是一个开源、内存中的键值对存储系统，凭借其极低的延迟和高并发读写能力，主要承担数据缓存、会话管理和分布式锁等核心功能。Redis灵活的数据结构也使其能够胜任消息中间件的职能。本项目将采用go-redis与Redis进行交互，实现与Gin框架的高效集成。

=== 认证与鉴权

本系统的认证与鉴权采用JWT（JSON Web Token）技术。JWT是一种定义紧凑、自包含的开放标准（RFC 7519），用于在客户端与服务端之间安全地传输信息。由于其无状态性和多平台适用性，JWT能够简洁而高效地保障用户身份验证与权限分发过程，满足现代Web应用的访问控制需求。本项目将使用`github.com/golang-jwt/jwt/v4`来实现JWT的生成、解析和验证逻辑，确保认证机制的安全性和高效性。

=== 第三方API

考虑到用户交互和数据呈现的便利性，系统地址信息展示不仅使用原始的经纬度坐标。本项目选用了高德地图的逆地理编码（Reverse Geocoding）API。该API能够根据输入的经纬度坐标，精确解析并返回完整的国内省市区位置信息。

本系统的AI图像识别功能将使用阿里云视觉智能开放平台的图像分类API服务。该服务基于阿里云成熟的深度学习模型，能够对上传图片进行高精度、多维度的主题识别和内容分析，并自动生成精确的标签集合。

=== MCP

MCP（Model Context Protocol，模型上下文协议）是一种开放协议，旨在实现大型语言模型（LLM）应用与外部数据源、工具和服务之间的无缝集成。本项目将在前端使用TypeScript实现MCP客户端，然后使用Python实现MCP服务器，使系统能够支持用户通过对话方式查询网站图片。TypeScript和Python均使用官方的MCP SDK。

== 架构图

=== 系统总体架构图

系统的总体架构如下图所示，其中客户端支持PC和手机等多种设备访问，后端服务器可以接收前端网页发出的HTTP请求并进行处理，将需要持久化的数据存储在数据库中。

#image("images/架构图/总体架构图.png")

=== 前端架构图

该项目前端主要采用Vue 3框架，并采用由饿了么团队维护的Element Plus作为组件库。包管理器选择npm，打包工具选择webpack，构建工具选择Vite。

#image("images/架构图/前端架构图.png")

=== 后端架构图

#image("images/架构图/后端架构图.png")

=== MCP架构图

#image("images/架构图/mcp架构图.png")

== 系统运行环境

系统运行环境包括硬件环境和软件环境，以下按照客户端和服务器端分别列举运行环境要求：

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

- 操作系统：Windows 7以上、Linux或macOS
- 浏览器：主流浏览器，如Edge、Firefox和Chrome等

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
  [Windows 10/11、Linux或macOS],
  [],
  [关系型数据库],
  [MySQL],
  [9.0以上],
  [key-value数据库],
  [Redis],
  [8.0以上],
  [后端环境],
  [Go],
  [1.24以上],
  [前端环境],
  [npm],
  [10.0以上]
)

上述软件依赖均通过Docker部署，最终项目中会提供Docker Compose文件，所以服务端需要有较新版本（v25.0以上）的Docker。

= 数据库设计与ER图

== 数据库表设计

本系统的数据库表设计如下三张，具体说明见下方表格：

1. User(#underline("id"), name, password_hash, email, is_admin, is_suspended, created_time, updated_time)
2. Image(#underline("id"), name, time, longitude, latitude, province, city, district, formatted_address, width_pixel, height_pixel, user_id, created_time, updated_time)
3. Tag(#underline("id"), name, image_id, created_time, updated_time)

=== User

User表用于存储用户的基本信息。

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
  [password_hash],
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
  [is_suspended],
  [bool],
  [是否被封禁],
  [非空，默认false],
  [created_time],
  [datetime],
  [创建时间],
  [GORM自动设置],
  [updated_time],
  [datetime],
  [更新时间],
  [GORM自动设置]
)

=== Image

Image表用于存储图片的元信息。由于图片文件体积较大，不适合直接存储在数据库中，因此本系统采用文件系统存储图片文件，数据库中仅存储图片的文件名。通过图片存储目录和文件名即可定位到对应的图片文件。

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
  [GORM自动设置],
  [updated_time],
  [datetime],
  [更新时间],
  [GORM自动设置]
)

=== Tag

Tag表用于存储图片的标签信息，实现图片的分类管理功能。

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
  [GORM自动设置],
  [updated_time],
  [datetime],
  [更新时间],
  [GORM自动设置]
)

=== 缩略图存储

本项目中图片存储在后端服务器，在Gin框架中配置静态文件目录，以便前端直接根据URL显示图片。本项目将每位用户的图片放在各自的文件夹（名称为用户ID）下存储，原图和缩略图分别存储在不同的文件夹中。后端大致的目录结构如下所示：

```sh
server/ # 后端目录
|______ public/ # 静态文件目录
|       |______ user_id # 具体的用户ID
|               |______ origin # 原图目录
|               |______ thumbnail # 缩略图目录
|______ # 其它目录
```

由于原图和缩略图的路径都可以通过数据库中的图片元信息计算得出，因此不需要额外在数据库中存储缩略图的路径。此外，由于用户可能上传同名图片，为了避免文件覆盖，本项目会为每张图片的名称加上一个UUID（Universally Unique Identifier）。

== ER 图

#image("images/ER图.png")

= 系统接口设计

== 统一说明

系统接口主要用于前后端交互，本项目前端使用axios作为HTTP客户端，与后端进行异步交互，而后端则提供一系列API和Service，并绑定到相应路由上。本系统的API设计一定程度上遵循RESTful规范，具体设计如下。

注意：由于目前系统仍处于开发阶段，API接口将在编码过程中进行调整，最终的API接口与本文档可能存在差异。

对于入参，本项目在参数较少时选择使用HTTP的query参数，在参数较多时选择在body中使用JSON格式传输DTO（数据传输对象）；对于出参，本项目规定统一的返回格式如下所示：

```json
{
  "code": 200, // HTTP返回码
  "message": "", // 返回的消息
  "data": [] // 返回的数据，可以是任意类型
}
```

其中code是HTTP返回码，与实际返回码保持一致；message是后端返回的提示信息，这两部分的作用比较固定；data是后端返回的具体数据，以下接口返回参数说明中只解释data部分。

== 用户信息相关接口

用户信息类接口主要使用User表的数据，并对User表进行增删改查操作，提供的URL均以user开头，具体的API主要有：

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
  [id是用户ID，username是用户名，email是用户邮箱，isAdmin表示用户是否是管理员，createdTime和updatedTime分别是用户的创建和修改时间],
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
  [用于用户登录，验证用户输入的密码是否正确，若正确则返回JWT和用户信息，这两个东西将存储在浏览器localStorage中以便后续使用]
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
  [用于用户登录，验证用户输入的密码是否正确，若正确则返回JWT和用户信息，这两个东西将存储在浏览器localStorage中以便后续使用]
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
  [id是用户id，duration是持续时间（按分钟计算）],
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

=== 解封用户

#table(
  columns: (150pt, 300pt),
  align: (center, center),
  [接口URL],
  [/user/unsuspend],
  [HTTP请求类型],
  [POST],
  [请求参数类型],
  [query],
  [请求参数示例],
  [id:integer],
  [请求参数说明],
  [id是待解封用户id],
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
  [用于提前解封一个被封禁了的账户，仅限管理员使用]
)

== 图片类接口

图片类接口主要使用Image和Tag表的数据，并对这两张表进行增删改查操作，提供的URL均以image开头，具体的API主要有：

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
  imageName:string],
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
  [用于上传图片，后端会将图片通过文件系统保存到对应目录，图片EXIF信息和其他元信息会被保存到数据库中]
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
  [imageID是图片ID，tags是想要添加的标签（字符串数组）],
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
  [用于为指定图片设置标签，采用直接覆盖方式，所以增删标签均使用这个接口]
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
  [image是新的图片的二进制数据，imageID是图片ID],
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
  [用于更新图片，主要用于用户在前端编辑图片后的更新操作]
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
  userID:integer],
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
  [imageID:integer],
  [请求参数说明],
  [图片ID],
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
  [图片的详细信息，包括标签和EXIF信息等],
  [接口简介],
  [用于根据图片ID查询图片的详细信息]
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
  [imageID:integer],
  [请求参数说明],
  [图片ID],
  [返回值示例],
  [string\<binary\>],
  [返回值说明],
  [真正的图片文件],
  [接口简介],
  [用于获取图片文件，主要是下载图片时需要]
)

= 界面原型

以下是项目的界面原型，部分界面可能还比较简单，主要展示核心交互逻辑，后续会进行美化优化。

== 注册界面

在这个界面用户可以填写自己的基本信息，包括用户名、邮箱和密码，然后注册自己的账户。如果注册成功（用户名和邮箱都不重复，并且密码符合要求），系统会自动跳转到登录页面。

#image("images/界面原型/注册界面.png")

== 登录界面

在这个界面用户可以输入自己的用户名和密码，然后点击登录按钮。如果登录成功（密码正确，且账号未被封禁），则跳转到网站首页。如果其他页面产生了HTTP状态码为401 Unauthorized的请求（说明未登录或者账号被封禁），则会被重定向到登录页面。

#image("images/界面原型/登录界面.png")

== 网站首页

在这个界面会显示一些系统信息和用户的个人信息，可以点击上边栏跳转到其他页面。上边栏中的"用户管理"只会在管理员的前端渲染，并且在通过上边栏跳转时也会进行权限验证。

#image("images/界面原型/网站首页.png")

== 上传图片界面

这个页面用于上传图片，支持批量上传（一次最多10张）。点击加号方框即可从本地挑选图片上传，已经上传到前端的图片会在前端显示预览图，点击提交上传按钮才会正式将图片发送到后端存储。

#image("images/界面原型/上传图片.png")

== 查看图片界面

这个界面的功能比较丰富，下面分别介绍各个功能：

可以在这个界面按照不同条件查询图片，支持的查询条件包括：

- 个人所有图片
- 图片标签
- 拍摄时间区间
- 省市区地理位置

#image("images/界面原型/查询条件.png")

查询结果如下图所示，图片会依次排列并显示最核心的信息，即图片ID、图片名和标签。

#image("images/界面原型/查询结果.png")
点击图片可以进行图片轮播，并且会自动切换到被点击的图片。

#image("images/界面原型/图片轮播.png")

当鼠标悬停在图片上时，图片右上角会浮现四个图标，分别对应查看图片详情、添加标签、编辑图片、下载图片和删除图片。

#image("images/界面原型/图片功能.png")

点击查看图片详情图标后，会弹出对话框显示图片原图和详细信息（主要是图片的EXIF信息）。

#image("images/界面原型/图片详情.png")

点击添加标签图标后，会弹出对话框显示当前图片的标签，可以添加新的标签，也可以删除已有的标签。

#image("images/界面原型/添加标签.png")

点击编辑图片图标后，会弹出对话框显示当前图片原图，并且可以裁剪图片，调整图片的亮度、对比度、饱和度和色相等属性。

#image("images/界面原型/编辑图片.png")

点击下载图片图标后，图片原始文件会被发送到前端，然后自动在浏览器中下载。点击删除图片图标后，如果确认删除，则图片文件和元信息将分别从文件系统和数据库中删除。

= 系统出错设计

== 出错信息

系统通过统一的错误码和友好的信息反馈异常情况，关键错误方式及处理方式如下：

#table(
  columns: (100pt, 120pt, 120pt, 120pt),
  align: (center, center, center, center),
  table.header("错误类型", "错误信息提示", "可能原因", "处理方法"),
  [后端连接失败],
  [弹窗提示“系统繁忙，请稍后再试”],
  [1. 后端启动失败
  2. ip配置错误],
  [1. 检查配置文件
  2. 重启服务],
  [权限不足],
  [页面提示“您无权限执行此操作”],
  [普通用户视图越权访问管理员功能],
  [记录日志并跳转至默认页面],
  [未登录直接访问业务页面],
  [弹窗提示“您尚未登录”],
  [直接通过URL访问业务页面],
  [自动返回登录页面],
  [已在业务页面操作但JWT无效],
  [弹窗提示“您已退出登录”],
  [长时间未操作导致JWT过期],
  [自动返回登录页面],
  [访问不存在的页面],
  [弹窗提示“您访问的页面不存在”],
  [URL输入了错误地址],
  [跳转到404页面],
  [文件格式错误],
  [提示“图片格式仅支持png和jpg”],
  [上传其它格式的图片或信息缺损的图片],
  [强制重新选择图片]
)

== 补救措施

本项目主要通过系统日志来监控项目状态：

1. 在日志中记录所有HTTP请求的信息，包含发起操作者IP、时间戳及详细参数
2. 日志存储在后端文件系统中，保留周期为2个月

= 附录

== 项目进度安排

项目的进度安排如下，目前正在稳步推进中：

#table(
  columns: (150pt, 250pt),
  align: (center, center),
  table.header("时间", "任务"),
  [2025年9月], [项目初期分析和调研],
  [2025年10月], [完成设计报告],
  [2025年11月], [基本完成项目编码],
  [2025年12月], [调试和完善项目，并完成测试报告和产品说明书]
)

== 备注

本设计报告为项目初期的评估，后续还会根据实际情况进行调整。