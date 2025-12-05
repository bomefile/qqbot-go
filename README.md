

'''mysql
CREATE TABLE `counter_model` (
  `id` int NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `count` int NOT NULL DEFAULT '0' COMMENT '计数器的值',
  `created_at` datetime NOT NULL COMMENT '创建时间',
  `updated_at` datetime NOT NULL COMMENT '更新时间',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='计数器模型表';
'''