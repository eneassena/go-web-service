drop database if exists bootcampgo;

create schema bootcampgo 
	default character set utf8 
	default collate utf8_general_ci;
    
use bootcampgo;


create table products (
	`id` int not null auto_increment primary key,
	`name` varchar(100) not null,
	`type` varchar(100) not null,
	`count` int(11) not null,
	`price` decimal(5,2) not null
);