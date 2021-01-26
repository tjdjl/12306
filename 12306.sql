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
  `order_id` varchar(255) DEFAULT NULL,
  `trip_id` int(10) unsigned DEFAULT NULL COMMENT '车次id',
  `user_id` int(10) unsigned DEFAULT NULL COMMENT '用户id',
  `seat_id` int(10) unsigned DEFAULT NULL COMMENT '座位id',
  `start_no` int(10) unsigned DEFAULT NULL COMMENT '开始站点号',
  `end_no` int(10) unsigned DEFAULT NULL COMMENT '结束站点号',
  `start_station` varchar(255) DEFAULT NULL COMMENT '开始站点名字',
  `end_station` varchar(255) DEFAULT NULL COMMENT '结束站点名字',
  `date` datetime DEFAULT NULL COMMENT '车票日期',
  `status` varchar(255) DEFAULT NULL COMMENT '状态',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8;

-- 数据导出被取消选择。

-- 导出  表 12306.stations 结构
CREATE TABLE IF NOT EXISTS `stations` (
  `id` int(255) unsigned NOT NULL AUTO_INCREMENT,
  `spell` varchar(255) DEFAULT NULL COMMENT '拼音',
  `name` varchar(255) DEFAULT NULL COMMENT '名字',
  `is_hot` int(11) DEFAULT NULL COMMENT '是否热门站点',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=4 DEFAULT CHARSET=utf8;

-- 数据导出被取消选择。

-- 导出  表 12306.trains 结构
CREATE TABLE IF NOT EXISTS `trains` (
  `id` int(255) NOT NULL AUTO_INCREMENT,
  `name` varchar(255) DEFAULT NULL COMMENT '列车名',
  `type` int(8) DEFAULT NULL COMMENT '列车类型',
  `carriage_num` int(8) DEFAULT NULL COMMENT '车厢数',
  `status` int(8) DEFAULT NULL COMMENT '列车状态',
  PRIMARY KEY (`id`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- 数据导出被取消选择。

-- 导出  表 12306.trips 结构
CREATE TABLE IF NOT EXISTS `trips` (
  `id` int(255) NOT NULL AUTO_INCREMENT,
  `tripId` varchar(255) DEFAULT NULL COMMENT '车次号',
  `trainId` int(11) DEFAULT NULL COMMENT '列车id',
  `start_station_Id` int(11) DEFAULT NULL COMMENT '起点站id',
  `start_station` varchar(255) DEFAULT NULL COMMENT '起点站',
  `end_station_Id` int(11) DEFAULT NULL COMMENT '终点站id',
  `end_station` varchar(255) DEFAULT NULL COMMENT '终点站',
  `departure_time` datetime DEFAULT NULL COMMENT '发车时间',
  `station_num` int(11) DEFAULT NULL COMMENT '站台总数目',
  PRIMARY KEY (`id`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- 数据导出被取消选择。

-- 导出  表 12306.trip_price 结构
CREATE TABLE IF NOT EXISTS `trip_price` (
  `id` int(255) NOT NULL AUTO_INCREMENT,
  `tripId` int(255) DEFAULT NULL COMMENT '车次id',
  `station_first` int(255) DEFAULT NULL COMMENT '站点一',
  `station_second` int(255) DEFAULT NULL COMMENT '站点二',
  `price` decimal(10,2) DEFAULT NULL COMMENT '价格',
  PRIMARY KEY (`id`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- 数据导出被取消选择。

-- 导出  表 12306.trip_seat_segment 结构
CREATE TABLE IF NOT EXISTS `trip_seat_segment` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `trip_id` int(10) unsigned DEFAULT NULL COMMENT '车次id',
  `status` int(10) unsigned DEFAULT NULL COMMENT '状态',
  `seat_id` int(10) unsigned DEFAULT NULL COMMENT '座位id',
  `segment` int(10) unsigned DEFAULT NULL COMMENT '第几段',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=4 DEFAULT CHARSET=utf8;

-- 数据导出被取消选择。

-- 导出  表 12306.trip_segment 结构
CREATE TABLE IF NOT EXISTS `trip_segment` (
  `id` int(11) DEFAULT NULL,
  `trip_id` int(11) DEFAULT NULL,
  `segment_no` int(5) DEFAULT NULL,
  `business_seats` varbinary(500) DEFAULT NULL,
  `first_seats` varbinary(500) DEFAULT NULL,
  `second_seats` varbinary(500) DEFAULT NULL,
  `no_satas` varbinary(500) DEFAULT NULL,
  `hard_berth` varbinary(500) DEFAULT NULL,
  `soft_berth` varbinary(500) DEFAULT NULL,
  `senior_soft_berth` varbinary(500) DEFAULT NULL,
  KEY `索引 1` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- 数据导出被取消选择。

-- 导出  表 12306.trip_station 结构
CREATE TABLE IF NOT EXISTS `trip_station` (
  `id` int(255) NOT NULL AUTO_INCREMENT,
  `trip_id` int(255) DEFAULT NULL COMMENT '车次id',
  `station_id` int(255) DEFAULT NULL COMMENT '站点id',
  `station_name` varchar(50) DEFAULT NULL,
  `sequence` int(8) DEFAULT NULL COMMENT '站点序列（该车次的第几个站点）',
  `start_time` datetime DEFAULT NULL COMMENT '启动时间',
  `arrive_time` datetime DEFAULT NULL COMMENT '到站时间',
  PRIMARY KEY (`id`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8;

-- 数据导出被取消选择。

-- 导出  表 12306.users 结构
CREATE TABLE IF NOT EXISTS `users` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `user_name` varchar(20) NOT NULL COMMENT '用户名',
  `telephone` varchar(255) NOT NULL COMMENT '电话号码',
  `password` varchar(255) NOT NULL COMMENT '密码',
  PRIMARY KEY (`id`),
  UNIQUE KEY `telephone` (`telephone`)
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8;

-- 数据导出被取消选择。

/*!40101 SET SQL_MODE=IFNULL(@OLD_SQL_MODE, '') */;
/*!40014 SET FOREIGN_KEY_CHECKS=IF(@OLD_FOREIGN_KEY_CHECKS IS NULL, 1, @OLD_FOREIGN_KEY_CHECKS) */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
