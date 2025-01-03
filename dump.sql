-- MySQL dump 10.13  Distrib 8.2.0, for Linux (x86_64)
--
-- Host: localhost    Database: club
-- ------------------------------------------------------
-- Server version	8.2.0

/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!50503 SET NAMES utf8mb4 */;
/*!40103 SET @OLD_TIME_ZONE=@@TIME_ZONE */;
/*!40103 SET TIME_ZONE='+00:00' */;
/*!40014 SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0 */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;

--
-- Table structure for table `clubs`
--

DROP TABLE IF EXISTS `clubs`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `clubs` (
  `club_id` varchar(5) NOT NULL,
  `club` varchar(255) DEFAULT NULL,
  `Info` text,
  `Img_link` varchar(255) DEFAULT NULL,
  `email` varchar(255) DEFAULT NULL,
  `unique_id` varchar(10) DEFAULT NULL,
  PRIMARY KEY (`club_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `clubs`
--

LOCK TABLES `clubs` WRITE;
/*!40000 ALTER TABLE `clubs` DISABLE KEYS */;
INSERT INTO `clubs` VALUES ('CL01','Robotics Society','To be Updated','https://images.unsplash.com/photo-1485827404703-89b55fcc595e?q=80&w=2070&auto=format&fit=crop&ixlib=rb-4.0.3&ixid=M3wxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8fA%3D%3D','aich.1@iitj.ac.in','RoCL01'),('CL02','Sangam','Sangam is the music society. Join your musical journey with sangam.','https://images.unsplash.com/photo-1511379938547-c1f69419868d?q=80&w=2070&auto=format&fit=crop&ixlib=rb-4.0.3&ixid=M3wxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8fA%3D%3D','aichrishav2003@gmail.com','SaCL02'),('CL03','Nexus','To be Updated','https://images.unsplash.com/photo-1608178398319-48f814d0750c?q=80&w=1779&auto=format&fit=crop&ixlib=rb-4.0.3&ixid=M3wxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8fA%3D%3D',NULL,'NeCL03'),('CL04','FrameX','To be Updated','https://images.unsplash.com/photo-1502209877429-d7c6df9eb3f9?q=80&w=2066&auto=format&fit=crop&ixlib=rb-4.0.3&ixid=M3wxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8fA%3D%3D',NULL,'FaCL04'),('CL05','Shutterbugs','To be Updated','https://images.unsplash.com/photo-1502982720700-bfff97f2ecac?q=80&w=2070&auto=format&fit=crop&ixlib=rb-4.0.3&ixid=M3wxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8fA%3D%3D',NULL,'ShCL05'),('CL06','Arteliers','To be Updated','https://plus.unsplash.com/premium_photo-1673514503542-3953ca855a03?q=80&w=2070&auto=format&fit=crop&ixlib=rb-4.0.3&ixid=M3wxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8fA%3D%3D',NULL,'ArCL06');
/*!40000 ALTER TABLE `clubs` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `items`
--

DROP TABLE IF EXISTS `items`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `items` (
  `item_id` varchar(50) NOT NULL,
  `item` varchar(500) DEFAULT NULL,
  `club_id` varchar(10) DEFAULT NULL,
  `quantity` int DEFAULT NULL,
  PRIMARY KEY (`item_id`),
  KEY `club_id` (`club_id`),
  CONSTRAINT `items_ibfk_1` FOREIGN KEY (`club_id`) REFERENCES `clubs` (`club_id`) ON DELETE RESTRICT ON UPDATE RESTRICT
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `items`
--

LOCK TABLES `items` WRITE;
/*!40000 ALTER TABLE `items` DISABLE KEYS */;
INSERT INTO `items` VALUES ('IT01','Motor Driver','CL01',5),('IT02','Wheels','CL01',0),('IT03','Telescope','CL03',1),('IT05','Camera','CL04',1),('IT06','Drum','CL02',2),('IT07','Flute','CL02',24),('IT08','Mic','CL02',49),('IT09','Piano','CL02',37),('IT10','Tabla','CL02',66),('IT12','Temperature sensor','CL01',3),('IT13','Saxophone','CL02',31),('IT14','Mic stand','CL02',100);
/*!40000 ALTER TABLE `items` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `student`
--

DROP TABLE IF EXISTS `student`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `student` (
  `username` varchar(50) NOT NULL,
  `name` varchar(100) DEFAULT NULL,
  `password` varchar(100) DEFAULT NULL,
  `Institute_id` varchar(50) DEFAULT NULL,
  PRIMARY KEY (`username`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `student`
--

LOCK TABLES `student` WRITE;
/*!40000 ALTER TABLE `student` DISABLE KEYS */;
INSERT INTO `student` VALUES ('aich.1','Rishav','$2a$08$tCpMHY5yEiQaitEC/fYfCuEmscGtJG.e2Gsb0GkKivQjoXQIEXHbK','B70CS499'),('ak','anku','$2a$08$aIf2NLxnU/s9Si0BnxUHluGgGfRtM5OWrqz3w/f75Ebo6aP.DJaPy','B37CS395'),('ank','ankit','$2a$08$XaZJIQhWCWODfty7D0Bk8ejFOoxTHDYkO9cpZNSbYXwFx9dz3xqyu','B84CS021'),('ar','Arun','$2a$08$MvIqmDDS0ZOi53IRhEqCvexWcp.1RJdjP8AkkMlw6XHFKFMu//R5W','B21CS098'),('as','Ashu','$2a$08$eMMlA/97xO7gvqVrVei5bO1pAO8PpcPPctVsY352TrQFTGnfwoLYS','B21AI007'),('cv','titu','$2a$08$69hUM5WEDZZfcj0yevT.Ge28fXCNPjokvebJ.bmqzqWXDR1xj.Pn.','b22ru78'),('de','deepak','$2a$08$97xGBmN1.jD/UwO/mEa37u6iMfaSgb7iJB.4WrG1KBr8EC852ovEe','B21CH089'),('git','Anshu','$2a$08$lhFV1gmDM/1TAfghhZF49epCiFOkarvCL6h9EoeVcnvO/zmmABCAa','B58CS844'),('Ind','Indu','$2a$08$MEm0qGpIfCA2lM4/Bv5tyuRnAVKEhXaqFbyh42FX3.IIsey5oZy22','B21UI098'),('j','jiju','$2a$08$hTldo0prIG4qlaDOmvh5PunTRmKUtGvA/VMlvOd3/VYdo32CGFLeK','B47CS849'),('joint','kamble','$2a$08$ZKE9k5mECEkbR5q8EnL5TOhumDBbtNNy6.SiDiIZ0wRcQpB0xw/8.','B21CI021'),('ku','kiku','$2a$08$pJ5zFBRcxh/MQgQWqPzz7exWBJAE1WQ3N/.49ollB36//tWkclEkm','B24RT09'),('ml','l','$2a$08$0a2uPTUNSH3Q9ckbGeIDf.4BrqbY2rVwfj9V.fHsZ.ikyaidXnQPK','b45yu89'),('moni','nino','$2a$08$AcGoKjsVoxFqx2ot4Rm9TujaIEMfp1DcPt93Z1o5ksax6rKt80.62','B81CS545'),('nu','tinu','$2a$08$iwJKKUPs7UXG1ClMciodZ.UED5smnHimmN2jZIKoyT/yrot6LZcw.','B23CS098'),('old','queen','$2a$08$cargY8vObrrLEFhvku/LDu5iFIt4GPEq9/MC5Qb1z9HISPwoEbumG','b21ai098'),('Pu','Puneet','$2a$08$0tLUe5EgTwGZSpGt3eyCauVZxeTNL/UBgfVG/ZT1WOWMNzc1ZHzVO','B21AI009'),('qt','ritu','$2a$08$N0NoKoIAkxEoxDSCruKRsOzwOrR1ugbLnpwWwUZrUQPcbux644OHC','B27CS716'),('r','r','$2a$08$h9feLsYwPLf0w7VTaP9W9O5x1rsqUo0Q.FmLcnY/J4oYxNvZ2DBcy','B77CS712'),('ridu','Rishav Aich','$2a$08$yHldRsTmqxnOfjJPRCUEjO4aaF9ugxjyrkWzSqi3vu4YCfR/CDJdC','B24CS082'),('ritu','Ashish','$2a$08$qMhsRwEKCvYHBVvPJHvfKug5FNg0oWjCzXJNTIjP0RNx.uHzPEi2y','B68CS155'),('sun','arjun','$2a$08$qu1xp2gKdB1Bp32UIev6weu/Np509FwvBdHyml2TwekQrIagbb4cW','B73CS210'),('tinu','tinu','$2a$08$5ExDYArsvb68c4O5etWInOsKtjGEtGdN0lpftG8NxIedgFUR/q2E.','B23CS019'),('tk','titu','$2a$08$WPXB3NDJUbcBZW4BLbUkcOMruz4pm5cyvgRFoHgmcPjvOXRrNq5mK','B84CS605'),('u','yuji','$2a$08$IApjFw/d7GTPoUQXeL8P6ecRTkcBZKqQlX57ZxjdsQCmIuIiQ6qXa','B48CS608'),('un','sunil','$2a$08$ulZjeq5j6kcm7Ri.AmyE/.z2.7ExkmpSL6JjNQgNQWQrq2cqbVv5y','B58CS111');
/*!40000 ALTER TABLE `student` ENABLE KEYS */;
UNLOCK TABLES;
/*!40103 SET TIME_ZONE=@OLD_TIME_ZONE */;

/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40014 SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;

-- Dump completed on 2024-12-25  6:55:49
