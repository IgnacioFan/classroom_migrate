create table chapters (
  id serial primary key,
  name varchar(255) not null,
  course_id int references courses(id),
  created_at timestamp default current_timestamp,
  updated_at timestamp default current_timestamp,
  deleted_at timestamp
) 
