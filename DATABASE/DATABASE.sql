CREATE DATABASE `registro`;

USE `registro`;

CREATE TABLE clientes (
    id INT AUTO_INCREMENT PRIMARY KEY,
    email VARCHAR(255) NOT NULL,
    senha VARCHAR(255) NOT NULL
);

