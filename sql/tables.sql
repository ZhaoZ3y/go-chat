-- 用户表
CREATE TABLE `users` (
                         `id` bigint(20) NOT NULL AUTO_INCREMENT,
                         `username` varchar(50) NOT NULL COMMENT '用户名',
                         `password` varchar(255) NOT NULL COMMENT '密码',
                         `email` varchar(100) NOT NULL COMMENT '邮箱',
                         `nickname` varchar(50) DEFAULT '' COMMENT '昵称',
                         `avatar` varchar(255) DEFAULT '' COMMENT '头像',
                         `phone` varchar(20) DEFAULT '' COMMENT '手机号',
                         `status` tinyint(4) DEFAULT '1' COMMENT '状态 1:正常 2:禁用',
                         `create_time` bigint(20) NOT NULL COMMENT '创建时间',
                         `update_time` bigint(20) NOT NULL COMMENT '更新时间',
                         PRIMARY KEY (`id`),
                         UNIQUE KEY `idx_username` (`username`),
                         UNIQUE KEY `idx_email` (`email`),
                         KEY `idx_phone` (`phone`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='用户表';

-- 好友表
CREATE TABLE `friends` (
                           `id` bigint(20) NOT NULL AUTO_INCREMENT,
                           `user_id` bigint(20) NOT NULL COMMENT '用户ID',
                           `friend_id` bigint(20) NOT NULL COMMENT '好友ID',
                           `remark` varchar(50) DEFAULT '' COMMENT '备注名',
                           `status` tinyint(4) DEFAULT '1' COMMENT '状态 1:正常 2:拉黑',
                           `create_time` bigint(20) NOT NULL COMMENT '创建时间',
                           `update_time` bigint(20) NOT NULL COMMENT '更新时间',
                           PRIMARY KEY (`id`),
                           UNIQUE KEY `idx_user_friend` (`user_id`,`friend_id`),
                           KEY `idx_friend_id` (`friend_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='好友表';

-- 好友申请表
CREATE TABLE `friend_requests` (
                                   `id` bigint(20) NOT NULL AUTO_INCREMENT,
                                   `from_user_id` bigint(20) NOT NULL COMMENT '申请人ID',
                                   `to_user_id` bigint(20) NOT NULL COMMENT '被申请人ID',
                                   `message` varchar(200) DEFAULT '' COMMENT '申请消息',
                                   `status` tinyint(4) DEFAULT '1' COMMENT '状态 1:待处理 2:已同意 3:已拒绝',
                                   `create_time` bigint(20) NOT NULL COMMENT '创建时间',
                                   `update_time` bigint(20) NOT NULL COMMENT '更新时间',
                                   PRIMARY KEY (`id`),
                                   KEY `idx_from_user` (`from_user_id`),
                                   KEY `idx_to_user` (`to_user_id`),
                                   KEY `idx_status` (`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='好友申请表';

-- 群组表
CREATE TABLE `groups` (
                          `id` bigint(20) NOT NULL AUTO_INCREMENT,
                          `name` varchar(100) NOT NULL COMMENT '群组名称',
                          `description` varchar(500) DEFAULT '' COMMENT '群组描述',
                          `avatar` varchar(255) DEFAULT '' COMMENT '群组头像',
                          `owner_id` bigint(20) NOT NULL COMMENT '群主ID',
                          `member_count` int(11) DEFAULT '1' COMMENT '成员数量',
                          `max_member_count` int(11) DEFAULT '500' COMMENT '最大成员数量',
                          `status` tinyint(4) DEFAULT '1' COMMENT '状态 1:正常 2:禁用',
                          `create_time` bigint(20) NOT NULL COMMENT '创建时间',
                          `update_time` bigint(20) NOT NULL COMMENT '更新时间',
                          PRIMARY KEY (`id`),
                          KEY `idx_owner` (`owner_id`),
                          KEY `idx_status` (`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='群组表';

-- 群组成员表
CREATE TABLE `group_members` (
                                 `id` bigint(20) NOT NULL AUTO_INCREMENT,
                                 `group_id` bigint(20) NOT NULL COMMENT '群组ID',
                                 `user_id` bigint(20) NOT NULL COMMENT '用户ID',
                                 `role` tinyint(4) DEFAULT '3' COMMENT '角色 1:群主 2:管理员 3:普通成员',
                                 `nickname` varchar(50) DEFAULT '' COMMENT '群昵称',
                                 `status` tinyint(4) DEFAULT '1' COMMENT '状态 1:正常 2:禁言',
                                 `join_time` bigint(20) NOT NULL COMMENT '加入时间',
                                 `update_time` bigint(20) NOT NULL COMMENT '更新时间',
                                 PRIMARY KEY (`id`),
                                 UNIQUE KEY `idx_group_user` (`group_id`,`user_id`),
                                 KEY `idx_user_id` (`user_id`),
                                 KEY `idx_role` (`role`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='群组成员表';

-- 消息表
CREATE TABLE `messages` (
                            `id` bigint(20) NOT NULL AUTO_INCREMENT,
                            `from_user_id` bigint(20) NOT NULL COMMENT '发送者ID',
                            `to_user_id` bigint(20) DEFAULT '0' COMMENT '接收者ID(私聊)',
                            `group_id` bigint(20) DEFAULT '0' COMMENT '群组ID(群聊)',
                            `type` tinyint(4) NOT NULL DEFAULT '0' COMMENT '消息类型 0:文本 1:图片 2:文件 3:音频 4:视频 5:系统',
                            `content` text NOT NULL COMMENT '消息内容',
                            `extra` text COMMENT '额外信息JSON',
                            `chat_type` tinyint(4) NOT NULL DEFAULT '0' COMMENT '聊天类型 0:私聊 1:群聊',
                            `status` tinyint(4) DEFAULT '1' COMMENT '状态 1:正常 2:已删除 3:已撤回',
                            `create_time` bigint(20) NOT NULL COMMENT '创建时间',
                            PRIMARY KEY (`id`),
                            KEY `idx_from_user` (`from_user_id`),
                            KEY `idx_to_user` (`to_user_id`),
                            KEY `idx_group` (`group_id`),
                            KEY `idx_chat_type` (`chat_type`),
                            KEY `idx_create_time` (`create_time`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='消息表';

-- 会话表
CREATE TABLE `conversations` (
                                 `id` bigint(20) NOT NULL AUTO_INCREMENT,
                                 `user_id` bigint(20) NOT NULL COMMENT '用户ID',
                                 `target_id` bigint(20) NOT NULL COMMENT '目标ID(用户ID或群组ID)',
                                 `type` tinyint(4) NOT NULL DEFAULT '0' COMMENT '类型 0:私聊 1:群聊',
                                 `last_message` text COMMENT '最后一条消息',
                                 `last_message_time` bigint(20) DEFAULT '0' COMMENT '最后消息时间',
                                 `unread_count` int(11) DEFAULT '0' COMMENT '未读消息数',
                                 `create_time` bigint(20) NOT NULL COMMENT '创建时间',
                                 `update_time` bigint(20) NOT NULL COMMENT '更新时间',
                                 PRIMARY KEY (`id`),
                                 UNIQUE KEY `idx_user_target_type` (`user_id`,`target_id`,`type`),
                                 KEY `idx_update_time` (`update_time`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='会话表';

-- 文件表
CREATE TABLE `files` (
                         `id` bigint(20) NOT NULL AUTO_INCREMENT,
                         `filename` varchar(255) NOT NULL COMMENT '文件名',
                         `original_name` varchar(255) NOT NULL COMMENT '原始文件名',
                         `file_path` varchar(500) NOT NULL COMMENT '文件路径',
                         `file_url` varchar(500) NOT NULL COMMENT '文件URL',
                         `file_type` varchar(50) NOT NULL COMMENT '文件类型',
                         `file_size` bigint(20) NOT NULL COMMENT '文件大小',
                         `mime_type` varchar(100) NOT NULL COMMENT 'MIME类型',
                         `hash` varchar(64) NOT NULL COMMENT '文件哈希',
                         `user_id` bigint(20) NOT NULL COMMENT '上传用户ID',
                         `status` tinyint(4) DEFAULT '1' COMMENT '状态 1:正常 2:已删除',
                         `create_time` bigint(20) NOT NULL COMMENT '创建时间',
                         `update_time` bigint(20) NOT NULL COMMENT '更新时间',
                         PRIMARY KEY (`id`),
                         KEY `idx_user` (`user_id`),
                         KEY `idx_hash` (`hash`),
                         KEY `idx_type` (`file_type`),
                         KEY `idx_create_time` (`create_time`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='文件表';

-- 文件分块上传表
CREATE TABLE `file_uploads` (
                                `id` bigint(20) NOT NULL AUTO_INCREMENT,
                                `upload_id` varchar(64) NOT NULL COMMENT '上传ID',
                                `filename` varchar(255) NOT NULL COMMENT '文件名',
                                `file_size` bigint(20) NOT NULL COMMENT '文件大小',
                                `chunk_size` int(11) NOT NULL COMMENT '分块大小',
                                `total_chunks` int(11) NOT NULL COMMENT '总分块数',
                                `uploaded_chunks` text COMMENT '已上传分块列表',
                                `user_id` bigint(20) NOT NULL COMMENT '用户ID',
                                `status` tinyint(4) DEFAULT '1' COMMENT '状态 1:上传中 2:已完成 3:已取消',
                                `create_time` bigint(20) NOT NULL COMMENT '创建时间',
                                `update_time` bigint(20) NOT NULL COMMENT '更新时间',
                                PRIMARY KEY (`id`),
                                UNIQUE KEY `idx_upload_id` (`upload_id`),
                                KEY `idx_user` (`user_id`),
                                KEY `idx_status` (`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='文件分块上传表';