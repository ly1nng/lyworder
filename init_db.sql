CREATE DATABASE `ticket_system` DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci;  /*!40100 DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci */ /*!80016 DEFAULT ENCRYPTION='N' */;

-- ticket_system.tickets definition

CREATE TABLE IF NOT EXISTS `tickets` (
  `id` varchar(20) NOT NULL,
  `title` varchar(100)  NOT NULL COMMENT '工单标题',
  `content` text  NOT NULL COMMENT '工单内容',
  `status` varchar(20)  NOT NULL COMMENT '工单状态',
  `user_id` varchar(50)  NOT NULL COMMENT '用户ID',
  `user_name` varchar(100)  NOT NULL COMMENT '用户名称',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` datetime DEFAULT NULL COMMENT '更新时间',
  `operator_name` varchar(255) DEFAULT NULL COMMENT '处理人',
  `remark` text COMMENT '工单备注',
  `solution` text COMMENT '工单解决方案',
  `ticket_type` varchar(255) DEFAULT NULL COMMENT '工单类型',
  `screen_shots` json DEFAULT NULL COMMENT '截图路径数组，JSON格式存储多张图片路径',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;


-- ticket_system.users definition

CREATE TABLE IF NOT EXISTS `users`  (
  `id` int NOT NULL AUTO_INCREMENT,
  `ops_type` varchar(50) DEFAULT NULL COMMENT '运维类型',
  `user_name` varchar(100) NOT NULL COMMENT '运维用户名',
  `role` varchar(50) NOT NULL DEFAULT 'user' COMMENT '运维角色',
  `status` tinyint NOT NULL DEFAULT 1 COMMENT '工作状态，1表示工作中，0表示非工作中',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` datetime DEFAULT NULL COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `user_name` (`user_name`)
) ENGINE=InnoDB AUTO_INCREMENT=4 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

-- 初始化运维账号
INSERT INTO `users` (`user_name`, `role`, `status`) VALUES ('admin', 'admin', 1);



