create database if not exists test_user_service;
use test_user_service;
create table if not exists users
(
    id     int primary key auto_increment,
    name   varchar(255),
    status ENUM ('ACTIVE', 'INACTIVE') NOT NULL DEFAULT 'ACTIVE',
    gender ENUM ('FEMALE','MALE'),
    unique (name)
);