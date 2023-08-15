-- +migrate Up
create table users (
  id serial primary key,
  username varchar unique,
  password varchar,
  fullname varchar,
  email varchar,
  phone varchar,
  created_at timestamptz default now(),
  updated_at timestamptz,
  deleted_at timestamptz
);

-- +migrate Down
drop table users;
