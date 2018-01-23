/*
Navicat MySQL Data Transfer

Source Server         : 192.168.200.248
Source Server Version : 50711
Source Host           : 192.168.200.248:3306
Source Database       : test

Target Server Type    : MYSQL
Target Server Version : 50711
File Encoding         : 65001

Date: 2018-01-18 17:14:18
*/

SET FOREIGN_KEY_CHECKS=0;

-- ----------------------------
-- Table structure for area
-- ----------------------------
DROP TABLE IF EXISTS `area`;
CREATE TABLE `area` (
  `uid` int(11) NOT NULL,
  `province` varchar(255) DEFAULT '湖北',
  `city` varchar(255) DEFAULT '武汉',
  PRIMARY KEY (`uid`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- ----------------------------
-- Records of area
-- ----------------------------
INSERT INTO `area` VALUES ('1', '湖北', '武汉');
INSERT INTO `area` VALUES ('2', '湖北', '武汉');
INSERT INTO `area` VALUES ('3', '湖北', '武汉');
INSERT INTO `area` VALUES ('4', '湖北', '武汉');
INSERT INTO `area` VALUES ('5', '湖北', '武汉');
