-- 用户表
CREATE TABLE `users` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `username` varchar(50) NOT NULL COMMENT '用户名',
  `password` varchar(100) NOT NULL COMMENT '密码',
  `nickname` varchar(50) NOT NULL COMMENT '昵称',
  `avatar` varchar(255) DEFAULT '' COMMENT '头像',
  `email` varchar(100) DEFAULT '' COMMENT '邮箱',
  `phone` varchar(20) DEFAULT '' COMMENT '手机号',
  `status` tinyint DEFAULT '1' COMMENT '状态：1正常 2禁用',
  `last_login_at` datetime DEFAULT NULL COMMENT '最后登录时间',
  `create_at` datetime DEFAULT CURRENT_TIMESTAMP,
  `update_at` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_username` (`username`),
  UNIQUE KEY `uk_email` (`email`),
  KEY `idx_status` (`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='用户表';

-- 好友关系表
CREATE TABLE `friendships` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `user_id` bigint NOT NULL COMMENT '用户ID',
  `friend_id` bigint NOT NULL COMMENT '好友ID',
  `status` tinyint DEFAULT '1' COMMENT '状态：1正常 2已删除',
  `create_at` datetime DEFAULT CURRENT_TIMESTAMP,
  `update_at` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_user_friend` (`user_id`,`friend_id`),
  KEY `idx_friend_id` (`friend_id`),
  KEY `idx_status` (`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='好友关系表';

