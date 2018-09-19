/*
Navicat MySQL Data Transfer

Source Server         : 192.168.200.248
Source Server Version : 50711
Source Host           : 192.168.200.248:3306
Source Database       : test

Target Server Type    : MYSQL
Target Server Version : 50711
File Encoding         : 65001

Date: 2018-01-30 14:17:45
*/

SET FOREIGN_KEY_CHECKS=0;

-- ----------------------------
-- Table structure for users
-- ----------------------------
DROP TABLE IF EXISTS `users`;
CREATE TABLE `users` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `name` varchar(255) DEFAULT 'fizzday',
  `age` int(11) DEFAULT '18',
  `website` varchar(255) DEFAULT 'http://fizzday.net',
  `job` varchar(255) DEFAULT 'go',
  `money` decimal(10,2) DEFAULT '5.00',
  `created_at` datetime DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=17 DEFAULT CHARSET=utf8;

-- ----------------------------
-- Records of users
-- ----------------------------
INSERT INTO `users` VALUES ('1', 'fizzday', '18', 'http://fizzday.net', 'go', '88.00', '2018-01-26 20:00:31', '2018-01-30 13:34:54');
INSERT INTO `users` VALUES ('2', 'fizzday', '18', 'http://fizzday.net', 'go', '2.00', '2018-01-26 20:00:34', '2018-01-26 20:00:34');
INSERT INTO `users` VALUES ('3', 'fizzday', '18', 'http://fizzday.net', 'go', '3.00', '2018-01-26 20:00:37', '2018-01-26 20:00:37');
INSERT INTO `users` VALUES ('4', 'fizzday', '18', 'http://fizzday.net', 'go', '4.00', '2018-01-26 20:00:43', '2018-01-26 20:00:43');
INSERT INTO `users` VALUES ('5', 'fizzday', '18', 'http://fizzday.net', 'go', '5.00', '2018-01-26 20:00:45', '2018-01-26 20:00:45');
INSERT INTO `users` VALUES ('6', 'fizz5', '17', 'http://fizzday.net', 'it3', '5.00', '2018-01-30 12:41:52', '2018-01-30 12:41:52');
INSERT INTO `users` VALUES ('7', 'fizz5', '17', 'http://fizzday.net', 'it3', '5.00', '2018-01-30 12:43:49', '2018-01-30 12:43:49');
