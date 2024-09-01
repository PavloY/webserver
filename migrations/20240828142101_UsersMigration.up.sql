CREATE TABLE users (
    id bigserial not null primary key,
    login varchar not null unique
);

CREATE TABLE products (
    id bigserial not null primary key,
    name varchar not null 
);