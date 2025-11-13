-- +goose Up
create table order_items (
    order_uuid varchar(36),
    part_uuid varchar(36) not null,
    created_at timestamp not null default now(),
    updated_at timestamp not null default now(),
    primary key (order_uuid, part_uuid),
    foreign key (order_uuid) references orders (order_uuid) on delete cascade
);

-- +goose Down
drop table order_items;
