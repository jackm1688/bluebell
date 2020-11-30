
#CREATE USER 'bluebell'@'localhost' IDENTIFIED BY 'bluebell@2020';
# Grant all privileges on bluebell.* to 'bluebell'@'localhost';
DROP DATABASE  IF EXISTS  `bluebell`;
CREATE DATABASE `bluebell` DEFAULT CHARACTER SET utf8mb4;

use `bluebell`;
DROP TABLE IF EXISTS  `user`;
CREATE TABLE `user`(
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `user_id` bigint(20) NOT NULL,
  `username` VARCHAR(64) COLLATE utf8mb4_general_ci NOT NULL,
  `password` VARCHAR(64) COLLATE utf8mb4_general_ci NOT NULL,
  `email` VARCHAR(64) COLLATE utf8mb4_general_ci,
  `gender` tinyint(64) NOT NULL DEFAULT '0',
  `create_time` TIMESTAMP NULL DEFAULT CURRENT_TIMESTAMP,
  `udpate_time` TIMESTAMP NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE current_timestamp ,
  PRIMARY KEY(`id`),
  UNIQUE  key idx_username (`username`) USING BTREE,
  UNIQUE KEY idx_user_id(`user_id`) USING BTREE
) engine=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

#创建社区
DROP TABLE IF EXISTS  `community`;
CREATE TABLE `community`(
    `id` int(11) NOT NULL AUTO_INCREMENT,
    `community_id` int(10) unsigned NOT NULL,
    `community_name` VARCHAR(128) COLLATE utf8mb4_general_ci NOT NULL,
    `introduction` VARCHAR(256) COLLATE utf8mb4_general_ci NOT NULL,
    `create_time` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `update_time` TIMESTAMP NOT NULL DEFAULT  CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`),
    UNIQUE KEY `idx_community_id`(`community_id`),
    UNIQUE KEY `idx_community_name`(`community_name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

INSERT INTO `community` (community_id,community_name,introduction) VALUES(1,"GO","Golang");
INSERT INTO `community` (community_id,community_name,introduction) VALUES(2,"MySQL","Database");
INSERT INTO `community` (community_id,community_name,introduction) VALUES(3,"Java","Java dev");
INSERT INTO `community` (community_id,community_name,introduction) VALUES(4,"Code","code dev");

DROP TABLE IF EXISTS `post`;
CREATE TABLE `post`(
    `id` bigint(20) NOT NULL AUTO_INCREMENT,
    `post_id` bigint(20) NOT NULL COMMENT '帖子ID',
    `title` VARCHAR(128) NOT NULL COLLATE utf8mb4_general_ci COMMENT '标题',
    `content` varchar(8192)  NOT NULL COLLATE utf8mb4_general_ci comment '内容',
    `author_id` bigint(20) NOT NULL COMMENT '作者的用户ID',
    `community_id` bigint(20) NOT NULL COMMENT '社区id',
    `status` tinyint(4) NOT NULL DEFAULT '1' COMMENT '帖子状态',
    `create_time` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `update_time` TIMESTAMP NOT NULL DEFAULT  CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP  COMMENT '更新时间',
    PRIMARY KEY (`id`),
    UNIQUE KEY `idx_post_id`(`post_id`),
    KEY `idx_author_id`(`author_id`),
    KEY `idx_community_id`(`community_id`)
) engine=InnoDB charset=utf8mb4 Collate=utf8mb4_general_ci;