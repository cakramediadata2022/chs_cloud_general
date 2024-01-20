-- MySQL dump 10.13  Distrib 5.7.30, for Win32 (AMD64)
--
-- Host: localhost    Database: db_main_cloud
-- ------------------------------------------------------
-- Server version	5.5.5-10.6.11-MariaDB

/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!40101 SET NAMES utf8 */;
/*!40103 SET @OLD_TIME_ZONE=@@TIME_ZONE */;
/*!40103 SET TIME_ZONE='+00:00' */;
/*!40014 SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0 */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;

--
-- Current Database: `db_main_cloud`
--

/*!40000 DROP DATABASE IF EXISTS `db_main_cloud`*/;

CREATE DATABASE /*!32312 IF NOT EXISTS*/ `db_main_cloud` /*!40100 DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci */;

USE `db_main_cloud`;

--
-- Table structure for table `chain`
--

DROP TABLE IF EXISTS `chain`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `chain` (
  `chain_id` int(11) NOT NULL AUTO_INCREMENT,
  `name` varchar(100) NOT NULL,
  `address` varchar(100) DEFAULT NULL,
  `city_code` varchar(20) DEFAULT NULL,
  `state_code` varchar(20) DEFAULT NULL,
  `country_code` varchar(20) DEFAULT NULL,
  `created_at` timestamp NULL DEFAULT current_timestamp(),
  `updated_at` timestamp NULL DEFAULT '0000-00-00 00:00:00' ON UPDATE current_timestamp(),
  PRIMARY KEY (`chain_id`)
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `chain`
--

LOCK TABLES `chain` WRITE;
/*!40000 ALTER TABLE `chain` DISABLE KEYS */;
INSERT INTO `chain` VALUES (1,'MPHG',NULL,NULL,NULL,NULL,'2023-06-13 03:36:18','0000-00-00 00:00:00');
/*!40000 ALTER TABLE `chain` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `client`
--

DROP TABLE IF EXISTS `client`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `client` (
  `client_id` int(11) NOT NULL AUTO_INCREMENT,
  `name` varchar(100) DEFAULT NULL,
  `email` varchar(100) DEFAULT NULL,
  `phone` varchar(50) DEFAULT NULL,
  `username` varchar(50) DEFAULT NULL,
  `password` varchar(100) DEFAULT NULL,
  `contact_info` varchar(100) DEFAULT NULL,
  `created_at` timestamp NULL DEFAULT current_timestamp(),
  `updated_at` timestamp NULL DEFAULT '0000-00-00 00:00:00' ON UPDATE current_timestamp(),
  PRIMARY KEY (`client_id`)
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `client`
--

LOCK TABLES `client` WRITE;
/*!40000 ALTER TABLE `client` DISABLE KEYS */;
INSERT INTO `client` VALUES (1,'Cakrasoft','admin@cakrasoft.com',NULL,'cakrasoft',NULL,NULL,'2023-06-30 09:48:05','0000-00-00 00:00:00');
/*!40000 ALTER TABLE `client` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `company`
--

DROP TABLE IF EXISTS `company`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `company` (
  `company_id` bigint(20) NOT NULL AUTO_INCREMENT,
  `client_id` int(11) DEFAULT NULL,
  `chain_id` int(11) DEFAULT NULL,
  `company_code` varchar(10) NOT NULL,
  `name` varchar(100) NOT NULL,
  `property_name` varchar(100) NOT NULL,
  `domain` varchar(100) DEFAULT NULL,
  `subdomain` varchar(100) NOT NULL,
  `city_code` varchar(20) DEFAULT NULL,
  `state_code` varchar(20) DEFAULT NULL,
  `country_code` varchar(20) DEFAULT NULL,
  `address` varchar(100) DEFAULT NULL,
  `rooms` int(11) DEFAULT NULL,
  `max_user` int(11) DEFAULT NULL,
  `is_complete_wizard` tinyint(1) DEFAULT 0,
  `created_at` timestamp NULL DEFAULT current_timestamp(),
  `updated_at` timestamp NULL DEFAULT '0000-00-00 00:00:00' ON UPDATE current_timestamp(),
  PRIMARY KEY (`company_id`),
  UNIQUE KEY `subdomain` (`subdomain`),
  UNIQUE KEY `company_code` (`company_code`),
  UNIQUE KEY `company_code_2` (`company_code`),
  UNIQUE KEY `subdomain_2` (`subdomain`),
  UNIQUE KEY `domain` (`domain`),
  UNIQUE KEY `domain_2` (`domain`),
  KEY `company_fk1` (`client_id`),
  KEY `company_fk2` (`chain_id`),
  CONSTRAINT `company_fk1` FOREIGN KEY (`client_id`) REFERENCES `client` (`client_id`),
  CONSTRAINT `company_fk2` FOREIGN KEY (`chain_id`) REFERENCES `chain` (`chain_id`)
) ENGINE=InnoDB AUTO_INCREMENT=5 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `company`
--

LOCK TABLES `company` WRITE;
/*!40000 ALTER TABLE `company` DISABLE KEYS */;
INSERT INTO `company` VALUES (1,1,NULL,'CMD','PT Cakra Media Data','Cakra Hotel','localhost:3000','demo.cakrasoft.net','SS','SS','ID','Jl Mambal',50,NULL,NULL,'2023-06-30 09:46:56','0000-00-00 00:00:00'),(2,1,NULL,'CMA','Bakra','Bakra','cakrasoft.net','cmd.localhost:3000',NULL,NULL,NULL,NULL,NULL,NULL,NULL,'2023-06-30 09:46:56','0000-00-00 00:00:00'),(3,1,NULL,'CMC','Bakra','Bakra','s','bmd.localhost:3000',NULL,NULL,NULL,NULL,NULL,NULL,NULL,'2023-06-30 09:46:56','0000-00-00 00:00:00'),(4,1,NULL,'CMB','Bakra','Bakra',NULL,'dmd',NULL,NULL,NULL,NULL,NULL,NULL,NULL,'2023-06-30 09:46:56','0000-00-00 00:00:00');
/*!40000 ALTER TABLE `company` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `company_database`
--

DROP TABLE IF EXISTS `company_database`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `company_database` (
  `company_id` bigint(11) NOT NULL,
  `database_name` varchar(100) NOT NULL,
  `host` varchar(100) NOT NULL,
  `port` int(11) NOT NULL,
  `username` varchar(50) NOT NULL,
  `password` varchar(100) NOT NULL,
  PRIMARY KEY (`company_id`),
  UNIQUE KEY `database_name` (`database_name`),
  CONSTRAINT `company_database_fk1` FOREIGN KEY (`company_id`) REFERENCES `company` (`company_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci ROW_FORMAT=DYNAMIC;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `company_database`
--

LOCK TABLES `company_database` WRITE;
/*!40000 ALTER TABLE `company_database` DISABLE KEYS */;
INSERT INTO `company_database` VALUES (1,'db_chs_cloud','localhost',3308,'root','kalomang'),(2,'db_client1','localhost',3308,'root','kalomang'),(3,'db_client2','localhost',3308,'root','kalomang'),(4,'db_client3','localhost',3308,'root','kalomang');
/*!40000 ALTER TABLE `company_database` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `feature`
--

DROP TABLE IF EXISTS `feature`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `feature` (
  `feature_id` int(11) NOT NULL AUTO_INCREMENT,
  `feature_name` varchar(50) DEFAULT NULL,
  `description` varchar(100) DEFAULT NULL,
  PRIMARY KEY (`feature_id`)
) ENGINE=InnoDB AUTO_INCREMENT=6 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `feature`
--

LOCK TABLES `feature` WRITE;
/*!40000 ALTER TABLE `feature` DISABLE KEYS */;
INSERT INTO `feature` VALUES (1,'Keylock',NULL),(2,'IPTV',NULL),(3,'PABX',NULL),(4,'MikroTik',NULL),(5,'Channel Manager',NULL);
/*!40000 ALTER TABLE `feature` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `module`
--

DROP TABLE IF EXISTS `module`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `module` (
  `module_id` int(11) NOT NULL AUTO_INCREMENT,
  `name` varchar(100) DEFAULT NULL,
  `description` text DEFAULT NULL,
  PRIMARY KEY (`module_id`)
) ENGINE=InnoDB AUTO_INCREMENT=11 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `module`
--

LOCK TABLES `module` WRITE;
/*!40000 ALTER TABLE `module` DISABLE KEYS */;
INSERT INTO `module` VALUES (6,'Front Desk',NULL),(7,'Point of Sales',NULL),(8,'Banquet',NULL),(9,'Accounting',NULL),(10,'Inventory & Fix Asset ',NULL);
/*!40000 ALTER TABLE `module` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `plan`
--

DROP TABLE IF EXISTS `plan`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `plan` (
  `plan_id` int(11) NOT NULL AUTO_INCREMENT,
  `plan_name` varchar(50) DEFAULT NULL,
  `description` varchar(100) DEFAULT NULL,
  `price` decimal(10,2) DEFAULT NULL,
  PRIMARY KEY (`plan_id`)
) ENGINE=InnoDB AUTO_INCREMENT=7 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `plan`
--

LOCK TABLES `plan` WRITE;
/*!40000 ALTER TABLE `plan` DISABLE KEYS */;
INSERT INTO `plan` VALUES (1,'Free Plan',NULL,NULL),(2,'Basic Plan',NULL,NULL),(3,'Standard Plan',NULL,NULL),(4,'Professional Plan',NULL,NULL),(5,'Advanced Plan',NULL,NULL),(6,'Enterprise Plan',NULL,NULL);
/*!40000 ALTER TABLE `plan` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `subscription`
--

DROP TABLE IF EXISTS `subscription`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `subscription` (
  `subscription_id` int(11) NOT NULL AUTO_INCREMENT,
  `client_id` int(11) DEFAULT NULL,
  `company_id` bigint(20) DEFAULT NULL,
  `plan_id` int(11) DEFAULT NULL,
  `start_date` date DEFAULT NULL,
  `end_date` date DEFAULT NULL,
  `payment_status` varchar(20) DEFAULT NULL,
  PRIMARY KEY (`subscription_id`),
  KEY `subscription_ibfk_2` (`plan_id`),
  KEY `subscription_fk1` (`client_id`),
  KEY `subscription_fk2` (`company_id`),
  CONSTRAINT `subscription_fk1` FOREIGN KEY (`client_id`) REFERENCES `client` (`client_id`),
  CONSTRAINT `subscription_fk2` FOREIGN KEY (`company_id`) REFERENCES `company` (`company_id`),
  CONSTRAINT `subscription_ibfk_2` FOREIGN KEY (`plan_id`) REFERENCES `plan` (`plan_id`)
) ENGINE=InnoDB AUTO_INCREMENT=4 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `subscription`
--

LOCK TABLES `subscription` WRITE;
/*!40000 ALTER TABLE `subscription` DISABLE KEYS */;
INSERT INTO `subscription` VALUES (1,1,1,1,'2023-01-01','2024-07-06','paid'),(2,1,2,1,'2023-01-01','2023-09-16','paid'),(3,1,3,1,'2023-01-01','2023-02-01','paid');
/*!40000 ALTER TABLE `subscription` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `subscription_feature`
--

DROP TABLE IF EXISTS `subscription_feature`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `subscription_feature` (
  `subscription_id` int(11) DEFAULT NULL,
  `feature_id` int(11) DEFAULT NULL,
  KEY `subscription_feature_ibfk_1` (`subscription_id`),
  KEY `subscription_feature_ibfk_2` (`feature_id`),
  CONSTRAINT `subscription_feature_ibfk_1` FOREIGN KEY (`subscription_id`) REFERENCES `subscription` (`subscription_id`),
  CONSTRAINT `subscription_feature_ibfk_2` FOREIGN KEY (`feature_id`) REFERENCES `feature` (`feature_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `subscription_feature`
--

LOCK TABLES `subscription_feature` WRITE;
/*!40000 ALTER TABLE `subscription_feature` DISABLE KEYS */;
/*!40000 ALTER TABLE `subscription_feature` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `timezone`
--

DROP TABLE IF EXISTS `timezone`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `timezone` (
  `timezone` varchar(50) DEFAULT NULL,
  `offset` varchar(20) DEFAULT NULL,
  `offset_dst` varchar(20) DEFAULT NULL,
  `id` int(11) NOT NULL AUTO_INCREMENT,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=589 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `timezone`
--

LOCK TABLES `timezone` WRITE;
/*!40000 ALTER TABLE `timezone` DISABLE KEYS */;
INSERT INTO `timezone` VALUES ('Africa/Abidjan','0','0',1),('Africa/Accra','0','0',2),('Africa/Addis_Ababa','10800','10800',3),('Africa/Algiers','3600','3600',4),('Africa/Asmara','10800','10800',5),('Africa/Asmera','10800','10800',6),('Africa/Bamako','0','0',7),('Africa/Bangui','3600','3600',8),('Africa/Banjul','0','0',9),('Africa/Bissau','0','0',10),('Africa/Blantyre','7200','7200',11),('Africa/Brazzaville','3600','3600',12),('Africa/Bujumbura','7200','7200',13),('Africa/Cairo','7200','7200',14),('Africa/Casablanca','0','3600',15),('Africa/Ceuta','3600','7200',16),('Africa/Conakry','0','0',17),('Africa/Dakar','0','0',18),('Africa/Dar_es_Salaam','10800','10800',19),('Africa/Djibouti','10800','10800',20),('Africa/Douala','3600','3600',21),('Africa/El_Aaiun','0','3600',22),('Africa/Freetown','0','0',23),('Africa/Gaborone','7200','7200',24),('Africa/Harare','7200','7200',25),('Africa/Johannesburg','7200','7200',26),('Africa/Juba','10800','10800',27),('Africa/Kampala','10800','10800',28),('Africa/Khartoum','10800','10800',29),('Africa/Kigali','7200','7200',30),('Africa/Kinshasa','3600','3600',31),('Africa/Lagos','3600','3600',32),('Africa/Libreville','3600','3600',33),('Africa/Lome','0','0',34),('Africa/Luanda','3600','3600',35),('Africa/Lubumbashi','7200','7200',36),('Africa/Lusaka','7200','7200',37),('Africa/Malabo','3600','3600',38),('Africa/Maputo','7200','7200',39),('Africa/Maseru','7200','7200',40),('Africa/Mbabane','7200','7200',41),('Africa/Mogadishu','10800','10800',42),('Africa/Monrovia','0','0',43),('Africa/Nairobi','10800','10800',44),('Africa/Ndjamena','3600','3600',45),('Africa/Niamey','3600','3600',46),('Africa/Nouakchott','0','0',47),('Africa/Ouagadougou','0','0',48),('Africa/Porto-Novo','3600','3600',49),('Africa/Sao_Tome','0','0',50),('Africa/Timbuktu','0','0',51),('Africa/Tripoli','7200','7200',52),('Africa/Tunis','3600','3600',53),('Africa/Windhoek','3600','7200',54),('America/Adak','-36000','-32400',55),('America/Anchorage','-32400','-28800',56),('America/Anguilla','-14400','-14400',57),('America/Antigua','-14400','-14400',58),('America/Araguaina','-10800','-10800',59),('America/Argentina/Buenos_Aires','-10800','-10800',60),('America/Argentina/Catamarca','-10800','-10800',61),('America/Argentina/ComodRivadavia','-10800','-10800',62),('America/Argentina/Cordoba','-10800','-10800',63),('America/Argentina/Jujuy','-10800','-10800',64),('America/Argentina/La_Rioja','-10800','-10800',65),('America/Argentina/Mendoza','-10800','-10800',66),('America/Argentina/Rio_Gallegos','-10800','-10800',67),('America/Argentina/Salta','-10800','-10800',68),('America/Argentina/San_Juan','-10800','-10800',69),('America/Argentina/San_Luis','-10800','-10800',70),('America/Argentina/Tucuman','-10800','-10800',71),('America/Argentina/Ushuaia','-10800','-10800',72),('America/Aruba','-14400','-14400',73),('America/Asuncion','-10800','-10800',74),('America/Atikokan','-18000','-18000',75),('America/Atka','-36000','-32400',76),('America/Bahia','-10800','-10800',77),('America/Bahia_Banderas','-21600','-18000',78),('America/Barbados','-14400','-14400',79),('America/Belem','-10800','-10800',80),('America/Belize','-21600','-21600',81),('America/Blanc-Sablon','-14400','-14400',82),('America/Boa_Vista','-14400','-14400',83),('America/Bogota','-18000','-18000',84),('America/Boise','-25200','-21600',85),('America/Buenos_Aires','-10800','-10800',86),('America/Cambridge_Bay','-25200','-21600',87),('America/Campo_Grande','-14400','-10800',88),('America/Cancun','-18000','-18000',89),('America/Caracas','-16200','-16200',90),('America/Catamarca','-10800','-10800',91),('America/Cayenne','-10800','-10800',92),('America/Cayman','-18000','-18000',93),('America/Chicago','-21600','-18000',94),('America/Chihuahua','-25200','-21600',95),('America/Coral_Harbour','-18000','-18000',96),('America/Cordoba','-10800','-10800',97),('America/Costa_Rica','-21600','-21600',98),('America/Creston','-25200','-25200',99),('America/Cuiaba','-14400','-10800',100),('America/Curacao','-14400','-14400',101),('America/Danmarkshavn','0','0',102),('America/Dawson','-28800','-25200',103),('America/Dawson_Creek','-25200','-25200',104),('America/Denver','-25200','-21600',105),('America/Detroit','-18000','-14400',106),('America/Dominica','-14400','-14400',107),('America/Edmonton','-25200','-21600',108),('America/Eirunepe','-18000','-18000',109),('America/El_Salvador','-21600','-21600',110),('America/Ensenada','-28800','-25200',111),('America/Fort_Nelson','-25200','-25200',112),('America/Fort_Wayne','-18000','-14400',113),('America/Fortaleza','-10800','-10800',114),('America/Glace_Bay','-14400','-10800',115),('America/Godthab','-10800','-7200',116),('America/Goose_Bay','-14400','-10800',117),('America/Grand_Turk','-14400','-14400',118),('America/Grenada','-14400','-14400',119),('America/Guadeloupe','-14400','-14400',120),('America/Guatemala','-21600','-21600',121),('America/Guayaquil','-18000','-18000',122),('America/Guyana','-14400','-14400',123),('America/Halifax','-14400','-10800',124),('America/Havana','-18000','-14400',125),('America/Hermosillo','-25200','-25200',126),('America/Indiana/Indianapolis','-18000','-14400',127),('America/Indiana/Knox','-21600','-18000',128),('America/Indiana/Marengo','-18000','-14400',129),('America/Indiana/Petersburg','-18000','-14400',130),('America/Indiana/Tell_City','-21600','-18000',131),('America/Indiana/Vevay','-18000','-14400',132),('America/Indiana/Vincennes','-18000','-14400',133),('America/Indiana/Winamac','-18000','-14400',134),('America/Indianapolis','-18000','-14400',135),('America/Inuvik','-25200','-21600',136),('America/Iqaluit','-18000','-14400',137),('America/Jamaica','-18000','-18000',138),('America/Jujuy','-10800','-10800',139),('America/Juneau','-32400','-28800',140),('America/Kentucky/Louisville','-18000','-14400',141),('America/Kentucky/Monticello','-18000','-14400',142),('America/Knox_IN','-21600','-18000',143),('America/Kralendijk','-14400','-14400',144),('America/La_Paz','-14400','-14400',145),('America/Lima','-18000','-18000',146),('America/Los_Angeles','-28800','-25200',147),('America/Louisville','-18000','-14400',148),('America/Lower_Princes','-14400','-14400',149),('America/Maceio','-10800','-10800',150),('America/Managua','-21600','-21600',151),('America/Manaus','-14400','-14400',152),('America/Marigot','-14400','-14400',153),('America/Martinique','-14400','-14400',154),('America/Matamoros','-21600','-18000',155),('America/Mazatlan','-25200','-21600',156),('America/Mendoza','-10800','-10800',157),('America/Menominee','-21600','-18000',158),('America/Merida','-21600','-18000',159),('America/Metlakatla','-28800','-28800',160),('America/Mexico_City','-21600','-18000',161),('America/Miquelon','-10800','-7200',162),('America/Moncton','-14400','-10800',163),('America/Monterrey','-21600','-18000',164),('America/Montevideo','-10800','-7200',165),('America/Montreal','-18000','-14400',166),('America/Montserrat','-14400','-14400',167),('America/Nassau','-18000','-14400',168),('America/New_York','-18000','-14400',169),('America/Nipigon','-18000','-14400',170),('America/Nome','-32400','-28800',171),('America/Noronha','-7200','-7200',172),('America/North_Dakota/Beulah','-21600','-18000',173),('America/North_Dakota/Center','-21600','-18000',174),('America/North_Dakota/New_Salem','-21600','-18000',175),('America/Ojinaga','-25200','-21600',176),('America/Panama','-18000','-18000',177),('America/Pangnirtung','-18000','-14400',178),('America/Paramaribo','-10800','-10800',179),('America/Phoenix','-25200','-25200',180),('America/Port-au-Prince','-18000','-14400',181),('America/Port_of_Spain','-14400','-14400',182),('America/Porto_Acre','-18000','-18000',183),('America/Porto_Velho','-14400','-14400',184),('America/Puerto_Rico','-14400','-14400',185),('America/Rainy_River','-21600','-18000',186),('America/Rankin_Inlet','-21600','-18000',187),('America/Recife','-10800','-10800',188),('America/Regina','-21600','-21600',189),('America/Resolute','-21600','-18000',190),('America/Rio_Branco','-18000','-18000',191),('America/Rosario','-10800','-10800',192),('America/Santa_Isabel','-28800','-25200',193),('America/Santarem','-10800','-10800',194),('America/Santiago','-10800','-10800',195),('America/Santo_Domingo','-14400','-14400',196),('America/Sao_Paulo','-10800','-7200',197),('America/Scoresbysund','-3600','0',198),('America/Shiprock','-25200','-21600',199),('America/Sitka','-32400','-28800',200),('America/St_Barthelemy','-14400','-14400',201),('America/St_Johns','-12600','-9000',202),('America/St_Kitts','-14400','-14400',203),('America/St_Lucia','-14400','-14400',204),('America/St_Thomas','-14400','-14400',205),('America/St_Vincent','-14400','-14400',206),('America/Swift_Current','-21600','-21600',207),('America/Tegucigalpa','-21600','-21600',208),('America/Thule','-14400','-10800',209),('America/Thunder_Bay','-18000','-14400',210),('America/Tijuana','-28800','-25200',211),('America/Toronto','-18000','-14400',212),('America/Tortola','-14400','-14400',213),('America/Vancouver','-28800','-25200',214),('America/Virgin','-14400','-14400',215),('America/Whitehorse','-28800','-25200',216),('America/Winnipeg','-21600','-18000',217),('America/Yakutat','-32400','-28800',218),('America/Yellowknife','-25200','-21600',219),('Antarctica/Casey','28800','28800',220),('Antarctica/Davis','25200','25200',221),('Antarctica/DumontDUrville','36000','36000',222),('Antarctica/Macquarie','39600','39600',223),('Antarctica/Mawson','18000','18000',224),('Antarctica/McMurdo','43200','46800',225),('Antarctica/Palmer','-10800','-10800',226),('Antarctica/Rothera','-10800','-10800',227),('Antarctica/South_Pole','43200','46800',228),('Antarctica/Syowa','10800','10800',229),('Antarctica/Troll','0','7200',230),('Antarctica/Vostok','21600','21600',231),('Arctic/Longyearbyen','3600','7200',232),('Asia/Aden','10800','10800',233),('Asia/Almaty','21600','21600',234),('Asia/Amman','7200','10800',235),('Asia/Anadyr','43200','43200',236),('Asia/Aqtau','18000','18000',237),('Asia/Aqtobe','18000','18000',238),('Asia/Ashgabat','18000','18000',239),('Asia/Ashkhabad','18000','18000',240),('Asia/Baghdad','10800','10800',241),('Asia/Bahrain','10800','10800',242),('Asia/Baku','14400','18000',243),('Asia/Bangkok','25200','25200',244),('Asia/Beirut','7200','10800',245),('Asia/Bishkek','21600','21600',246),('Asia/Brunei','28800','28800',247),('Asia/Calcutta','19800','19800',248),('Asia/Chita','28800','28800',249),('Asia/Choibalsan','28800','32400',250),('Asia/Chongqing','28800','28800',251),('Asia/Chungking','28800','28800',252),('Asia/Colombo','19800','19800',253),('Asia/Dacca','21600','21600',254),('Asia/Damascus','7200','10800',255),('Asia/Dhaka','21600','21600',256),('Asia/Dili','32400','32400',257),('Asia/Dubai','14400','14400',258),('Asia/Dushanbe','18000','18000',259),('Asia/Gaza','7200','10800',260),('Asia/Harbin','28800','28800',261),('Asia/Hebron','7200','10800',262),('Asia/Ho_Chi_Minh','25200','25200',263),('Asia/Hong_Kong','28800','28800',264),('Asia/Hovd','25200','28800',265),('Asia/Irkutsk','28800','28800',266),('Asia/Istanbul','7200','10800',267),('Asia/Jakarta','25200','25200',268),('Asia/Jayapura','32400','32400',269),('Asia/Jerusalem','7200','10800',270),('Asia/Kabul','16200','16200',271),('Asia/Kamchatka','43200','43200',272),('Asia/Karachi','18000','18000',273),('Asia/Kashgar','21600','21600',274),('Asia/Kathmandu','20700','20700',275),('Asia/Katmandu','20700','20700',276),('Asia/Khandyga','32400','32400',277),('Asia/Kolkata','19800','19800',278),('Asia/Krasnoyarsk','25200','25200',279),('Asia/Kuala_Lumpur','28800','28800',280),('Asia/Kuching','28800','28800',281),('Asia/Kuwait','10800','10800',282),('Asia/Macao','28800','28800',283),('Asia/Macau','28800','28800',284),('Asia/Magadan','36000','36000',285),('Asia/Makassar','28800','28800',286),('Asia/Manila','28800','28800',287),('Asia/Muscat','14400','14400',288),('Asia/Nicosia','7200','10800',289),('Asia/Novokuznetsk','25200','25200',290),('Asia/Novosibirsk','21600','21600',291),('Asia/Omsk','21600','21600',292),('Asia/Oral','18000','18000',293),('Asia/Phnom_Penh','25200','25200',294),('Asia/Pontianak','25200','25200',295),('Asia/Pyongyang','30600','30600',296),('Asia/Qatar','10800','10800',297),('Asia/Qyzylorda','21600','21600',298),('Asia/Rangoon','23400','23400',299),('Asia/Riyadh','10800','10800',300),('Asia/Riyadh87','10800','10800',301),('Asia/Riyadh88','10800','10800',302),('Asia/Riyadh89','10800','10800',303),('Asia/Saigon','25200','25200',304),('Asia/Sakhalin','36000','36000',305),('Asia/Samarkand','18000','18000',306),('Asia/Seoul','32400','32400',307),('Asia/Shanghai','28800','28800',308),('Asia/Singapore','28800','28800',309),('Asia/Srednekolymsk','39600','39600',310),('Asia/Taipei','28800','28800',311),('Asia/Tashkent','18000','18000',312),('Asia/Tbilisi','14400','14400',313),('Asia/Tehran','12600','16200',314),('Asia/Tel_Aviv','7200','10800',315),('Asia/Thimbu','21600','21600',316),('Asia/Thimphu','21600','21600',317),('Asia/Tokyo','32400','32400',318),('Asia/Ujung_Pandang','28800','28800',319),('Asia/Ulaanbaatar','28800','32400',320),('Asia/Ulan_Bator','28800','32400',321),('Asia/Urumqi','21600','21600',322),('Asia/Vientiane','25200','25200',323),('Asia/Vladivostok','36000','36000',324),('Asia/Yakutsk','32400','32400',325),('Asia/Yekaterinburg','18000','18000',326),('Asia/Yerevan','14400','14400',327),('Atlantic/Azores','-3600','0',328),('Atlantic/Bermuda','-14400','-10800',329),('Atlantic/Canary','0','3600',330),('Atlantic/Cape_Verde','-3600','-3600',331),('Atlantic/Faeroe','0','3600',332),('Atlantic/Faroe','0','3600',333),('Atlantic/Jan_Mayen','3600','7200',334),('Atlantic/Madeira','0','3600',335),('Atlantic/Reykjavik','0','0',336),('Atlantic/South_Georgia','-7200','-7200',337),('Atlantic/St_Helena','0','0',338),('Atlantic/Stanley','-10800','-10800',339),('Australia/ACT','36000','39600',340),('Australia/Adelaide','34200','37800',341),('Australia/Brisbane','36000','36000',342),('Australia/Broken_Hill','34200','37800',343),('Australia/Canberra','36000','39600',344),('Australia/Currie','36000','39600',345),('Australia/Darwin','34200','34200',346),('Australia/Eucla','31500','31500',347),('Australia/Hobart','36000','39600',348),('Australia/LHI','37800','39600',349),('Australia/Lindeman','36000','36000',350),('Australia/Lord_Howe','37800','39600',351),('Australia/Melbourne','36000','39600',352),('Australia/NSW','36000','39600',353),('Australia/North','34200','34200',354),('Australia/Perth','28800','28800',355),('Australia/Queensland','36000','36000',356),('Australia/South','34200','37800',357),('Australia/Sydney','36000','39600',358),('Australia/Tasmania','36000','39600',359),('Australia/Victoria','36000','39600',360),('Australia/West','28800','28800',361),('Australia/Yancowinna','34200','37800',362),('Brazil/Acre','-18000','-18000',363),('Brazil/DeNoronha','-7200','-7200',364),('Brazil/East','-10800','-7200',365),('Brazil/West','-14400','-14400',366),('CET','3600','7200',367),('CST6CDT','-21600','-18000',368),('Canada/Atlantic','-14400','-10800',369),('Canada/Central','-21600','-18000',370),('Canada/East-Saskatchewan','-21600','-21600',371),('Canada/Eastern','-18000','-14400',372),('Canada/Mountain','-25200','-21600',373),('Canada/Newfoundland','-12600','-9000',374),('Canada/Pacific','-28800','-25200',375),('Canada/Saskatchewan','-21600','-21600',376),('Canada/Yukon','-28800','-25200',377),('Chile/Continental','-10800','-10800',378),('Chile/EasterIsland','-18000','-18000',379),('Cuba','-18000','-14400',380),('EET','7200','10800',381),('EST','-18000','-18000',382),('EST5EDT','-18000','-14400',383),('Egypt','7200','7200',384),('Eire','0','3600',385),('Etc/GMT','0','0',386),('Etc/GMT+0','0','0',387),('Etc/GMT+1','-3600','-3600',388),('Etc/GMT+10','-36000','-36000',389),('Etc/GMT+11','-39600','-39600',390),('Etc/GMT+12','-43200','-43200',391),('Etc/GMT+2','-7200','-7200',392),('Etc/GMT+3','-10800','-10800',393),('Etc/GMT+4','-14400','-14400',394),('Etc/GMT+5','-18000','-18000',395),('Etc/GMT+6','-21600','-21600',396),('Etc/GMT+7','-25200','-25200',397),('Etc/GMT+8','-28800','-28800',398),('Etc/GMT+9','-32400','-32400',399),('Etc/GMT-0','0','0',400),('Etc/GMT-1','3600','3600',401),('Etc/GMT-10','36000','36000',402),('Etc/GMT-11','39600','39600',403),('Etc/GMT-12','43200','43200',404),('Etc/GMT-13','46800','46800',405),('Etc/GMT-14','50400','50400',406),('Etc/GMT-2','7200','7200',407),('Etc/GMT-3','10800','10800',408),('Etc/GMT-4','14400','14400',409),('Etc/GMT-5','18000','18000',410),('Etc/GMT-6','21600','21600',411),('Etc/GMT-7','25200','25200',412),('Etc/GMT-8','28800','28800',413),('Etc/GMT-9','32400','32400',414),('Etc/GMT0','0','0',415),('Etc/Greenwich','0','0',416),('Etc/UCT','0','0',417),('Etc/UTC','0','0',418),('Etc/Universal','0','0',419),('Etc/Zulu','0','0',420),('Europe/Amsterdam','3600','7200',421),('Europe/Andorra','3600','7200',422),('Europe/Athens','7200','10800',423),('Europe/Belfast','0','3600',424),('Europe/Belgrade','3600','7200',425),('Europe/Berlin','3600','7200',426),('Europe/Bratislava','3600','7200',427),('Europe/Brussels','3600','7200',428),('Europe/Bucharest','7200','10800',429),('Europe/Budapest','3600','7200',430),('Europe/Busingen','3600','7200',431),('Europe/Chisinau','7200','10800',432),('Europe/Copenhagen','3600','7200',433),('Europe/Dublin','0','3600',434),('Europe/Gibraltar','3600','7200',435),('Europe/Guernsey','0','3600',436),('Europe/Helsinki','7200','10800',437),('Europe/Isle_of_Man','0','3600',438),('Europe/Istanbul','7200','10800',439),('Europe/Jersey','0','3600',440),('Europe/Kaliningrad','7200','10800',441),('Europe/Kiev','7200','10800',442),('Europe/Lisbon','0','3600',443),('Europe/Ljubljana','3600','7200',444),('Europe/London','0','3600',445),('Europe/Luxembourg','3600','7200',446),('Europe/Madrid','3600','7200',447),('Europe/Malta','3600','7200',448),('Europe/Mariehamn','7200','10800',449),('Europe/Minsk','10800','10800',450),('Europe/Monaco','3600','7200',451),('Europe/Moscow','10800','10800',452),('Europe/Nicosia','7200','10800',453),('Europe/Oslo','3600','7200',454),('Europe/Paris','3600','7200',455),('Europe/Podgorica','3600','7200',456),('Europe/Prague','3600','7200',457),('Europe/Riga','7200','10800',458),('Europe/Rome','3600','7200',459),('Europe/Samara','14400','14400',460),('Europe/San_Marino','3600','7200',461),('Europe/Sarajevo','3600','7200',462),('Europe/Simferopol','7200','10800',463),('Europe/Skopje','3600','7200',464),('Europe/Sofia','7200','10800',465),('Europe/Stockholm','3600','7200',466),('Europe/Tallinn','7200','10800',467),('Europe/Tirane','3600','7200',468),('Europe/Tiraspol','7200','10800',469),('Europe/Uzhgorod','7200','10800',470),('Europe/Vaduz','3600','7200',471),('Europe/Vatican','3600','7200',472),('Europe/Vienna','3600','7200',473),('Europe/Vilnius','7200','10800',474),('Europe/Volgograd','10800','10800',475),('Europe/Warsaw','3600','7200',476),('Europe/Zagreb','3600','7200',477),('Europe/Zaporozhye','7200','10800',478),('Europe/Zurich','3600','7200',479),('Factory','0','0',480),('GB','0','3600',481),('GB-Eire','0','3600',482),('GMT','0','0',483),('GMT+0','0','0',484),('GMT-0','0','0',485),('GMT0','0','0',486),('Greenwich','0','0',487),('HST','-36000','-36000',488),('Hongkong','28800','28800',489),('Iceland','0','0',490),('Indian/Antananarivo','10800','10800',491),('Indian/Chagos','21600','21600',492),('Indian/Christmas','25200','25200',493),('Indian/Cocos','23400','23400',494),('Indian/Comoro','10800','10800',495),('Indian/Kerguelen','18000','18000',496),('Indian/Mahe','14400','14400',497),('Indian/Maldives','18000','18000',498),('Indian/Mauritius','14400','14400',499),('Indian/Mayotte','10800','10800',500),('Indian/Reunion','14400','14400',501),('Iran','12600','16200',502),('Israel','7200','10800',503),('Jamaica','-18000','-18000',504),('Japan','32400','32400',505),('Kwajalein','43200','43200',506),('Libya','7200','7200',507),('MET','3600','7200',508),('MST','-25200','-25200',509),('MST7MDT','-25200','-21600',510),('Mexico/BajaNorte','-28800','-25200',511),('Mexico/BajaSur','-25200','-21600',512),('Mexico/General','-21600','-18000',513),('Mideast/Riyadh87','0','0',514),('Mideast/Riyadh88','0','0',515),('Mideast/Riyadh89','0','0',516),('NZ','43200','46800',517),('NZ-CHAT','45900','49500',518),('Navajo','0','0',519),('PRC','0','0',520),('PST8PDT','0','0',521),('Pacific/Apia','46800','50400',522),('Pacific/Auckland','43200','46800',523),('Pacific/Bougainville','39600','39600',524),('Pacific/Chatham','45900','49500',525),('Pacific/Chuuk','36000','36000',526),('Pacific/Easter','-18000','-18000',527),('Pacific/Efate','39600','39600',528),('Pacific/Enderbury','46800','46800',529),('Pacific/Fakaofo','46800','46800',530),('Pacific/Fiji','43200','46800',531),('Pacific/Funafuti','43200','43200',532),('Pacific/Galapagos','-21600','-21600',533),('Pacific/Gambier','-32400','-32400',534),('Pacific/Guadalcanal','39600','39600',535),('Pacific/Guam','36000','36000',536),('Pacific/Honolulu','-36000','-36000',537),('Pacific/Johnston','-36000','-36000',538),('Pacific/Kiritimati','50400','50400',539),('Pacific/Kosrae','39600','39600',540),('Pacific/Kwajalein','43200','43200',541),('Pacific/Majuro','43200','43200',542),('Pacific/Marquesas','-34200','-34200',543),('Pacific/Midway','-39600','-39600',544),('Pacific/Nauru','43200','43200',545),('Pacific/Niue','-39600','-39600',546),('Pacific/Norfolk','39600','39600',547),('Pacific/Noumea','39600','39600',548),('Pacific/Pago_Pago','-39600','-39600',549),('Pacific/Palau','32400','32400',550),('Pacific/Pitcairn','-28800','-28800',551),('Pacific/Ponape','39600','39600',552),('Pacific/Port_Moresby','36000','36000',553),('Pacific/Rarotonga','-36000','-36000',554),('Pacific/Saipan','36000','36000',555),('Pacific/Samoa','-39600','-39600',556),('Pacific/Tahiti','-36000','-36000',557),('Pacific/Tarawa','43200','43200',558),('Pacific/Tongatapu','46800','46800',559),('Pacific/Truk','36000','36000',560),('Pacific/Wake','43200','43200',561),('Pacific/Wallis','43200','43200',562),('Pacific/Yap','36000','36000',563),('Poland','3600','7200',564),('Portugal','0','3600',565),('ROC','28800','28800',566),('ROK','32400','32400',567),('Singapore','28800','28800',568),('Turkey','7200','10800',569),('UCT','0','0',570),('US/Alaska','-32400','-28800',571),('US/Aleutian','-36000','-32400',572),('US/Arizona','-25200','-25200',573),('US/Central','-21600','-18000',574),('US/East-Indiana','-18000','-14400',575),('US/Eastern','-18000','-14400',576),('US/Hawaii','-36000','-36000',577),('US/Indiana-Starke','-21600','-18000',578),('US/Michigan','-18000','-14400',579),('US/Mountain','-25200','-21600',580),('US/Pacific','-28800','-25200',581),('US/Pacific-New','-28800','-25200',582),('US/Samoa','-39600','-39600',583),('UTC','0','0',584),('Universal','0','0',585),('W-SU','10800','10800',586),('WET','0','3600',587),('Zulu','0','0',588);
/*!40000 ALTER TABLE `timezone` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Dumping events for database 'db_main_cloud'
--

--
-- Dumping routines for database 'db_main_cloud'
--
/*!40103 SET TIME_ZONE=@OLD_TIME_ZONE */;

/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40014 SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;

-- Dump completed on 2023-07-29 10:52:26
