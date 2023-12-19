create table units (
  id serial primary key,
  name varchar(255) not null,
  description text,
  content text not null,
  sort_key int not null,
  chapter_id int references chapters(id),
  created_at timestamp default current_timestamp,
  updated_at timestamp default current_timestamp,
  deleted_at timestamp
)
