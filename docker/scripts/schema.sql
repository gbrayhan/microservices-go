CREATE DATABASE IF NOT EXISTS boilerplate_go;

USE boilerplate_go;

create table if not exists medicines
(
    id          bigint auto_increment primary key,
    name        varchar(30) charset utf8mb3         null,
    ean_code    varchar(30) charset utf8mb3         null,
    description varchar(150) charset utf8mb3        null,
    laboratory  varchar(50) charset utf8mb3         null,
    created_at  timestamp default CURRENT_TIMESTAMP null,
    updated_at  timestamp default CURRENT_TIMESTAMP null on update CURRENT_TIMESTAMP,
    constraint medicines_UN_ean unique (ean_code),
    constraint medicines_UN_name unique (name)
) CHARSET = utf8mb4
  COLLATE = utf8mb4_general_ci;

create table if not exists users
(
    id            bigint auto_increment,
    user          varchar(100) charset utf8mb4         not null,
    email         varchar(100) charset utf8mb4         not null,
    first_name    varchar(100) charset utf8mb4         null,
    last_name     varchar(100) charset utf8mb4         null,
    status        tinyint(1) default 1                 null,
    hash_password varchar(255) charset utf8mb4         not null,
    created_at    timestamp  default CURRENT_TIMESTAMP not null on update CURRENT_TIMESTAMP,
    updated_at    timestamp  default CURRENT_TIMESTAMP not null on update CURRENT_TIMESTAMP,
    constraint users_UN_email unique (email),
    constraint users_UN_user unique (user),
    constraint users_id_IDX unique (id)
) CHARSET = utf8mb4
  COLLATE = utf8mb4_general_ci;

FLUSH PRIVILEGES;




