-- +goose Up
-- +goose StatementBegin
create schema if not exists users_schema;

CREATE EXTENSION IF NOT EXISTS "pgcrypto";

drop table if exists users_schema.users cascade;

CREATE TABLE if not exists users_schema.users (
    id            SERIAL PRIMARY KEY,
    login         TEXT unique not null,
    password_hash TEXT        NOT NULL
);

drop table if exists users_schema.account;

create table if not exists users_schema.account (
    login text primary key not null references users_schema.users(login),
    balance int default 1000 CHECK (balance >= 0)
);

drop table if exists users_schema.user_items;

create table if not exists users_schema.user_items (
    login text not null references users_schema.users(login),
    item text not null
);

CREATE INDEX IF NOT EXISTS idx_user_items ON users_schema.user_items (login);

drop table if exists users_schema.user_operations;

create table if not exists users_schema.user_operations (
    sender text not null references users_schema.users(login),
    recipient text not null references users_schema.users(login),
    amount int not null
);

CREATE INDEX IF NOT EXISTS idx_user_sender ON users_schema.user_operations (sender);
CREATE INDEX IF NOT EXISTS idx_user_recipient ON users_schema.user_operations (recipient);

create schema if not exists auth_schema;

drop table if exists auth_schema.users_secrets;

create table if not exists auth_schema.users_secrets (
    login text not null references users_schema.users(login),
    secret text not null,
    session_id text not null,
    unique (login, session_id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop schema users_schema cascade;
drop schema auth_schema cascade;
-- +goose StatementEnd
