-- +goose Up
create table users (
    id        bigserial primary key,
    name      varchar(255) not null,
    password  varchar(255) not null,
    email     varchar(255) unique not null,
    address   varchar(255) not null,
    role      int not null,    
    updated_at timestamptz default now(),
    created_at timestamptz default now(),
    deleted_at timestamptz default null
);

insert into users(name, password, email, address, role) values ('Seller','$2a$12$yT.dJTZnu4FRJq9zXw0mBOA/xmZHJPVi5ni13Zk9Pn6E0QmwKkZTu','seller@mail.com','Address',1);
insert into users(name, password, email, address, role) values ('Buyer','$2a$12$yT.dJTZnu4FRJq9zXw0mBOA/xmZHJPVi5ni13Zk9Pn6E0QmwKkZTu','buyer@mail.com','Address',2);

-- +goose Down
drop table users;