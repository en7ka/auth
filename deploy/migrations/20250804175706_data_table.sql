-- +goose Up
create table users (
      id SERIAL PRIMARY KEY,
      username VARCHAR(100) NOT NULL UNIQUE,
      email VARCHAR(100) NOT NULL UNIQUE,
      role VARCHAR(20) not null,
      created_at timestamp not null default now()
);

-- +goose Down
drop table users;