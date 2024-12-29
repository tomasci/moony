-- +goose Up
-- +goose StatementBegin
create table users
(
    id         uuid not null default gen_random_uuid() primary key,
    username   text not null
        constraint users_username_unique unique,
    password   text not null,
    email      text not null
        constraint users_email_unique unique,
    created_at timestamp with time zone,
    updated_at timestamp with time zone,
    deleted_at timestamp with time zone
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table if exists users;
-- +goose StatementEnd
