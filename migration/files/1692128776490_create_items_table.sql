-- +migrate Up
create table items (
  id serial primary key,
  item_name varchar,
  item_price int,
  item_stock int,
  created_at timestamptz default now(),
  updated_at timestamptz,
  deleted_at timestamptz
);

-- +migrate Down
drop table items;
