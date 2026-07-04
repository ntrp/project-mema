-- +goose Up
alter table app.users
    add column if not exists display_name text not null default '',
    add column if not exists picture_url text not null default '';
