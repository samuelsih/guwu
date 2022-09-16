CREATE TABLE IF NOT EXISTS users (
    id varchar(100) not null primary key,
    username varchar(255) not null,
    email varchar(255) not null unique,
    password varchar(255) not null,
    created_at timestamp not null default now(),
    updated_at timestamp default null
);