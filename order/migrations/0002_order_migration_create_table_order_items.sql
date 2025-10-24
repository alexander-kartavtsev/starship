-- +goose Up
create table if not exists order_items (
    id integer,
    order_id integer,
    part_id integer not null,
    quantity integer not null default 1,
    price numeric(10,2) not null,
    created_at timestamp not null default now(),
    updated_at timestamp not null on update now(),
    primary key (id, order_id),
    foreign key (part_id) references parts (id),
    foreign key (order_id) references orders (id),
);

-- +goose Down
drop table order_items;
