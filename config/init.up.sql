CREATE TABLE IF NOT EXISTS users (
    id varchar(100) not null primary key,
    name varchar(255) not null,
    email varchar(255) not null unique,
    created_at timestamp not null default now(),
    updated_at timestamp default null
);

CREATE TABLE IF NOT EXISTS posts (
    id varchar(100) not null primary key,
    user_id varchar(100) not null,
    description varchar(255) not null,
    FOREIGN KEY (user_id) REFERENCES users(id),
    created_at timestamp not null default now(),
    updated_at timestamp default null
);