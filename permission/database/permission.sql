-- ----------------------------
-- Table structure for admins_role
-- ----------------------------
DROP TABLE IF EXISTS `admin_role`;
CREATE TABLE `admin_role`
(
    `id`       int NOT NULL AUTO_INCREMENT COMMENT 'ID',
    `admin_id` int NOT NULL,
    `role_id`  int NOT NULL,
    PRIMARY KEY (`id`) USING BTREE,
    UNIQUE KEY `uni` (`admin_id`, `role_id`) USING BTREE
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_unicode_ci
  ROW_FORMAT = DYNAMIC;

-- ----------------------------
-- Table structure for roles
-- ----------------------------
DROP TABLE IF EXISTS `role`;
CREATE TABLE `role`
(
    `id`          int          NOT NULL AUTO_INCREMENT,
    `name`        varchar(64)  NOT NULL DEFAULT '' COMMENT '角色名称',
    `description` varchar(128) NOT NULL DEFAULT '' COMMENT '角色描述',
    `created_at`  datetime              DEFAULT CURRENT_TIMESTAMP,
    `updated_at`  datetime              DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    `deleted_at`  datetime              DEFAULT NULL,
    PRIMARY KEY (`id`) USING BTREE,
    UNIQUE KEY `name` (`name`) USING BTREE
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_unicode_ci
  ROW_FORMAT = DYNAMIC;

-- ----------------------------
-- Table structure for permissions
-- ----------------------------
DROP TABLE IF EXISTS `permission`;
CREATE TABLE `permission`
(
    `id`          int          NOT NULL AUTO_INCREMENT,
    `name`        varchar(64)  NOT NULL DEFAULT '' COMMENT '权限名称',
    `description` varchar(128) NOT NULL DEFAULT '' COMMENT '描述',
    `path`        VARCHAR(128) NULL DEFAULT NULL COMMENT '关联的路径',
    `method`      VARCHAR(16) NULL DEFAULT NULL COMMENT '关联的请求方法',
    `created_at`  datetime              DEFAULT CURRENT_TIMESTAMP,
    `updated_at`  datetime              DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    `deleted_at`  datetime              DEFAULT NULL,
    PRIMARY KEY (`id`) USING BTREE,
    UNIQUE KEY `name` (`name`) USING BTREE,
    UNIQUE KEY `uni` (`path`, `method`) USING BTREE

) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_unicode_ci
  ROW_FORMAT = DYNAMIC;

-- ----------------------------
-- Table structure for role_permissions
-- ----------------------------
DROP TABLE IF EXISTS `role_permission`;
CREATE TABLE `role_permission`
(
    `id`            int NOT NULL AUTO_INCREMENT,
    `role_id`       int NOT NULL DEFAULT '0' COMMENT '角色ID',
    `permission_id` int NOT NULL DEFAULT '0' COMMENT '权限ID',
    PRIMARY KEY (`id`),
    UNIQUE KEY `uni`(`role_id`,`permission_id`) USING BTREE

) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_unicode_ci
  ROW_FORMAT = DYNAMIC;

-- ----------------------------
-- Table structure for user
-- ----------------------------
DROP TABLE IF EXISTS `menu`;
CREATE TABLE `menu`
(
    `id`            int         NOT NULL AUTO_INCREMENT,
    `pid`           int         NOT NULL COMMENT '父ID',
    `permission_id` int         NOT NULL,
    `name`          VARCHAR(50) NOT NULL COMMENT '标签name',
    `icon`          VARCHAR(32) NULL DEFAULT NULL COMMENT '图标',
    `type`          tinyint(2) NOT NULL COMMENT '类型 0菜单 1按钮',
    `order`         int NULL DEFAULT NULL COMMENT '排序,倒序',
    `created_at`    datetime DEFAULT CURRENT_TIMESTAMP,
    `updated_at`    datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    `deleted_at`    datetime DEFAULT NULL,
    PRIMARY KEY (`id`),
    KEY             `permission_id`(`permission_id`) USING BTREE
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_unicode_ci
  ROW_FORMAT = DYNAMIC;

SET
FOREIGN_KEY_CHECKS = 1;