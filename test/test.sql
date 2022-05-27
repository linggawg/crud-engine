/*
 Navicat Premium Data Transfer

 Source Server         : mysql_local
 Source Server Type    : MySQL
 Source Server Version : 80027
 Source Host           : localhost:3306
 Source Schema         : crud

 Target Server Type    : MySQL
 Target Server Version : 80027
 File Encoding         : 65001

 Date: 24/05/2022 09:18:40
*/

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- Table structure for province
-- ----------------------------
DROP TABLE IF EXISTS `province`;
CREATE TABLE `province`  (
  `id` int NOT NULL,
  `name` varchar(255) CHARACTER SET latin1 COLLATE latin1_swedish_ci NULL DEFAULT NULL,
  `created_by` varchar(255) CHARACTER SET latin1 COLLATE latin1_swedish_ci NULL DEFAULT NULL,
  `created_date` timestamp NULL DEFAULT NULL,
  `is_deleted` smallint NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB CHARACTER SET = latin1 COLLATE = latin1_swedish_ci ROW_FORMAT = DYNAMIC;

-- ----------------------------
-- Records of province
-- ----------------------------
INSERT INTO `province` VALUES (1, 'Aceh', 'system', '2021-10-15 21:09:06', 0);
INSERT INTO `province` VALUES (2, 'Sumatra Utara', 'system', '2021-10-15 21:09:07', 0);
INSERT INTO `province` VALUES (3, 'Sumatra Barat', 'system', '2021-10-15 21:09:08', 0);
INSERT INTO `province` VALUES (4, 'Riau', 'system', '2021-10-15 21:09:08', 0);
INSERT INTO `province` VALUES (5, 'Kepulauan Riau', 'system', '2021-10-15 21:09:09', 0);
INSERT INTO `province` VALUES (6, 'Jambi', 'system', '2021-10-15 21:09:09', 0);
INSERT INTO `province` VALUES (7, 'Bengkulu', 'system', '2021-10-15 21:09:10', 0);
INSERT INTO `province` VALUES (8, 'Sumatra Selatan', 'system', '2021-10-15 21:09:10', 0);
INSERT INTO `province` VALUES (9, 'Kepulauan Bangka Belitung', 'system', '2021-10-15 21:09:11', 0);
INSERT INTO `province` VALUES (10, 'Lampung', 'system', '2021-10-15 21:09:11', 0);
INSERT INTO `province` VALUES (11, 'DKI Jakarta', 'system', '2021-10-15 21:09:12', 0);
INSERT INTO `province` VALUES (12, 'Banten', 'system', '2021-10-15 21:09:12', 0);
INSERT INTO `province` VALUES (13, 'Jawa Barat', 'system', '2021-10-15 21:09:12', 0);
INSERT INTO `province` VALUES (14, 'Jawa Tengah', 'system', '2021-10-15 21:09:13', 0);
INSERT INTO `province` VALUES (15, 'Daerah Istimewa Yogyakarta', 'system', '2021-10-15 21:09:14', 0);
INSERT INTO `province` VALUES (16, 'Jawa Timur', 'system', '2021-10-15 21:09:14', 0);
INSERT INTO `province` VALUES (17, 'Bali', 'system', '2021-10-15 21:09:15', 0);
INSERT INTO `province` VALUES (18, 'Nusa Tenggara Barat', 'system', '2021-10-15 21:09:15', 0);
INSERT INTO `province` VALUES (19, 'Nusa Tenggara Timur', 'system', '2021-10-15 21:09:16', 0);
INSERT INTO `province` VALUES (20, 'Kalimantan Barat', 'system', '2021-10-15 21:09:16', 0);
INSERT INTO `province` VALUES (21, 'Kalimantan Tengah', 'system', '2021-10-15 21:09:17', 0);
INSERT INTO `province` VALUES (22, 'Kalimantan Selatan', 'system', '2021-10-15 21:09:17', 0);
INSERT INTO `province` VALUES (23, 'Kalimantan Timur', 'system', '2021-10-15 21:09:18', 0);
INSERT INTO `province` VALUES (24, 'Kalimantan Utara', 'system', '2021-10-15 21:09:18', 0);
INSERT INTO `province` VALUES (25, 'Sulawesi Utara', 'system', '2021-10-15 21:09:19', 0);
INSERT INTO `province` VALUES (26, 'Gorontalo', 'system', '2021-10-15 21:09:19', 0);
INSERT INTO `province` VALUES (27, 'Sulawesi Tengah', 'system', '2021-10-15 21:09:20', 0);
INSERT INTO `province` VALUES (28, 'Sulawesi Barat', 'system', '2021-10-15 21:09:20', 0);
INSERT INTO `province` VALUES (29, 'Sulawesi Selatan', 'system', '2021-10-15 21:09:21', 0);
INSERT INTO `province` VALUES (30, 'Sulawesi Tenggara', 'system', '2021-10-15 21:09:21', 0);
INSERT INTO `province` VALUES (31, 'Maluku Utara', 'system', '2021-10-15 21:09:22', 0);
INSERT INTO `province` VALUES (32, 'Maluku', 'system', '2021-10-15 21:09:22', 0);
INSERT INTO `province` VALUES (33, 'Papua Barat', 'system', '2021-10-15 21:09:22', 0);
INSERT INTO `province` VALUES (34, 'Papua', 'system', '2021-10-15 21:09:23', 0);

SET FOREIGN_KEY_CHECKS = 1;
