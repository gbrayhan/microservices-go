CREATE DATABASE IF NOT EXISTS boilerplate_go;

USE boilerplate_go;
CREATE USER IF NOT EXISTS 'appuser'@'localhost' IDENTIFIED WITH mysql_native_password BY 'youShouldChangeThisPassword';
CREATE USER IF NOT EXISTS 'appuser'@'%' IDENTIFIED WITH mysql_native_password BY 'youShouldChangeThisPassword';
CREATE USER IF NOT EXISTS 'appuser'@'mysqldb' IDENTIFIED WITH mysql_native_password BY 'youShouldChangeThisPassword';


create table if not exists medicines
(
    id          bigint auto_increment primary key,
    name        varchar(30) charset utf8mb4         null,
    ean_code    varchar(30) charset utf8mb4         null,
    description varchar(150) charset utf8mb4        null,
    laboratory  varchar(50) charset utf8mb4         null,
    created_at  timestamp default CURRENT_TIMESTAMP null,
    updated_at  timestamp default CURRENT_TIMESTAMP null on update CURRENT_TIMESTAMP,
    constraint medicines_UN_ean unique (ean_code),
    constraint medicines_UN_name unique (name)
) CHARSET = utf8mb4
  COLLATE = utf8mb4_general_ci;

create table if not exists users
(
    id            bigint auto_increment,
    user_name     varchar(100) charset utf8mb4         not null,
    email         varchar(100) charset utf8mb4         not null,
    first_name    varchar(100) charset utf8mb4         null,
    last_name     varchar(100) charset utf8mb4         null,
    status        tinyint(1) default 1                 null,
    hash_password varchar(255) charset utf8mb4         not null,
    created_at    timestamp  default CURRENT_TIMESTAMP not null on update CURRENT_TIMESTAMP,
    updated_at    timestamp  default CURRENT_TIMESTAMP not null on update CURRENT_TIMESTAMP,
    constraint users_UN_email unique (email),
    constraint users_UN_user unique (user_name),
    constraint users_id_IDX unique (id)
) CHARSET = utf8mb4
  COLLATE = utf8mb4_general_ci;


-- TODO: you should consider to change this user and password ons your production environment
INSERT INTO boilerplate_go.users (user_name,email,first_name,last_name,status,hash_password,created_at,updated_at) VALUES
    ('gbrayhan','gbrayhan@gmail.com','Alejandro','Gabriel',1,'$2a$10$ARGDNUz.xsfWAaS2KCG2T.h5N3d9NTf77i0Q5dp6FdpYXSJI08ijW','2024-01-23 03:23:20','2024-01-23 03:23:20');


GRANT ALL ON *.* to 'appuser'@'localhost' IDENTIFIED WITH mysql_native_password BY 'youShouldChangeThisPassword';
GRANT ALL ON *.* to 'appuser'@'mysqldb' IDENTIFIED WITH mysql_native_password BY 'youShouldChangeThisPassword';
GRANT ALL ON *.* to 'appuser'@'%' IDENTIFIED WITH mysql_native_password BY 'youShouldChangeThisPassword';

GRANT ALL ON *.* to 'root'@'localhost' IDENTIFIED WITH mysql_native_password BY 'youShouldChangeThisPassword';
GRANT ALL ON *.* to 'root'@'%' IDENTIFIED WITH mysql_native_password BY 'youShouldChangeThisPassword';
FLUSH PRIVILEGES;




