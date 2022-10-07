-- Adminer 4.8.1 MySQL 8.0.30 dump

SET NAMES utf8;
SET time_zone = '+00:00';
SET foreign_key_checks = 0;
SET sql_mode = 'NO_AUTO_VALUE_ON_ZERO';

SET NAMES utf8mb4;

CREATE DATABASE `project_truthful` /*!40100 DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci */ /*!80016 DEFAULT ENCRYPTION='N' */;
USE `project_truthful`;

DROP TABLE IF EXISTS `user`;
CREATE TABLE `user` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `username` varchar(30) NOT NULL,
  `email` varchar(319) NOT NULL,
  `display_name` varchar(30) NOT NULL,
  `password` char(60) NOT NULL,
  `birthdate` date NOT NULL,
  `creation_date` timestamp DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

DROP TABLE IF EXISTS `question`;
CREATE TABLE `question` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `receiver_id` int unsigned NOT NULL,
  `author_id` int unsigned NOT NULL,
  `author_ip_address` varchar(45) NOT NULL,
  `text` varchar(500) NOT NULL,
  PRIMARY KEY (`id`),
  KEY `author_id` (`author_id`),
  CONSTRAINT `question_ibfk_1` FOREIGN KEY (`author_id`) REFERENCES `user` (`id`),
  KEY `receiver_id` (`receiver_id`),
  CONSTRAINT `question_ibfk_2` FOREIGN KEY (`receiver_id`) REFERENCES `user` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

DROP TABLE IF EXISTS `answer`;
CREATE TABLE `answer` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `user_id` int unsigned NOT NULL,
  `question_id` int unsigned NOT NULL,
  `text` varchar(1000) NOT NULL,
  `time_created` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `question_id` (`question_id`),
  CONSTRAINT `answer_ibfk_1` FOREIGN KEY (`question_id`) REFERENCES `question` (`id`),
  KEY `user_id` (`user_id`),
  CONSTRAINT `answer_ibfk_2` FOREIGN KEY (`user_id`) REFERENCES `user` (`id`)
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


DROP TABLE IF EXISTS `follow`;
CREATE TABLE `follow` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `follower` int unsigned NOT NULL,
  `followed` int unsigned NOT NULL,
  PRIMARY KEY (`id`),
  KEY `follower` (`follower`),
  KEY `followed` (`followed`),
  CONSTRAINT `follow_ibfk_1` FOREIGN KEY (`follower`) REFERENCES `user` (`id`),
  CONSTRAINT `follow_ibfk_2` FOREIGN KEY (`followed`) REFERENCES `user` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;


-- 2022-10-02 16:12:21