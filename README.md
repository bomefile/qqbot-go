


CREATE TABLE `counter_model` (
  `id` int NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `count` int NOT NULL DEFAULT '0' COMMENT '计数器的值',
  `created_at` datetime NOT NULL COMMENT '创建时间',
  `updated_at` datetime NOT NULL COMMENT '更新时间',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='计数器模型表';


DROP TABLE IF EXISTS `user_show_view`;
CREATE TABLE `user_show_view_0` (
`id` bigint NOT NULL AUTO_INCREMENT COMMENT '主键ID，UUID',
`video_content_id` varchar(50) NOT NULL COMMENT '业务id视频内容id， 如： clip_id，moment_id',
`video_content_type` tinyint NOT NULL COMMENT '1 clip，2 moment，3 show，4 eposide，5 movie ,6 series',
`video_id` varchar(50) NOT NULL DEFAULT '' COMMENT '视频id',
`uid` bigint DEFAULT NULL COMMENT '用户id',
`view_time` bigint DEFAULT '0' COMMENT '观看时长 单位秒',
`view_status` tinyint DEFAULT NULL COMMENT '完成状态 0 位置  1 完成 2 未完成',
`ext` varchar(100) NOT NULL DEFAULT '' COMMENT '扩展信息',
`created_at` bigint NOT NULL DEFAULT '0' COMMENT '创建时间（毫秒时间戳）',
`updated_at` bigint NOT NULL DEFAULT '0' COMMENT '更新时间（毫秒时间戳）',
PRIMARY KEY (`id`),
KEY `idx_uid_view_time` (`uid`,`view_status`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb3 COMMENT='俱乐部信息表';


{
acme_ca https://acme-v02.api.letsencrypt.org/directory
email cn.wangliangliang@gmail.com
auto_https off
}

https://8.140.17.9 {
encode gzip
log {
output stdout
format console
}
reverse_proxy qqbot:8080
}