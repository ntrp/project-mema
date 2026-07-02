-- +goose Up
alter table app.media_release_candidates
    add column if not exists published_at timestamptz;
