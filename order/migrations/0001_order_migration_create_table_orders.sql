-- +goose Up
create table orders (
    order_uuid varchar(36) primary key,
    user_uuid varchar(36) not null,
    transaction_uuid varchar(36) not null default '',
    total_price float not null,
    payment_method varchar(20),
    status varchar(20) not null,
    created_at timestamp not null default current_timestamp,
    updated_at timestamp not null default current_timestamp
);

-- +goose Down
drop table orders;
