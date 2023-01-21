DROP TABLE IF EXISTS posts cascade;
DROP TABLE IF EXISTS users cascade;
DROP TABLE IF EXISTS posts cascade;
DROP TABLE IF EXISTS user_follows cascade;

CREATE TABLE IF NOT EXISTS users (
    id varchar(100) not null primary key default uuid_generate_v4(),
    username varchar(255) not null,
    email varchar(255) not null unique,
    password varchar(255) not null,
    created_at timestamp not null default now(),
    updated_at timestamp default null
);

CREATE TABLE IF NOT EXISTS posts (
    id varchar(100) not null primary key,
    user_id varchar(100) not null,
    description text not null,
    FOREIGN KEY (user_id) REFERENCES users(id),
    created_at timestamp not null default now(),
    updated_at timestamp default null
);

CREATE TABLE IF NOT EXISTS user_follows (
    user_id varchar(100) not null,
    user_follow_id varchar(100) not null,
    created_at timestamp not null default now(),
    FOREIGN KEY (user_id) REFERENCES users(id),
    FOREIGN KEY (user_follow_id) REFERENCES users(id)
);