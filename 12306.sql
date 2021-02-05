-- --------------------------------------------------------
-- 主机:                           127.0.0.1
-- 服务器版本:                        5.7.22-log - MySQL Community Server (GPL)
-- 服务器操作系统:                      Win64
-- HeidiSQL 版本:                  11.0.0.5919
-- --------------------------------------------------------

/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET NAMES utf8 */;
/*!50503 SET NAMES utf8mb4 */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;

-- 导出  表 12306.orders 结构
CREATE TABLE IF NOT EXISTS `orders` (
  `id` int(255) unsigned NOT NULL AUTO_INCREMENT,
  `user_id` int(10) unsigned DEFAULT NULL COMMENT '用户id',
  `seat_no` int(10) unsigned DEFAULT NULL COMMENT '座位id',
  `start_station_no` int(10) unsigned DEFAULT NULL COMMENT '开始站点号',
  `end_station_no` int(10) unsigned DEFAULT NULL COMMENT '结束站点号',
  `start_station` varchar(255) DEFAULT NULL COMMENT '开始站点名字',
  `end_station` varchar(255) DEFAULT NULL COMMENT '结束站点名字',
  `date` datetime DEFAULT NULL COMMENT '车票日期',
  `status` varchar(255) DEFAULT NULL COMMENT '状态',
  `order_id` varchar(255) DEFAULT NULL,
  `trip_id` int(10) unsigned DEFAULT NULL,
  `seat_id` int(10) unsigned DEFAULT NULL,
  `seat_catogory` varchar(50) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=62 DEFAULT CHARSET=utf8mb4;

-- 数据导出被取消选择。

-- 导出  表 12306.stations 结构
CREATE TABLE IF NOT EXISTS `stations` (
  `id` int(255) unsigned NOT NULL AUTO_INCREMENT,
  `spell` varchar(255) DEFAULT NULL COMMENT '拼音',
  `name` varchar(255) DEFAULT NULL COMMENT '名字',
  `is_hot` int(11) DEFAULT NULL COMMENT '是否热门站点',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=4 DEFAULT CHARSET=utf8mb4;

-- 数据导出被取消选择。

-- 导出  表 12306.trains 结构
CREATE TABLE IF NOT EXISTS `trains` (
  `id` int(255) NOT NULL AUTO_INCREMENT,
  `train_number` varchar(50) DEFAULT NULL,
  `catogory` varchar(50) DEFAULT NULL COMMENT '列车类型',
  `status` int(8) DEFAULT NULL COMMENT '列车状态',
  `station_nums` int(11) DEFAULT NULL COMMENT '经过的站点数量',
  `is_fast` int(11) DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=21 DEFAULT CHARSET=utf8mb4;

-- 数据导出被取消选择。

-- 导出  表 12306.train_seat 结构
CREATE TABLE IF NOT EXISTS `train_seat` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `trip_id` int(11) DEFAULT NULL,
  `seat_catogory` varchar(50) DEFAULT NULL,
  `seat_no` int(11) DEFAULT NULL COMMENT '从1开始',
  `seat_name` varchar(50) DEFAULT NULL COMMENT '比如12A，表示12排的A座位号',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- 数据导出被取消选择。

-- 导出  表 12306.train_station_pair 结构
CREATE TABLE IF NOT EXISTS `train_station_pair` (
  `id` int(255) NOT NULL AUTO_INCREMENT,
  `train_id` int(255) DEFAULT NULL COMMENT 'tripID',
  `start_station_id` int(255) DEFAULT NULL COMMENT '站点id',
  `end_station_id` int(255) DEFAULT NULL,
  `start_station_no` int(8) DEFAULT NULL COMMENT '站点序列（该车次的第几个站点，1为起始站）',
  `end_station_no` int(8) DEFAULT NULL,
  `price` int(8) DEFAULT NULL,
  `start_station_name` varchar(50) DEFAULT NULL COMMENT '站点名字',
  `end_station_name` varchar(50) DEFAULT NULL,
  `start_time` time DEFAULT NULL COMMENT '启动时间',
  `end_time` time DEFAULT NULL COMMENT '到站时间',
  PRIMARY KEY (`id`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8mb4 COMMENT='一个train经过的站点对们；？';

-- 数据导出被取消选择。

-- 导出  表 12306.trips 结构
CREATE TABLE IF NOT EXISTS `trips` (
  `id` int(255) NOT NULL AUTO_INCREMENT,
  `train_id` int(11) DEFAULT NULL COMMENT '列车id',
  `start_date` date DEFAULT NULL COMMENT '发车时间',
  PRIMARY KEY (`id`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=21 DEFAULT CHARSET=utf8mb4;

-- 数据导出被取消选择。

-- 导出  表 12306.trip_segment 结构
CREATE TABLE IF NOT EXISTS `trip_segment` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `trip_id` int(11) DEFAULT NULL COMMENT 'TripID',
  `segment_no` int(5) DEFAULT NULL COMMENT '段号，1为第一段',
  `seat_bytes` varbinary(500) DEFAULT NULL COMMENT '该段对应的seat_catogory座位类型的座位',
  `seat_catogory` varchar(50) DEFAULT NULL COMMENT '座位类型',
  KEY `索引 1` (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=15 DEFAULT CHARSET=utf8mb4;

-- 数据导出被取消选择。

-- 导出  表 12306.users 结构
CREATE TABLE IF NOT EXISTS `users` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `user_name` varchar(20) NOT NULL COMMENT '用户名',
  `telephone` varchar(255) NOT NULL COMMENT '电话号码',
  `password` varchar(255) NOT NULL COMMENT '密码',
  PRIMARY KEY (`id`),
  UNIQUE KEY `telephone` (`telephone`)
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8mb4;

-- 数据导出被取消选择。

/*!40101 SET SQL_MODE=IFNULL(@OLD_SQL_MODE, '') */;
/*!40014 SET FOREIGN_KEY_CHECKS=IF(@OLD_FOREIGN_KEY_CHECKS IS NULL, 1, @OLD_FOREIGN_KEY_CHECKS) */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
