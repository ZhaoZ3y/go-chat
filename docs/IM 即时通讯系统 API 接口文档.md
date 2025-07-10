# IM 即时通讯系统 API 接口文档

## 基本信息

- **基础URL**: `127.0.0.1:8080`
- **认证方式**: JWT Token (Bearer Token)
- **请求格式**: JSON
- **响应格式**: JSON

## 通用成功响应格式

```json
{
  "code": 20000,
  "message": "success",
  "data": {}
}
```

## 通用失败响应格式

```json
{
  "code": XX,
  "message": "yyyy",
}
```

### 响应码说明

| 状态码 | 说明          |
| ------ | ------------- |
| 20000  | 成功          |
| 40000  | 参数错误      |
| 40001  | 未授权        |
| 50000  | 服务器错误    |
| 50001  | RPC客户端错误 |

## 1. 用户认证相关

### 1.1 用户注册

**接口地址**: `POST /register`

**请求参数**:

```json
{
    "username":"admin",
    "password":"123456",
    "email":"admin@admin.com",
    "nickname":"admin",
    "avatar":"https://admin.com"
}
```

**参数说明**:

| 字段     | 类型   | 必填 | 说明   | 限制         |
| -------- | ------ | ---- | ------ | ------------ |
| username | string | 是   | 用户名 | 3-20个字符   |
| password | string | 是   | 密码   | 6-50个字符   |
| email    | string | 是   | 邮箱   | 有效邮箱格式 |
| nickname | string | 否   | 昵称   | -            |
| avatar   | string | 否   | 头像   | 有效url格式  |

**响应示例**:

```json
{
  "code": 20000,
  "message": "success",
  "data": {
    "message": "注册成功"
  }
}
```

### 1.2 用户登录

**接口地址**: `POST /login`

**请求参数**:

```json
{
  "username": "testuser",
  "password": "123456"
}
```

**参数说明**:

| 字段     | 类型   | 必填 | 说明   |
| -------- | ------ | ---- | ------ |
| username | string | 是   | 用户名 |
| password | string | 是   | 密码   |

**响应示例**:

```json
{
  "code": 20000,
  "message": "success",
  "data": {
    "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
  }
}
```

### 1.3 刷新Token

**接口地址**: `POST /refresh_token`

**请求参数**:

