CREATE DATABASE IF NOT EXISTS boilerplate_go;

USE boilerplate_go;


CREATE TABLE IF NOT EXISTS  `medicines` (
                             `id` int(11) NOT NULL AUTO_INCREMENT,
                             `name` varchar(30) CHARACTER SET utf8 DEFAULT NULL,
                             `ean_code` varchar(30) CHARACTER SET utf8 DEFAULT NULL,
                             `description` varchar(150) CHARACTER SET utf8 DEFAULT NULL,
                             `laboratory` varchar(50) CHARACTER SET utf8 DEFAULT NULL,
                             `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
                             `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
                             PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=10873 DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;

FLUSH PRIVILEGES;
