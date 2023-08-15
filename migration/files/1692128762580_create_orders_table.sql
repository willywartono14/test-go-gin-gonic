-- +migrate Up
create table orders (
  id serial primary key,
  user_id int not null REFERENCES users(id),
  status varchar,
  invoice_number varchar unique,
  created_at timestamptz default now(),
  updated_at timestamptz,
  deleted_at timestamptz
);

-- +migrate Down
drop table orders;
