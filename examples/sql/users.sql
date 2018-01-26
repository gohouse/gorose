/*
Navicat MySQL Data Transfer

Source Server         : localhost
Source Server Version : 50631
Source Host           : localhost:3306
Source Database       : test

Target Server Type    : MYSQL
Target Server Version : 50631
File Encoding         : 65001

Date: 2018-01-16 13:25:28
*/

SET FOREIGN_KEY_CHECKS=0;

-- ----------------------------
-- Table structure for users
-- ----------------------------
DROP TABLE IF EXISTS `users`;
CREATE TABLE `users` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `name` varchar(255) DEFAULT NULL,
  `age` int(11) DEFAULT NULL,
  `website` varchar(255) DEFAULT NULL,
  `job` varchar(255) DEFAULT NULL,
  `money` decimal(10,2) DEFAULT '5.00',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=16 DEFAULT CHARSET=latin1;

-- ----------------------------
-- Records of users
-- ----------------------------
INSERT INTO `users` VALUES ('1', 'fizz', '18', 'fizzday.net', 'it', '5.00');
INSERT INTO `users` VALUES ('2', 'fizzday', '18', 'fizzday.net', 'engineer', '5.00');
INSERT INTO `users` VALUES ('3', 'gorose', '18', 'go-rose.com', 'go orm', '5.00');
INSERT INTO `users` VALUES ('4', 'fizz3', '3', null, null, '5.00');
INSERT INTO `users` VALUES ('5', 'fizz3', '3', null, 'sadf', '5.00');
INSERT INTO `users` VALUES ('8', null, null, null, 'sadf', '5.00');
INSERT INTO `users` VALUES ('9', null, null, '1', 'eee', '5.00');
INSERT INTO `users` VALUES ('10', null, null, 'asdf', 'eee', '5.00');
INSERT INTO `users` VALUES ('11', null, null, 'asdf', 'eee', '5.00');
INSERT INTO `users` VALUES ('12', 'fizz3', '3', null, null, '5.00');
INSERT INTO `users` VALUES ('13', 'fizz3', '3', null, null, '5.00');
INSERT INTO `users` VALUES ('14', 'fizz4', '17', null, 'it3', '5.00');
INSERT INTO `users` VALUES ('15', 'fizz4', '17', null, 'it3', '5.00');
