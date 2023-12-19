create table courses (
  id serial primary key,
  name varchar(255) not null,
  description text,
  created_at timestamp default current_timestamp,
  updated_at timestamp default current_timestamp,
  deleted_at timestamp
);