```json
{
  "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

**响应示例**:

```json
{
  "code": 20000,
  "message": "success",
  "data": {
    "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
  }
}
```

## 2. 用户管理

### 2.1 获取用户个人资料

**接口地址**: `GET /user/profile`

**请求头**: `Authorization: Bearer <token>`

**响应示例**:

```json
{
    "code": 20000,
    "message": "success",
    "data": {
        "id": 11,
        "username": "admin",
        "email": "admin@admin.com",
        "nickname": "admin",
        "avatar": "http://demo.com",
        "status": 1,
        "create_at": 1752122703,
        "update_at": 1752122952
    }
}
```

### 2.2 获取用户信息

**接口地址**: `GET /user/info`

**请求头**: `Authorization: Bearer <token>`

**请求参数**(query):

| 字段    | 类型 | 必填 | 说明   |
| ------- | ---- | ---- | ------ |
| user_id | int  | 是   | 用户ID |

**响应示例**:

```json
{
    "code": 20000,
    "message": "success",
    "data": {
        "id": 11,
        "username": "admin",
        "email": "admin@admin.com",
        "nickname": "admin",
        "avatar": "http://demo.com",
        "status": 1,
        "create_at": 1752122703,
        "update_at": 1752122952
    }
}
```

### 2.3 更新用户信息

**接口地址**: `PUT /user/update/info`

**请求头**: `Authorization: Bearer <token>`

**请求参数**:

```json
{
  "nickname": "新昵称",
  "avatar": "https://example.com/new-avatar.jpg",
  "email": "new@example.com"
}
```

**参数说明**:

| 字段     | 类型   | 必填 | 说明    |
| -------- | ------ | ---- | ------- |
| nickname | string | 否   | 昵称    |
| avatar   | string | 否   | 头像URL |
| email    | string | 否   | 邮箱    |

**响应示例**:

```json
{
  "code": 20000,
  "message": "success",
  "data": {
    "message": "更新成功"
  }
}
```

### 2.4 修改密码

**接口地址**: `PUT /user/update/password`

**请求头**: `Authorization: Bearer <token>`

**请求参数**:

```json
{
  "old_password": "123456",
  "new_password": "654321"
}
```

**参数说明**:

| 字段         | 类型   | 必填 | 说明   | 限制       |
| ------------ | ------ | ---- | ------ | ---------- |
| old_password | string | 是   | 旧密码 | -          |
| new_password | string | 是   | 新密码 | 6-50个字符 |

**响应示例**:

```json
{
  "code": 20000,
  "message": "success",
  "data": {
    "message": "密码修改成功"
  }
}
```

## 3. 搜索功能

### 3.1 搜索用户和群组

**接口地址**: `GET /search`

**请求头**: `Authorization: Bearer <token>`

**请求参数**:

| 字段    | 类型   | 必填 | 说明       |
| ------- | ------ | ---- | ---------- |
| keyword | string | 是   | 搜索关键词 |

**响应示例**:

```json
{
    "code": 20000,
    "message": "success",
    "data": {
        "groups": [
            {
                "id": 3,
                "name": "group3",
                "description": "2222",
                "avatar": "demo.com",
                "owner_id": 6,
                "member_count": 1,
                "max_member_count": 500,
                "status": 1,
                "create_at": 1751612966,
                "update_at": 1751612966
            },
            {
                "id": 2,
                "name": "group2",
                "description": "2222",
                "avatar": "demo.com",
                "owner_id": 6,
                "member_count": 1,
                "max_member_count": 500,
                "status": 1,
                "create_at": 1751610984,
                "update_at": 1751610984
            }
        ],
        "totalGroups": 2,
        "totalUsers": 2,
        "users": [
            {
                "id": 5,
                "username": "111111",
                "email": "demo@2.com",
                "nickname": "user1",
                "avatar": "http://demo.com",
                "status": 1,
                "create_at": 1748670274,
                "update_at": 1748672791
            },
            {
                "id": 1,
                "username": "123456",
                "email": "demo@1.com",
                "nickname": "user1",
                "status": 1,
                "create_at": 1748669232,
                "update_at": 1748669232
            }
        ]
    }
}
```

## 4. 好友管理

### 4.1 发送好友申请

**接口地址**: `POST /friend/request/send`

**请求头**: `Authorization: Bearer <token>`

**请求参数**:

```json
{
  "to_user_id": 2,
  "message": "我想加您为好友"
}
```

### 4.2 处理好友申请

**接口地址**: `POST /friend/request/handle`

**请求头**: `Authorization: Bearer <token>`

**请求参数**:

```json
{
  "request_id": 1,
  "action": 
}
```

**参数说明**:

| 字段       | 类型   | 必填 | 说明           |
| ---------- | ------ | ---- | -------------- |
| request_id | int    | 是   | 申请ID         |
| action     | int    | 是   | 2是同意3是拒绝 |
| message    | string | 否   | 原因           |

### 4.3 获取好友列表

**接口地址**: `GET /friend/list`

**请求头**: `Authorization: Bearer <token>`

**响应示例**:

```json
{
    "code": 20000,
    "message": "success",
    "data": {
        "friends": [
            {
                "id": 6,
                "user_id": 6,
                "friend_id": 8,
                "status": 1,
                "create_at": 1751277401,
                "update_at": 1751277401,
                "nickname": "user1",
                "avatar":"http://demo.com",
                "online_status":1
            }
        ],
        "total": 1
    }
}
```

### 4.4 删除好友

**接口地址**: `DELETE /friend/delete`

**请求头**: `Authorization: Bearer <token>`

**请求参数**:

```json
{
  "friend_id": 2
}
```

### 4.5 拉黑好友

**接口地址**: `POST /friend/block`

**请求头**: `Authorization: Bearer <token>`

**请求参数**:

```json
{
    "friend_id":5
}
```

### 4.6 获取拉黑好友列表

**接口地址**: `GET /friend/blocked/list`

**请求头**: `Authorization: Bearer <token>`

**响应示例**：

```json
{
    "code": 20000,
    "message": "success",
    "data": {
        "friends": [
            {
                "id": 5,
                "user_id": 8,
                "friend_id": 6,
                "remark": "11222",
                "status": 2,
                "create_at": 1751277401,
                "update_at": 1752125061
            }
        ],
        "total": 1
    }
}
```

### 4.7 更新好友备注

**接口地址**: `PUT /friend/remark`

**请求头**: `Authorization: Bearer <token>`

**请求参数**:

```json
{
  "user_id": 2,
  "remark": "新的备注"
}
```

## 5. 群组管理

### 5.1 创建群组

**接口地址**: `POST /group/create`

**请求头**: `Authorization: Bearer <token>`

**请求参数**:

```json
{
  "group_name": "测试群组",
  "description": "这是一个测试群组",
  "avatar": "https://example.com/group-avatar.jpg"
}
```

### 5.2 获取群组信息

**接口地址**: `GET /group/info`

**请求头**: `Authorization: Bearer <token>`

**请求参数**(query):

| 字段     | 类型 | 必填 | 说明   |
| -------- | ---- | ---- | ------ |
| group_id | int  | 是   | 群组ID |

**响应示例**

```json
{
    "code": 20000,
    "message": "success",
    "data": {
        "group_info": {
            "id": 1,
            "name": "group1",
            "description": "11111",
            "owner_id": 6,
            "member_count": 1,
            "max_member_count": 500,
            "status": 1,
            "create_at": 1751610303,
            "update_at": 1751612940
        }
    }
}
```

### 5.3 获取用户的群组列表

**接口地址**: `GET /group/list`

**请求头**: `Authorization: Bearer <token>`

**响应示例**

```json
{
    "code": 20000,
    "message": "success",
    "data": {
        "groups": [
            {
                "id": 1,
                "name": "group1",
                "description": "11111",
                "owner_id": 6,
                "member_count": 1,
                "max_member_count": 500,
                "status": 1,
                "create_at": 1751610303,
                "update_at": 1751612940
            },
            {
                "id": 2,
                "name": "group2",
                "description": "2222",
                "avatar": "demo.com",
                "owner_id": 6,
                "member_count": 1,
                "max_member_count": 500,
                "status": 1,
                "create_at": 1751610984,
                "update_at": 1751610984
            },
            {
                "id": 3,
                "name": "group3",
                "description": "2222",
                "avatar": "demo.com",
                "owner_id": 6,
                "member_count": 1,
                "max_member_count": 500,
                "status": 1,
                "create_at": 1751612966,
                "update_at": 1751612966
            },
            {
                "id": 4,
                "name": "group3",
                "description": "2222",
                "avatar": "demo.com",
                "owner_id": 6,
                "member_count": 1,
                "max_member_count": 500,
                "create_at": 1751613178,
                "update_at": 1751613178
            }
        ],
        "total": 4
    }
}
```

### 5. 获取群组成员列表

**接口地址**: `GET /group/members`

**请求头**: `Authorization: Bearer <token>`

**请求参数**(query):

| 字段     | 类型 | 必填 | 说明   |
| -------- | ---- | ---- | ------ |
| group_id | int  | 是   | 群组ID |

**响应示例**

```json
{
    "code": 20000,
    "message": "success",
    "data": {
        "members": [
            {
                "id": 1,
                "group_id": 1,
                "user_id": 6,
                "role": 2,
                "nickname": "user1",
                "join_time": 1751610303,
                "update_at": 1751782784
            }
        ],
        "total": 1
    }
}
```

### 5.5 更新群组信息

**接口地址**: `PUT /group/update/info`

**请求头**: `Authorization: Bearer <token>`

**请求参数**:

```json
{
  "group_id": 1,
  "name": "新的群组名称",
  "description": "新的群组描述",
  "avatar": "https://example.com/new-group-avatar.jpg"
}
```

### 5.6 设置管理员

**接口地址**: `POST /group/set/role`

**请求头**: `Authorization: Bearer <token>`

**请求参数**:

```json
{
  "group_id": 1,
  "user_id": 2,
  "role": 1
}
```

### 5.7 禁言群成员

**接口地址**: `POST /group/mute/member`

**请求头**: `Authorization: Bearer <token>`

**请求参数**:

```json
{
  "group_id": 1,
  "user_id": 2,
  "duration": 3600
}
```

### 5.8 邀请用户加入群组

**接口地址**: `POST /group/invite`

**请求头**: `Authorization: Bearer <token>`

**请求参数**:

```json
{
  "group_id": 1,
  "user_ids": [2, 3, 4]
}
```

### 5.9 申请加入群组

**接口地址**: `POST /group/join`

**请求头**: `Authorization: Bearer <token>`

**请求参数**:

```json
{
  "group_id": 1,
  "reason": "我想加入这个群组"
}
```

### 5.10 退出群组

**接口地址**: `POST /group/leave`

**请求头**: `Authorization: Bearer <token>`

**请求参数**:

```json
{
  "group_id": 1
}
```

### 5.11 解散群组

**接口地址**: `POST /group/dismiss`

**请求头**: `Authorization: Bearer <token>`

**请求参数**:

```json
{
  "group_id": 1
}
```

### 5.12 转让群组

**接口地址**: `POST /group/transfer`

**请求头**: `Authorization: Bearer <token>`

**请求参数**:

```json
{
  "group_id": 1,
  "new_owner_id": 2
}
```

### 5.13 踢出群成员

**接口地址**: `POST /group/kick`

**请求头**: `Authorization: Bearer <token>`

**请求参数**:

```json
{
  "group_id": 1,
  "user_id": 2
}
```

### 5.14 获取群成员信息

**接口地址**: `GET /group/member_info`

**请求头**: `Authorization: Bearer <token>`

**请求参数**（query）:

| 字段     | 类型 | 必填 | 说明   |
| -------- | ---- | ---- | ------ |
| group_id | int  | 是   | 群组ID |
| user_id  | int  | 是   | 用户ID |

**响应示例**

```json
{
    "code": 20000,
    "message": "success",
    "data": {
        "member_info": {
            "member_info": {
                "group_id": 1,
                "user_id": 6,
                "role": 2,
                "nickname": "user1",
                "join_time": 1751610303
            },
            "user_info": {
                "id": 6,
                "username": "222222",
                "nickname": "user2"
            }
        }
    }
}
```

### 5.15 更新自己群成员信息

**接口地址**: `PUT /group/update/member_info`

**请求头**: `Authorization: Bearer <token>`

**请求参数**:

```json
{
  "group_id": 1,
  "nickname": "新的群内昵称"
}
```

### 5.16 处理入群申请

**接口地址**: `POST /group/request/handle`

**请求头**: `Authorization: Bearer <token>`

**请求参数**:

```json
{
    "application_id":3,
    "approve":true
}
```

## 6. 通知管理

### 6.1 获取好友申请列表

**接口地址**: `GET /notification/request/friend`

**请求头**: `Authorization: Bearer <token>`

**响应示例**

```json
{
    "code": 20000,
    "message": "success",
    "data": {
        "requests": [
          {
            "id": 6,
            "from_user_id": 5,
            "to_user_id": 6,
            "message": "你好，我是user1",
            "status": 1,
            "create_at": 1752054949,
            "update_at": 1752054949,
            "from_nickname": "user1",
            "from_avatar": "http://demo.com",
            "to_nickname": "user2",
            "to_avatar": "http://demo.com"
          }
        ],
        "total": 1
    }
}
```

### 6.2 获取未读好友申请数量

**接口地址**: `GET /notification/request/friend/unread-count`

**请求头**: `Authorization: Bearer <token>`

**响应示例**

```json
{
    "code": 20000,
    "message": "success",
    "data": {
        "count": 1
    }
}
```

### 6.3 获取群组相关通知

**接口地址**: `GET /notification/group/notifications`

**请求头**: `Authorization: Bearer <token>`

**响应示例**

```json
{
    "code": 20000,
    "message": "success",
    "data": {
        "notifications": [
          {
            "id": 4,
            "type": 8,
            "group_id": 5,
            "operator_id": 6,
            "target_user_id": 8,
            "message": "您已被 'user1' 设置为群聊 'group4' 的新群主",
            "timestamp": 1751615721,
            "is_read": true,
            "operator_nickname": "user2",
            "operator_avatar": "http://demo.com",
            "target_user_nickname": "user1",
            "target_user_avatar": "http://demo.com",
          }
        ]
    }
}
```



### 6.4 获取入群申请列表

**接口地址**: `GET /notification/request/group`

**请求头**: `Authorization: Bearer <token>`

**响应示例**

```json
{
    "code": 20000,
    "message": "success",
    "data": {
        "applications": [
          {
            "id": 5,
            "group_id": 6,
            "user_id": 11,
            "reason": "12233",
            "apply_time": 1752126707,
            "user_nickname": "admin",
            "user_avatar": "http://demo.com",
            "group_name": "group2",
            "group_avatar": "demo.com"
          }
        ]
    }
}
```

### 6.5 获取未读入群申请数量

**接口地址**: `GET /notification/request/group/unread-count`

**请求头**: `Authorization: Bearer <token>`

**响应示例**

```json
{
    "code": 20000,
    "message": "success",
    "data": {
        "total_unread_count": 0
    }
}
```

## 7. 文件管理

### 7.1 上传文件

**接口地址**: `POST /file/upload`

**请求头**: `Authorization: Bearer <token>`

**请求类型**: `multipart/form-data`

### 7.2 下载文件

**接口地址**: `GET /file/download`

**请求头**: `Authorization: Bearer <token>`

**请求参数**(query):

| 字段    | 类型   | 必填 | 说明   |
| ------- | ------ | ---- | ------ |
| file_id | string | 是   | 文件ID |

### 7.3 删除文件

**接口地址**: `DELETE /file/delete`

**请求头**: `Authorization: Bearer <token>`
**请求参数**(query):

| 字段    | 类型   | 必填 | 说明   |
| ------- | ------ | ---- | ------ |
| file_id | string | 是   | 文件ID |


### 7.4 获取文件信息

**接口地址**: `GET /file/info`

**请求头**: `Authorization: Bearer <token>`

**请求参数**(query):

| 字段    | 类型   | 必填 | 说明   |
| ------- | ------ | ---- | ------ |
| file_id | string | 是   | 文件ID |

**响应示例**

```json
{
    "code": 20000,
    "message": "success",
    "data": {
        "file_id": "bfaed552-3886-4dc0-84bc-c6826a1fb3f9",
        "file_name": "1111.jpg",
        "file_size": 47531,
        "content_type": "image/jpeg",
        "user_id": 8,
        "created_at": 1752127645,
        "expire_at": 1752732445,
        "etag": "18f9827bdf038fa6bbb8a2c411fe1ed1",
        "file_type": "image"
    }
}
```

### 7.5 获取用户文件记录

**接口地址**: `GET /file/record`

**请求头**: `Authorization: Bearer <token>`

**响应示例**

~~~json
{
    "code": 20000,
    "message": "success",
    "data": {
        "file_records": [
            {
                "file_id": "bfaed552-3886-4dc0-84bc-c6826a1fb3f9",
                "file_name": "1111.jpg",
                "file_size": 47531,
                "content_type": "image/jpeg",
                "user_id": 8,
                "created_at": 1752127645,
                "expire_at": 1752732445,
                "etag": "18f9827bdf038fa6bbb8a2c411fe1ed1",
                "file_type": "image"
            }
        ]
    }
}
~~~

## 8. 聊天功能

### 8.1 发送消息

**接口地址**: `POST /chat/send`

**请求头**: `Authorization: Bearer <token>`

**请求参数**:

```json
{
    "to_user_id":5,
    "type":0,
    "content":"214124124",
    "chat_type":0
}(私聊)
{
    "group_id":6,
    "type":0,
    "content":"214124124",
    "chat_type":1
}(群聊)
```

**参数说明**:

| 字段          | 类型   | 必填 | 说明                                   |
| ------------- | ------ | ---- | -------------------------------------- |
| receiver_id   | int    | 是   | 接收者ID                               |
| receiver_type | string | 是   | 接收者类型 (user/group)                |
| content       | string | 是   | 消息内容                               |
| message_type  | string | 是   | 消息类型 (text/image/file/audio/video) |

### 8.2 获取历史消息

**接口地址**: `GET /chat/history`

**请求头**: `Authorization: Bearer <token>`

**请求参数**:

```json
{
    "target_id":5,
    "chat_type":0(0是私聊1是群聊)
}
```

**响应示例**

~~~json
{
    "code": 20000,
    "message": "success",
    "data": {
        "hasMore": false,
        "messages": [
            {
                "id": 7,
                "from_user_id": 6,
                "to_user_id": 5,
                "content": "214124124",
                "create_at": 1752128377
            },
            {
                "id": 8,
                "from_user_id": 6,
                "to_user_id": 5,
                "content": "214124124",
                "create_at": 1752128377
            }
        ]
    }
}
~~~



### 8.3 标记消息已读

**接口地址**: `POST /chat/read`

**请求头**: `Authorization: Bearer <token>`

**请求参数**:

```json
{
    "target_id":5,
    "chat_type":0(0是私聊1是群聊)
}
```

### 8.4 撤回消息

**接口地址**: `POST /chat/message/recall`

**请求头**: `Authorization: Bearer <token>`

**请求参数**:

```json
{
  "message_id": 123
}
```

### 8.5 删除消息

**接口地址**: `POST /chat/message/delete`

**请求头**: `Authorization: Bearer <token>`

**请求参数**:

```json
{
  "message_id": 123
}
```

### 8.6 获取会话列表

**接口地址**: `GET /chat/conversation/list`

**请求头**: `Authorization: Bearer <token>`

```json
{
    "code": 20000,
    "message": "success",
    "data": [
        {
            "id": 15,
            "user_id": 6,
            "target_id": 6,
            "type": 1,
            "last_message_id": 17,
            "last_message": "214124124",
            "last_message_time": 1752130226,
            "create_at": 1752130225,
            "update_at": 1752130225,
            "target_name": "group2",
            "target_avatar": "demo.com"
        },
        {
            "id": 11,
            "user_id": 6,
            "target_id": 5,
            "last_message_id": 8,
            "last_message": "214124124",
            "last_message_time": 1752128377,
            "create_at": 1752128377,
            "update_at": 1752128377,
            "target_name": "user1",
            "target_avatar": "http://demo.com"
        }
    ]
}
```

### 8.7 删除会话

**接口地址**: `POST /chat/conversation/delete`

**请求头**: `Authorization: Bearer <token>`

**请求参数**:

```json
{
  "conversation_id": "conv123"
}
```

### 8.8 置顶/取消置顶会话

**接口地址**: `POST /chat/conversation/pin`

**请求头**: `Authorization: Bearer <token>`

**请求参数**:

```json
{
  "target_id": 1,
  "chat_type": 1,
  "is_pinned": true
}
```

## 9. WebSocket 连接

### 9.1 建立WebSocket连接

**接口地址**: `GET /ws`

**请求头**: `Authorization: Bearer <token>`

**连接说明**:

- 需要在连接时提供有效的JWT Token
- 连接成功后可以实时收发消息
- 支持心跳检测保持连接活跃

**消息格式**:

```json
{
  "type": "message",
  "data": {
    "message_id": 123,
    "sender_id": 1,
    "receiver_id": 2,
    "receiver_type": "user",
    "content": "Hello, World!",
    "message_type": "text",
    "timestamp": "2024-01-01T00:00:00Z"
  }
}
```

## 错误处理

### 常见错误响应

#### 参数错误 (40000)

```json
{
  "code": 40000,
  "message": "参数错误: 用户名不能为空"
}
```

#### 未授权 (40001)

```json
{
  "code": 40001,
  "message": "用户未登录或ID无效"
}
```

#### 服务器错误 (50000)

```json
{
  "code": 50000,
  "message": "服务器内部错误"
}
```

#### RPC客户端错误 (50001)

```json
{
  "code": 50001,
  "message": "RPC服务调用失败"
}
```