-- 好友申请表
CREATE TABLE `friend_requests` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `from_user_id` bigint NOT NULL COMMENT '申请人ID',
  `to_user_id` bigint NOT NULL COMMENT '被申请人ID',
  `message` varchar(255) DEFAULT '' COMMENT '申请消息',
  `status` tinyint DEFAULT '0' COMMENT '状态：0待处理 1已同意 2已拒绝',
  `create_at` datetime DEFAULT CURRENT_TIMESTAMP,
  `update_at` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `idx_from_user` (`from_user_id`),
  KEY `idx_to_user` (`to_user_id`),
  KEY `idx_status` (`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='好友申请表';

-- 群组表
CREATE TABLE `groups` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `name` varchar(100) NOT NULL COMMENT '群名称',
  `description` varchar(500) DEFAULT '' COMMENT '群描述',
  `avatar` varchar(255) DEFAULT '' COMMENT '群头像',
  `owner_id` bigint NOT NULL COMMENT '群主ID',
  `member_count` int DEFAULT '0' COMMENT '成员数量',
  `max_member_count` int DEFAULT '500' COMMENT '最大成员数量',
  `status` tinyint DEFAULT '1' COMMENT '状态：1正常 2解散',
  `create_at` datetime DEFAULT CURRENT_TIMESTAMP,
  `update_at` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `idx_owner` (`owner_id`),
  KEY `idx_status` (`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='群组表';

-- 群成员表
CREATE TABLE `group_members` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `group_id` bigint NOT NULL COMMENT '群组ID',
  `user_id` bigint NOT NULL COMMENT '用户ID',
  `role` tinyint DEFAULT '0' COMMENT '角色：0普通成员 1管理员 2群主',
  `nickname` varchar(50) DEFAULT '' COMMENT '群昵称',
  `mute_until` datetime DEFAULT NULL COMMENT '禁言到期时间',
  `join_at` datetime DEFAULT CURRENT_TIMESTAMP COMMENT '加入时间',
  `status` tinyint DEFAULT '1' COMMENT '状态：1正常 2已退出 3被踢出',
  `create_at` datetime DEFAULT CURRENT_TIMESTAMP,
  `update_at` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_group_user` (`group_id`,`user_id`),
  KEY `idx_user_id` (`user_id`),
  KEY `idx_status` (`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='群成员表';

-- 消息表
CREATE TABLE `messages` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `from_user_id` bigint NOT NULL COMMENT '发送者ID',
  `to_user_id` bigint DEFAULT NULL COMMENT '接收者ID（私聊）',
  `to_group_id` bigint DEFAULT NULL COMMENT '群组ID（群聊）',
  `type` tinyint NOT NULL COMMENT '消息类型：1文本 2图片 3语音 4视频 5文件',
  `content` text COMMENT '消息内容',
  `file_id` bigint DEFAULT NULL COMMENT '文件ID',
  `reply_to` bigint DEFAULT NULL COMMENT '回复的消息ID',
  `status` tinyint DEFAULT '1' COMMENT '状态：1正常 2已撤回 3已删除',
  `create_at` datetime DEFAULT CURRENT_TIMESTAMP,
  `update_at` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `idx_from_user` (`from_user_id`),
  KEY `idx_to_user` (`to_user_id`),
  KEY `idx_to_group` (`to_group_id`),
  KEY `idx_create_at` (`create_at`),
  KEY `idx_reply_to` (`reply_to`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='消息表';

-- 消息读取状态表
CREATE TABLE `message_reads` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `message_id` bigint NOT NULL COMMENT '消息ID',
  `user_id` bigint NOT NULL COMMENT '用户ID',
  `read_at` datetime DEFAULT CURRENT_TIMESTAMP COMMENT '读取时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_message_user` (`message_id`,`user_id`),
  KEY `idx_user_id` (`user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='消息读取状态表';

-- 会话表
CREATE TABLE `conversations` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `user_id` bigint NOT NULL COMMENT '用户ID',
  `target_id` bigint NOT NULL COMMENT '目标ID（用户ID或群组ID）',
  `type` tinyint NOT NULL COMMENT '会话类型：1私聊 2群聊',
  `last_msg_id` bigint DEFAULT NULL COMMENT '最后一条消息ID',
  `unread_count` int DEFAULT '0' COMMENT '未读消息数',
  `is_muted` tinyint DEFAULT '0' COMMENT '是否免打扰',
  `is_pinned` tinyint DEFAULT '0' COMMENT '是否置顶',
  `update_at` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_user_target_type` (`user_id`,`target_id`,`type`),
  KEY `idx_update_at` (`update_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='会话表';

-- 文件表
CREATE TABLE `files` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `name` varchar(255) NOT NULL COMMENT '文件名',
  `original_name` varchar(255) NOT NULL COMMENT '原始文件名',
  `type` varchar(50) NOT NULL COMMENT '文件类型',
  `size` bigint NOT NULL COMMENT '文件大小',
  `path` varchar(500) NOT NULL COMMENT '文件路径',
  `url` varchar(500) NOT NULL COMMENT '文件URL',
  `hash` varchar(64) NOT NULL COMMENT '文件哈希',
  `upload_user_id` bigint NOT NULL COMMENT '上传用户ID',
  `status` tinyint DEFAULT '1' COMMENT '状态：1正常 2已删除',
  `create_at` datetime DEFAULT CURRENT_TIMESTAMP,
  `update_at` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_hash` (`hash`),
  KEY `idx_upload_user` (`upload_user_id`),
  KEY `idx_status` (`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='文件表';

-- 用户在线状态表（Redis缓存）
-- key: online:user:{user_id}
-- value: {server_id, connect_time, last_ping_time}

-- 用户连接会话表（Redis缓存）
-- key: connection:{connection_id}
-- value: {user_id, server_id, connect_time}

-- 消息队列（Redis Stream）
-- stream: message_queue
-- 消息格式: {type, from_user_id, to_user_id, to_group_id, message_id, content}

-- 创建索引优化查询性能
CREATE INDEX idx_messages_conversation ON messages (from_user_id, to_user_id, create_at DESC);
CREATE INDEX idx_messages_group_conversation ON messages (to_group_id, create_at DESC);
CREATE INDEX idx_conversations_user_update ON conversations (user_id, update_at DESC);