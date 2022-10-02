-- Adminer 4.8.1 MySQL 8.0.30 dump

SET NAMES utf8;
SET time_zone = '+00:00';
SET foreign_key_checks = 0;
SET sql_mode = 'NO_AUTO_VALUE_ON_ZERO';

SET NAMES utf8mb4;

DROP TABLE IF EXISTS `answer`;
CREATE TABLE `answer` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `questionId` int unsigned NOT NULL,
  `text` varchar(1000) NOT NULL,
  PRIMARY KEY (`id`),
  KEY `questionId` (`questionId`),
  CONSTRAINT `answer_ibfk_1` FOREIGN KEY (`questionId`) REFERENCES `question` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;


DROP TABLE IF EXISTS `answer_like`;
CREATE TABLE `answer_like` (
  `id` int NOT NULL AUTO_INCREMENT,
  `answerId` int unsigned NOT NULL,
  `like_author` int unsigned NOT NULL,
  PRIMARY KEY (`id`),
  KEY `answerId` (`answerId`),
  KEY `like_author` (`like_author`),
  CONSTRAINT `answer_like_ibfk_1` FOREIGN KEY (`answerId`) REFERENCES `answer` (`id`),
  CONSTRAINT `answer_like_ibfk_2` FOREIGN KEY (`like_author`) REFERENCES `user` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;


DROP TABLE IF EXISTS `question`;
CREATE TABLE `question` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `authorId` int unsigned NOT NULL,
  `authorIpAddress` varchar(45) NOT NULL,
  `text` varchar(500) NOT NULL,
  PRIMARY KEY (`id`),
  KEY `authorId` (`authorId`),
  CONSTRAINT `question_ibfk_1` FOREIGN KEY (`authorId`) REFERENCES `user` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;


DROP TABLE IF EXISTS `user`;
CREATE TABLE `user` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `username` varchar(30) NOT NULL,
  `email` varchar(319) NOT NULL,
  `display_name` varchar(30) NOT NULL,
  `password` char(60) NOT NULL,
  `birthdate` date NOT NULL,
  `creation_date` datetime DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;


-- 2022-10-02 16:12:21