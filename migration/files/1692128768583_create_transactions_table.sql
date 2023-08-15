-- +migrate Up
create table transactions (
  id serial primary key,
  order_id int not null REFERENCES orders(id),
  item_name varchar,
  item_price int,
  item_quantity int,
  created_at timestamptz default now(),
  updated_at timestamptz,
  deleted_at timestamptz
);

-- +migrate Down
drop table transactions;
