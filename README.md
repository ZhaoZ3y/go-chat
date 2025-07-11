# Go-Chat 项目文档

## 一、项目概述
Go-Chat 是一个基于 Go 语言开发的即时通讯项目，集成了多种服务，如文件管理、用户消息处理、好友关系管理、群组管理等。项目采用微服务架构，使用 gRPC 进行服务间通信，借助 Docker 进行容器化部署，同时使用了 MySQL 数据库、Redis 缓存、Kafka 消息队列和 MinIO 对象存储等技术。

## 二、项目结构
项目主要由以下几个部分组成：
1. **api**：包含 API 接口的控制器，负责处理客户端的请求。
2. **rpc**：包含多个 RPC 服务，如 `file`、`user`、`friend`、`message` 和 `group`，每个服务负责特定的业务逻辑。
3. **pkg**：包含公共的模型、工具函数和配置文件。
4. **cmd**：包含项目的入口文件。

## 三、服务说明

### 1. 文件服务（file.rpc）
- **功能**：处理文件的上传、下载、信息查询和删除等操作，同时支持用户头像的上传和更新。
- **配置文件**：`go-chat/rpc/file/etc/file.yaml`
- **主要逻辑文件**：
  - `go-chat/rpc/file/internal/logic/uploadfilelogic.go`：文件上传逻辑
  - `go-chat/rpc/file/internal/logic/getfilerecordlogic.go`：获取用户文件记录逻辑
  - `go-chat/rpc/file/internal/logic/uploadavatarlogic.go`：用户头像上传逻辑

### 2. 消息服务（message.rpc）
- **功能**：处理消息的标记已读操作，并发送 Kafka 事件通知对方。
- **配置文件**：`go-chat/rpc/message/etc/message.yaml`
- **主要逻辑文件**：`go-chat/rpc/message/internal/logic/markmessagereadlogic.go`

### 3. 群组服务（group.rpc）
- **功能**：获取群组的未读消息总数。
- **主要逻辑文件**：`go-chat/rpc/group/internal/logic/getunreadcountlogic.go`

### 4. 好友服务（friend.rpc）
- **功能**：更新好友备注信息。
- **主要逻辑文件**：`go-chat/rpc/friend/internal/logic/updatefriendremarklogic.go`

## 四、开发环境搭建

### 1. 安装依赖
确保已经安装了 Go 语言环境（版本 1.23.9），并在项目根目录下运行以下命令下载依赖：
```sh
go mod tidy
```

### 2. 生成 RPC 代码
在项目根目录下运行 `Makefile` 中的 `run` 命令，该命令会为每个 RPC 服务生成相应的代码，并启动 Docker 容器：
```sh
make run
```

### 3. 配置文件
根据实际情况修改各个服务的配置文件，如 `go-chat/rpc/file/etc/file.yaml` 和 `go-chat/rpc/message/etc/message.yaml`。

### 4. 启动服务
在项目根目录下运行以下命令启动 API 服务：
```sh
go run cmd/main.go -f pkg/config/config.yaml
```

## 五、API 接口说明

### 1. 文件上传
- **URL**：`/upload`
- **方法**：`POST`
- **请求参数**：
  - `file`：上传的文件
- **响应**：文件上传成功后的信息

### 2. 标记消息已读
- **URL**：`/mark-message-read`
- **方法**：`POST`
- **请求参数**：
  - `target_id`：目标 ID（私聊是对方用户 ID，群聊是群 ID）
  - `chat_type`：聊天类型，0 表示私聊，1 表示群聊
  - `last_read_message_id`：用户已读的最后一条消息的 ID
- **响应**：标记结果信息

### 3. 获取群组未读消息总数
- **URL**：`/group-unread-count`
- **方法**：`GET`
- **请求参数**：
  - `user_id`：用户 ID
- **响应**：群组未读消息总数

### 4. 更新好友备注
- **URL**：`/update-friend-remark`
- **方法**：`POST`
- **请求参数**：
  - `user_id`：用户 ID
  - `friend_id`：好友 ID
  - `remark`：新的备注信息
- **响应**：更新结果信息

## 六、注意事项
- 确保 Docker 服务已经启动，并且相关的容器（如 MySQL、Redis、Kafka、MinIO 等）正常运行。
- 在开发过程中，如需要修改 `proto` 文件，需要重新运行 `Makefile` 中的 `run` 命令生成新的代码。

## 七、贡献指南
如果你想为该项目做出贡献，请遵循以下步骤：
1. Fork 该项目到你的 GitHub 账户。
2. 创建一个新的分支进行开发。
3. 提交代码并发起 Pull Request。

## 八、许可证
该项目采用 [MIT 许可证](LICENSE)。
